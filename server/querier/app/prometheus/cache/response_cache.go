/*
 * Copyright (c) 2024 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cache

import (
	"fmt"
	"sort"
	"sync"
	"time"
	"unsafe"

	"github.com/deepflowio/deepflow/server/libs/lru"
	"github.com/deepflowio/deepflow/server/querier/config"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
)

const (
	pointSize    = int(unsafe.Sizeof(promql.Point{}))
	pointPtrSize = int(unsafe.Sizeof(&promql.Point{}))
)

type item[T any] struct {
	start int64
	end   int64
	step  int64
	vType parser.ValueType // origin value type
	data  T

	loadCompleted chan struct{}
}

func dataSizeOf(r *item[promql.Result]) uint64 {
	if r == nil {
		return 0
	}
	var size uintptr
	size = unsafe.Sizeof(r.start) + unsafe.Sizeof(r.end) + unsafe.Sizeof(r.step) + unsafe.Sizeof(r.vType)
	matrix, err := r.data.Matrix()
	if err != nil {
		totalPoints := 0
		for _, m := range matrix {
			for _, v := range m.Metric {
				size += uintptr(len(v.Name) + len(v.Value))
				size += unsafe.Sizeof((*string)(unsafe.Pointer(&v.Name))) + unsafe.Sizeof((*string)(unsafe.Pointer(&v.Value)))
			}
			totalPoints += len(m.Points)
		}
		size += uintptr(totalPoints*pointSize + totalPoints*pointPtrSize)
	}
	return uint64(size)
}

type Cacher struct {
	entries *lru.Cache[string, item[promql.Result]]
	lock    *sync.RWMutex

	cleanUpCache *time.Ticker
}

func NewCacher() *Cacher {
	c := &Cacher{
		lock:    &sync.RWMutex{},
		entries: lru.NewCache[string, item[promql.Result]](config.Cfg.Prometheus.Cache.CacheMaxCount),
	}
	go c.startUpCleanCache(config.Cfg.Prometheus.Cache.CacheCleanInterval)
	return c
}

func (c *Cacher) startUpCleanCache(cleanUpInterval int) {
	c.cleanUpCache = time.NewTicker(time.Duration(cleanUpInterval) * time.Second)
	defer func() {
		c.cleanUpCache.Stop()
		if err := recover(); err != nil {
			go c.startUpCleanCache(cleanUpInterval)
		}
	}()
	for range c.cleanUpCache.C {
		c.cleanCache()
	}
}

func (c *Cacher) cleanCache() {
	keys := c.entries.Keys()
	for _, k := range keys {
		item, ok := c.entries.Peek(k)
		if !ok {
			continue
		}
		size := dataSizeOf(&item)
		if size > config.Cfg.Prometheus.Cache.CacheItemSize {
			log.Infof("cache item remove: %s, real size: %d", k, size)
			c.entries.Remove(k)
		}
	}
}

func (c *Cacher) Fetch(key string, start, end int64) (r promql.Result, fixedStart int64, fixedEnd int64, queryRequired bool) {
	entry, ok := c.entries.Get(key)
	if !ok {
		c.lock.Lock()
		c.entries.Add(key, item[promql.Result]{vType: parser.ValueTypeNone, loadCompleted: make(chan struct{})})
		c.lock.Unlock()

		return promql.Result{Err: fmt.Errorf("key %s not found", key)}, start, end, true
	}

	if entry.vType == parser.ValueTypeNone {
		select {
		case <-time.After(time.Duration(config.Cfg.Prometheus.Cache.CacheFirstTimeout) * time.Second):
			log.Infof("req [%s:%d-%d] wait %d to get cache result", key, start, end, config.Cfg.Prometheus.Cache.CacheFirstTimeout)
			return promql.Result{Err: fmt.Errorf("key %s not found", key)}, start, end, false
		case <-entry.loadCompleted:
			entry, ok = c.entries.Get(key)
			if !ok {
				return promql.Result{Err: fmt.Errorf("key %s load failed", key)}, start, end, true
			}
			log.Debugf("req [%s:%d-%d] get cached result", key, start, end)
		}
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	if entry.vType == parser.ValueTypeVector {
		r, queryRequired = c.fetchInstant(&entry, start, end)
		return r, start, end, queryRequired
	} else if entry.vType == parser.ValueTypeMatrix {
		return c.fetchRange(&entry, start, end)
	}

	return promql.Result{Err: fmt.Errorf("value Type %s not found", key)}, start, end, true
}

func (c *Cacher) fetchRange(entry *item[promql.Result], start, end int64) (r promql.Result, fixedStart int64, fixedEnd int64, queryRequired bool) {
	fixedStart, fixedEnd = c.validateQueryTs(start, end, entry.start, entry.end)
	if fixedStart == 0 && fixedEnd == 0 {
		return c.extractSubData(entry.data, start, end), fixedStart, fixedEnd, false
	}
	return entry.data, fixedStart, fixedEnd, true
}

func (c *Cacher) fetchInstant(entry *item[promql.Result], start, end int64) (r promql.Result, queryRequired bool) {
	// for instant query, all data will storage as matrix, but depends on query promql, get difference result:
	// query node_cpu_seconds_total: get vector, find out the matched query time
	// query node_cpu_seconds_total[1m]: get matrix, find out the matched query time range
	samples, err := entry.data.Matrix()
	if err != nil {
		return promql.Result{Err: err}, true
	}
	result := make(promql.Vector, 0, samples.Len())

	// only when end == Points.T, can be added (time completely equal)
	for i := 0; i < samples.Len(); i++ {
		findEndTimeAt := sort.Search(len(samples[i].Points), func(j int) bool {
			return samples[i].Points[j].T == end
		})
		if findEndTimeAt == len(samples[i].Points) {
			// not found
			// fmt.Errorf("time %v not found", end)
			continue
		}

		result = append(result, promql.Sample{Metric: samples[i].Metric, Point: samples[i].Points[findEndTimeAt]})
	}
	if len(result) < samples.Len() {
		return promql.Result{}, true
	}
	return promql.Result{Value: result}, false
}

func (c *Cacher) validateQueryTs(start, end int64, cacheStart, cacheEnd int64) (int64, int64) {
	// cache miss:
	// left
	if end < cacheStart {
		return start, end
	}
	// right
	if start > cacheEnd {
		return start, end
	}
	// outside
	if start < cacheStart && end > cacheEnd {
		return start, end
	}

	// cache hit
	if start < cacheStart {
		return start, cacheStart
	}

	if end > cacheEnd {
		return cacheEnd, end
	}

	// cache hit, not query anything, return cache data
	if start >= cacheStart && cacheEnd >= end {
		return 0, 0
	}

	return start, end
}

func (c *Cacher) extractSubData(r promql.Result, start, end int64) promql.Result {
	matrix, err := r.Matrix()
	if err != nil {
		return promql.Result{Err: err}
	}
	result := make(promql.Matrix, 0, len(matrix))
	for _, series := range matrix {
		newSeries := promql.Series{Metric: series.Metric}
		begin := sort.Search(len(series.Points), func(i int) bool {
			return series.Points[i].T >= start
		})
		stop := sort.Search(len(series.Points), func(i int) bool {
			return series.Points[i].T > end // include end
		})
		newSeries.Points = series.Points[begin:stop]
		result = append(result, newSeries)
	}
	return promql.Result{Value: result}
}

func (c *Cacher) Merge(key string, start, end, step int64, res promql.Result) (promql.Result, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	entry, ok := c.entries.Get(key)
	if !ok || (entry.vType == parser.ValueTypeNone) {
		item := item[promql.Result]{
			start: start,
			end:   end,
			step:  step,
			data:  res,
			vType: res.Value.Type(),
		}

		if item.data.Value.Type() == parser.ValueTypeVector {
			v, err := res.Vector()
			if err == nil {
				item.data.Value = vectorTomatrix(&v)
				item.vType = parser.ValueTypeVector // mark origin value type
			}
			// else not vector, merge directly
		}
		c.entries.Add(key, item)
		return res, nil
	}

	var mergeResult promql.Result
	switch res.Value.Type() {
	case parser.ValueTypeVector:
		vector, err := res.Vector()
		if err != nil {
			return res, err
		}
		mergeResult, err = c.vectorMerge(vector, &entry.data)
		if err != nil {
			return promql.Result{Err: err}, err
		}
	case parser.ValueTypeMatrix:
		// replace
		if start > entry.end || end < entry.start {
			entry.start = start
			entry.end = end
			entry.data = res
		}

		if start < entry.start && end > entry.end {
			entry.start = start
			entry.end = end
			entry.data = res
		}

		matrix, err := res.Matrix()
		if err != nil {
			return res, err
		}
		mergeResult, err = c.matrixMerge(matrix, &entry.data)
		if err != nil {
			return promql.Result{Err: err}, err
		}
	default:
	}

	if entry.end < end {
		entry.end = end
	}
	if entry.start > start {
		entry.start = start
	}

	entry.data = mergeResult

	if entry.loadCompleted != nil {
		select {
		case _, ok := <-entry.loadCompleted:
			log.Debugf("entry loadCompleted closed: %v", ok)
		default:
			close(entry.loadCompleted)
		}
	}
	return entry.data, nil
}

func (c *Cacher) matrixMerge(resp promql.Matrix, cache *promql.Result) (promql.Result, error) {
	cacheMatrix, err := cache.Matrix()
	if err != nil {
		return promql.Result{Err: err}, err
	}
	output := make(promql.Matrix, 0, len(cacheMatrix))
	for _, cachedTs := range cacheMatrix {
		newSeries := promql.Series{Metric: cachedTs.Metric}
		newSeries.Points = cachedTs.Points
		for _, series := range resp {
			if promLabelsEqual(&cachedTs.Metric, &series.Metric) {
				existsStartT := newSeries.Points[0].T
				existsEndT := newSeries.Points[len(newSeries.Points)-1].T

				if existsEndT < series.Points[0].T {
					newSeries.Points = append(newSeries.Points, series.Points...)
				} else if existsStartT > series.Points[len(series.Points)-1].T {
					newSeries.Points = append(series.Points, newSeries.Points...)
				} else if existsEndT >= series.Points[0].T && existsEndT < series.Points[len(series.Points)-1].T {
					// cached data & resp overlap
					overlapPointAt := sort.Search(len(series.Points), func(i int) bool {
						return series.Points[i].T > existsEndT
					})
					newSeries.Points = append(newSeries.Points, series.Points[overlapPointAt:]...)
				} else if existsStartT <= series.Points[len(series.Points)-1].T && existsStartT > series.Points[0].T {
					overlapPointAt := sort.Search(len(series.Points), func(i int) bool {
						return series.Points[i].T >= existsStartT
					})
					newSeries.Points = append(series.Points[:overlapPointAt], newSeries.Points...)
				}

				sort.Slice(newSeries.Points, func(i, j int) bool {
					return newSeries.Points[i].T < newSeries.Points[j].T
				})
			}
		}
		output = append(output, newSeries)
	}

	return promql.Result{Value: output}, nil
}

func (c *Cacher) vectorMerge(resp promql.Vector, cached *promql.Result) (promql.Result, error) {
	cacheMatrix, err := cached.Matrix()
	if err != nil {
		return promql.Result{Err: err}, err
	}
	output := make(promql.Matrix, 0, len(cacheMatrix))
	for _, cachedTs := range cacheMatrix {
		newSeries := promql.Series{Metric: cachedTs.Metric}
		newSeries.Points = cachedTs.Points
		for _, samples := range resp {
			if promLabelsEqual(&cachedTs.Metric, &samples.Metric) {
				insertedPointAt := sort.Search(len(newSeries.Points), func(i int) bool {
					return newSeries.Points[i].T >= samples.Point.T
				})
				newSeries.Points = append(newSeries.Points, promql.Point{})
				copy(newSeries.Points[insertedPointAt+1:], newSeries.Points[insertedPointAt:])
				newSeries.Points[insertedPointAt] = samples.Point

				sort.Slice(newSeries.Points, func(i, j int) bool {
					return newSeries.Points[i].T < newSeries.Points[j].T
				})
			}
		}
	}
	return promql.Result{Value: output}, nil
}

func vectorTomatrix(v *promql.Vector) promql.Matrix {
	output := make(promql.Matrix, 0, len(*v))
	for _, m := range *v {
		output = append(output, promql.Series{
			Metric: m.Metric,
			Points: []promql.Point{m.Point},
		})
	}
	return output
}
