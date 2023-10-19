/**
 * Copyright (c) 2023 Yunshan Networks
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

package mysql

import (
	"time"

	"gorm.io/gorm"
)

type ResourceEvent struct {
	ID             int       `gorm:"primaryKey;autoIncrement;unique;column:id;type:int;not null" json:"ID"`
	Domain         string    `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
	SubDomain      string    `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	ResourceLcuuid string    `gorm:"column:resource_lcuuid;type:char(64);default:''" json:"RESOURCE_LCUUID"`
	Content        string    `gorm:"column:content;type:text" json:"CONTENT"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at;type:datetime" json:"CREATED_AT"`
}

type DomainAdditionalResource struct {
	ID                int             `gorm:"primaryKey;autoIncrement;unique;column:id;type:int;not null" json:"ID"`
	Domain            string          `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
	Content           string          `gorm:"column:content;type:longtext" json:"CONTENT"`
	CompressedContent compressedBytes `gorm:"column:compressed_content;type:longblob" json:"COMPRESSED_CONTENT"`
	CreatedAt         time.Time       `gorm:"autoCreateTime;column:created_at;type:datetime" json:"CREATED_AT"`
}

type Base struct {
	ID     int    `gorm:"primaryKey;autoIncrement;unique;column:id;type:int;not null" json:"ID"`
	Lcuuid string `gorm:"unique;column:lcuuid;type:char(64)" json:"LCUUID"`
	// TODO add CreatedAt/UpdatedAt/DeletedAt
}

func (b Base) GetID() int {
	return b.ID
}

func (b *Base) SetID(id int) {
	b.ID = id
}

func (b Base) GetLcuuid() string {
	return b.Lcuuid
}

type OperatedTime struct {
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;type:datetime" json:"CREATED_AT"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;type:datetime" json:"UPDATED_AT"`
}

type SoftDeleteBase struct {
	OperatedTime `json:",squash"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;default:null" json:"DELETED_AT"`
}

type Process struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string    `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	VTapID         uint32    `gorm:"column:vtap_id;type:int;not null;default:0" json:"VTAP_ID"`
	PID            uint64    `gorm:"column:pid;type:int;not null;default:0" json:"PID"`
	ProcessName    string    `gorm:"column:process_name;type:varchar(256);default:''" json:"PROCESS_NAME"`
	CommandLine    string    `gorm:"column:command_line;type:text" json:"COMMAND_LINE"`
	UserName       string    `gorm:"column:user_name;type:varchar(256);default:''" json:"USER_NAME"`
	StartTime      time.Time `gorm:"autoCreateTime;column:start_time;type:datetime" json:"START_TIME"`
	OSAPPTags      string    `gorm:"column:os_app_tags;type:text" json:"OS_APP_TAGS"`
	ContainerID    string    `gorm:"column:container_id;type:char(64);default:''" json:"CONTAINER_ID"`
	NetnsID        uint32    `gorm:"column:netns_id;type:int unsigned;default:0" json:"NETNS_ID"` // used to associate processes with cloud and container resources
	SubDomain      string    `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string    `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
}

type Domain struct {
	Base         `gorm:"embedded" json:",squash"`
	OperatedTime `gorm:"embedded" json:",squash"`
	SyncedAt     *time.Time `gorm:"column:synced_at" json:"SYNCED_AT"`
	Name         string     `gorm:"column:name;type:varchar(64)" json:"NAME"`
	IconID       int        `gorm:"column:icon_id;type:int" json:"ICON_ID"`
	DisplayName  string     `gorm:"column:display_name;type:varchar(64);default:''" json:"DISPLAY_NAME"`
	ClusterID    string     `gorm:"column:cluster_id;type:char(64)" json:"CLUSTER_ID"`
	Type         int        `gorm:"column:type;type:int;default:0" json:"TYPE"` // 1.openstack 2.vsphere 3.nsp 4.tencent 5.filereader 6.aws 7.pingan 8.zstack 9.aliyun 10.huawei prv 11.k8s 12.simulation 13.huawei 14.qingcloud 15.qingcloud_private 16.F5 17.CMB_CMDB 18.azure 19.apsara_stack 20.tencent_tce 21.qingcloud_k8s 22.kingsoft_private 23.genesis 24.microsoft_acs 25.baidu_bce
	Config       string     `gorm:"column:config;type:text" json:"CONFIG"`
	ErrorMsg     string     `gorm:"column:error_msg;type:text" json:"ERROR_MSG"`
	Enabled      int        `gorm:"column:enabled;type:int;not null;default:1" json:"ENABLED"` // 0.false 1.true
	State        int        `gorm:"column:state;type:int;not null;default:1" json:"STATE"`     // 1.normal 2.deleting 3.exception
	ControllerIP string     `gorm:"column:controller_ip;type:char(64)" json:"CONTROLLER_IP"`
}

// TODO 最终可以与cloud模块命名统一，Domain -> DomainLcuuid

type SubDomain struct {
	Base         `gorm:"embedded" json:",squash"`
	OperatedTime `gorm:"embedded" json:",squash"`
	SyncedAt     *time.Time `gorm:"column:synced_at" json:"SYNCED_AT"`
	Domain       string     `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
	Name         string     `gorm:"column:name;type:varchar(64);default:''" json:"NAME"`
	DisplayName  string     `gorm:"column:display_name;type:varchar(64);default:''" json:"DISPLAY_NAME"`
	CreateMethod int        `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	ClusterID    string     `gorm:"column:cluster_id;type:char(64);default:''" json:"CLUSTER_ID"`
	Config       string     `gorm:"column:config;type:text;default:''" json:"CONFIG"`
	ErrorMsg     string     `gorm:"column:error_msg;type:text;default:''" json:"ERROR_MSG"`
	Enabled      int        `gorm:"column:enabled;type:int;not null;default:1" json:"ENABLED"` // 0.false 1.true
	State        int        `gorm:"column:state;type:int;not null;default:1" json:"STATE"`     // 1.normal 2.deleting 3.exception
}

type Region struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string  `gorm:"column:name;type:varchar(64);default:''" json:"NAME"`
	CreateMethod   int     `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	Label          string  `gorm:"column:label;type:varchar(64);default:''" json:"LABEL"`
	Longitude      float64 `gorm:"column:longitude;type:double(7,4);default:null" json:"LONGITUDE"`
	Latitude       float64 `gorm:"column:latitude;type:double(7,4);default:null" json:"LATITUDE"`
}

type AZ struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(64);default:''" json:"NAME"`
	CreateMethod   int    `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	Label          string `gorm:"column:label;type:varchar(64);default:''" json:"LABEL"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	Domain         string `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
}

func (AZ) TableName() string {
	return "az"
}

type Host struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Type           int       `gorm:"column:type;type:int" json:"TYPE"`   // 1.Server 3.Gateway 4.DFI
	State          int       `gorm:"column:state;type:int" json:"STATE"` // 0.Temp 1.Creating 2.Complete 3.Modifying 4.Exception
	Name           string    `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Alias          string    `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Description    string    `gorm:"column:description;type:varchar(256);default:''" json:"DESCRIPTION"`
	IP             string    `gorm:"column:ip;type:char(64);default:''" json:"IP"`
	HType          int       `gorm:"column:htype;type:int" json:"HTYPE"`                           // 1. Xen host 2. VMware host 3. KVM host 4. Public cloud host 5. Hyper-V
	CreateMethod   int       `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	UserName       string    `gorm:"column:user_name;type:varchar(64);default:''" json:"USER_NAME"`
	UserPasswd     string    `gorm:"column:user_passwd;type:varchar(64);default:''" json:"USER_PASSWD"`
	VCPUNum        int       `gorm:"column:vcpu_num;type:int;default:0" json:"VCPU_NUM"`
	MemTotal       int       `gorm:"column:mem_total;type:int;default:0" json:"MEM_TOTAL"` // unit: M
	AZ             string    `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region         string    `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	Domain         string    `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
	SyncedAt       time.Time `gorm:"column:synced_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"SYNCED_AT"`
	ExtraInfo      string    `gorm:"column:extra_info;type:text;default:''" json:"EXTRA_INFO"`
}

func (Host) TableName() string {
	return "host_device"
}

type VM struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	State          int    `gorm:"index:state_server_index;column:state;type:int;not null" json:"STATE"` // 0.Temp 1.Creating 2.Created 3.To run 4.Running 5.To suspend 6.Suspended 7.To resume 8. To stop 9.Stopped 10.Modifing 11.Exception 12.Destroying
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	CreateMethod   int    `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	HType          int    `gorm:"column:htype;type:int;default:1" json:"HTYPE"`                 // 1.vm-c 2.bm-c 3.vm-n 4.bm-n 5.vm-s 6.bm-s
	LaunchServer   string `gorm:"index:state_server_index;column:launch_server;type:char(64);default:''" json:"LAUNCH_SERVER"`
	CloudTags      string `gorm:"column:cloud_tags;type:text;default:''" json:"CLOUD_TAGS"` // separated by ,
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"EPC_ID"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	UID            string `gorm:"column:uid;type:char(64);default:''" json:"UID"`
}

func (VM) TableName() string {
	return "vm"
}

type VMPodNodeConnection struct {
	Base      `gorm:"embedded" json:",squash"`
	VMID      int    `gorm:"column:vm_id;type:int;default:null" json:"VM_ID"`
	PodNodeID int    `gorm:"column:pod_node_id;type:int;default:null" json:"POD_NODE_ID"`
	Domain    string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	SubDomain string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
}

func (VMPodNodeConnection) TableName() string {
	return "vm_pod_node_connection"
}

type VMSecurityGroup struct {
	Base            `gorm:"embedded" json:",squash"`
	SecurityGroupID int `gorm:"column:sg_id;type:int;not null" json:"SG_ID"`
	VMID            int `gorm:"column:vm_id;type:int;not null" json:"VM_ID"`
	Priority        int `gorm:"column:priority;type:int;not null" json:"PRIORITY"`
}

func (VMSecurityGroup) TableName() string {
	return "vm_security_group"
}

type Contact struct {
	Base         `gorm:"embedded" json:",squash"`
	Name         string    `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	CreateMethod int       `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	Mobile       string    `gorm:"column:mobile;type:char(13);default:''" json:"MOBILE"`
	Email        string    `gorm:"column:email;type:varchar(128);default:''" json:"EMAIL"`
	Company      string    `gorm:"column:company;type:varchar(128);default:''" json:"COMPANY"`
	PushEmail    string    `gorm:"column:push_email;type:text;default:null" json:"PUSH_EMAIL"`
	Domain       string    `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
	AlarmPush    int       `gorm:"column:alarm_push;type:int;default:0" json:"ALARM_PUSH"`
	ReportPush   int       `gorm:"column:report_push;type:int;default:0" json:"REPORT_PUSH"`
	Deleted      int       `gorm:"column:deleted;type:int;default:0" json:"DELETED"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;default:null" json:"UPDATED_AT"`
}

type VPCContact struct {
	Base         `gorm:"embedded" json:",squash"`
	CreateMethod int `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	VPCID        int `gorm:"column:epc_id;type:int;default:0" json:"VPC_ID"`
	ContactID    int `gorm:"column:contact_id;type:int;default:0" json:"CONTACT_ID"`
}

func (VPCContact) TableName() string {
	return "epc_contact"
}

type VPC struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	CreateMethod   int    `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	Label          string `gorm:"column:label;type:varchar(64);default:''" json:"LABEL"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Domain         string `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"` // TODO delete in future
	TunnelID       int    `gorm:"column:tunnel_id;type:int;default:0" json:"TUNNEL_ID"`
	Mode           int    `gorm:"column:mode;type:int;default:2" json:"MODE"` //  1:route, 2:transparent
	CIDR           string `gorm:"column:cidr;type:char(64);default:''" json:"CIDR"`
	UID            string `gorm:"column:uid;type:char(64);default:''" json:"UID"`
}

func (VPC) TableName() string {
	return "epc"
}

type Network struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	State          int    `gorm:"column:state;type:int;not null" json:"STATE"`        // 0.Temp 1.Creating 2.Created 3.Exception 4.Modifing 5.Destroying 6.Destroyed
	NetType        int    `gorm:"column:net_type;type:int;default:4" json:"NET_TYPE"` // 1.CTRL 2.SERVICE 3.WAN 4.LAN
	Name           string `gorm:"column:name;type:varchar(256);not null" json:"NAME"`
	CreateMethod   int    `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	Label          string `gorm:"column:label;type:varchar(64);default:''" json:"LABEL"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Description    string `gorm:"column:description;type:varchar(256);default:''" json:"DESCRIPTION"`
	SubDomain      string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	ISP            int    `gorm:"column:isp;type:int;default:0" json:"ISP"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"VPC_ID"`
	SegmentationID int    `gorm:"column:segmentation_id;type:int;default:0" json:"SEGMENTATION_ID"`
	TunnelID       int    `gorm:"column:tunnel_id;type:int;default:0" json:"TUNNEL_ID"`
	Shared         bool   `gorm:"column:shared;type:int;default:0" json:"SHARED"`
	Topped         int    `gorm:"column:topped;type:int;default:0" json:"TOPPED"`
	IsVIP          int    `gorm:"column:is_vip;type:int;default:0" json:"IS_VIP"`
}

func (Network) TableName() string {
	return "vl2"
}

type Subnet struct {
	Base      `gorm:"embedded" json:",squash"`
	Prefix    string `gorm:"column:prefix;type:char(64);default:''" json:"PREFIX"`
	Netmask   string `gorm:"column:netmask;type:char(64);default:''" json:"NETMASK"`
	NetworkID int    `gorm:"column:vl2id;type:int;default:null" json:"VL2ID"`
	NetIndex  int    `gorm:"column:net_index;type:int;default:0" json:"NET_INDEX"`
	Name      string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label     string `gorm:"column:label;type:varchar(64);default:''" json:"LABEL"`
	SubDomain string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
}

func (Subnet) TableName() string {
	return "vl2_net"
}

type VRouter struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	State          int    `gorm:"index:state_server_index;column:state;type:int;not null" json:"STATE"` // 0.Temp 1.Creating 2.Created 3.Exception 4.Modifing 5.Destroying 6.To run 7.Running 8.To stop 9.Stopped
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	Description    string `gorm:"column:description;type:varchar(256);default:''" json:"DESCRIPTION"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"EPC_ID"`
	GWLaunchServer string `gorm:"index:state_server_index;column:gw_launch_server;type:char(64);default:''" json:"GW_LAUNCH_SERVER"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"` // TODO delete in future
}

func (VRouter) TableName() string {
	return "vnet"
}

type RoutingTable struct {
	Base        `gorm:"embedded" json:",squash"`
	VRouterID   int    `gorm:"column:vnet_id;type:int;default:null" json:"VNET_ID"`
	Destination string `gorm:"column:destination;type:text;default:''" json:"DESTINATION"`
	NexthopType string `gorm:"column:nexthop_type;type:text;default:''" json:"NEXTHOP_TYPE"`
	Nexthop     string `gorm:"column:nexthop;type:text;default:''" json:"NEXTHOP"`
}

type DHCPPort struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"VPC_ID"`
}

func (DHCPPort) TableName() string {
	return "dhcp_port"
}

type VInterface struct {
	Base         `gorm:"embedded" json:",squash"`
	Name         string    `gorm:"column:name;type:char(64);default:''" json:"NAME"`
	Index        int       `gorm:"column:ifindex;type:int;not null" json:"IFINDEX"`
	State        int       `gorm:"column:state;type:int;not null" json:"STATE"`                  // 1. Attached 2.Detached 3.Exception
	CreateMethod int       `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	Type         int       `gorm:"column:iftype;type:int;default:0" json:"IFTYPE"`               // 0.Unknown 1.Control 2.Service 3.WAN 4.LAN 5.Trunk 6.Tap 7.Tool
	Mac          string    `gorm:"index:mac_index;column:mac;type:char(32);default:''" json:"MAC"`
	VMac         string    `gorm:"column:vmac;type:char(32);default:''" json:"VMAC"`
	TapMac       string    `gorm:"column:tap_mac;type:char(32);default:''" json:"TAP_MAC"`
	NetworkID    int       `gorm:"column:subnetid;type:int;default:0" json:"SUBNETID"` // vl2 id
	VlanTag      int       `gorm:"column:vlantag;type:int;default:0" json:"VLANTAG"`
	DeviceType   int       `gorm:"column:devicetype;type:int;default:null" json:"DEVICETYPE"`   // Type 0.unknown 1.vm 2.vgw 3.third-party-device 4.vmwaf 5.NSP-vgateway 6.host-device 7.network-device 9.DHCP-port 10.pod 11.pod_service 12. redis_instance 13. rds_instance 14. pod_node 15. load_balance 16. nat_gateway
	DeviceID     int       `gorm:"column:deviceid;type:int;default:null" json:"DEVICEID"`       // unknown: Senseless ID, vm: vm ID, vgw/NSP-vgateway: vnet ID, third-party-device: third_party_device ID, vmwaf: vmwaf ID, host-device: host_device ID, network-device: network_device ID
	NetnsID      uint32    `gorm:"column:netns_id;type:int unsigned;default:0" json:"NETNS_ID"` // used to associate processes with cloud and container resources
	VtapID       uint32    `gorm:"column:vtap_id;type:int;default:0" json:"VTAP_ID"`
	SubDomain    string    `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain       string    `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	Region       string    `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
}

func (VInterface) TableName() string {
	return "vinterface"
}

type LANIP struct { // TODO 添加region字段
	Base         `gorm:"embedded" json:",squash"`
	IP           string    `gorm:"column:ip;type:char(64);default:''" json:"IP"`
	Netmask      string    `gorm:"column:netmask;type:char(64);default:''" json:"NETMASK"`
	Gateway      string    `gorm:"column:gateway;type:char(64);default:''" json:"GATEWAY"`
	CreateMethod int       `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	NetworkID    int       `gorm:"column:vl2id;type:int;default:null" json:"VL2ID"`
	NetIndex     int       `gorm:"column:net_index;type:int;default:0" json:"NET_INDEX"`
	SubDomain    string    `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain       string    `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	VInterfaceID int       `gorm:"column:vifid;type:int;default:null" json:"VINTERFACE_ID"`
	SubnetID     int       `gorm:"column:vl2_net_id;type:int;default:0" json:"SUBNET_ID"`
	ISP          int       `gorm:"column:isp;type:int;default:0" json:"ISP"` // Used for multi-ISP access
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
}

func (LANIP) TableName() string {
	return "vinterface_ip"
}

type WANIP struct {
	Base         `gorm:"embedded" json:",squash"`
	IP           string    `gorm:"column:ip;type:char(64);default:''" json:"IP"`
	Alias        string    `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Netmask      int       `gorm:"column:netmask;type:int;default:null" json:"NETMASK"`
	Gateway      string    `gorm:"column:gateway;type:char(64);default:''" json:"GATEWAY"`
	CreateMethod int       `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	ISP          int       `gorm:"column:isp;type:int;default:null" json:"ISP"`
	VInterfaceID int       `gorm:"column:vifid;type:int;default:0" json:"VINTERFACE_ID"`
	SubnetID     int       `gorm:"column:vl2_net_id;type:int;default:0" json:"SUBNET_ID"`
	SubDomain    string    `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain       string    `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	Region       string    `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
}

func (WANIP) TableName() string {
	return "ip_resource"
}

type FloatingIP struct {
	Base      `gorm:"embedded" json:",squash"`
	Domain    string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	Region    string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	VPCID     int    `gorm:"column:epc_id;type:int;default:0" json:"VPC_ID"`
	NetworkID int    `gorm:"column:vl2_id;type:int;default:null" json:"VL2_ID"` // TODO json字段是否能修改，需返回给前端？
	VMID      int    `gorm:"column:vm_id;type:int;default:null" json:"VM_ID"`
	IP        string `gorm:"column:ip;type:char(64);default:''" json:"IP"`
}

func (FloatingIP) TableName() string {
	return "floatingip"
}

type SecurityGroup struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label          string `gorm:"column:label;type:varchar(64);default:''" json:"LABEL"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"EPC_ID"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	Topped         int    `gorm:"column:topped;type:int;default:0" json:"TOPPED"`
}

type SecurityGroupRule struct {
	Base            `gorm:"embedded" json:",squash"`
	SecurityGroupID int    `gorm:"column:sg_id;type:int;not null" json:"SG_ID"`
	Direction       int    `gorm:"column:direction;type:tinyint(1);not null;default:0" json:"DIRECTION"` // 0.Unknow 1.Ingress 2.Egress
	Protocol        string `gorm:"column:protocol;type:char(64);default:''" json:"PROTOCOL"`
	EtherType       int    `gorm:"column:ethertype;type:tinyint(1);not null;default:0" json:"ETHERTYPE"` // 0.Unknow 1.IPv4 2.IPv6
	LocalPortRange  string `gorm:"column:local_port_range;type:text;default:''" json:"LOCAL_PORT_RANGE"`
	RemotePortRange string `gorm:"column:remote_port_range;type:text;default:''" json:"REMOTE_PORT_RANGE"`
	Local           string `gorm:"column:local;type:text;default:''" json:"LOCAL"`
	Remote          string `gorm:"column:remote;type:text;default:''" json:"REMOTE"`
	Priority        int    `gorm:"column:priority;type:int;not null" json:"PRIORITY"`
	Action          int    `gorm:"column:action;type:tinyint(1);not null;default:0" json:"ACTION"` // 0.Unknow 1.Accept 2.Drop
}

type NATGateway struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	FloatingIPs    string `gorm:"column:floating_ips;type:text;default:''" json:"FLOATING_IPS"` // separated by ,
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"EPC_ID"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"` // TODO delete in future
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	UID            string `gorm:"column:uid;type:char(64);default:''" json:"UID"`
}

func (NATGateway) TableName() string {
	return "nat_gateway"
}

type NATRule struct {
	Base           `gorm:"embedded" json:",squash"`
	NATGatewayID   int    `gorm:"column:nat_id;type:int;default:0" json:"NAT_ID"`
	Type           string `gorm:"column:type;type:char(16);default:''" json:"TYPE"`
	Protocol       string `gorm:"column:protocol;type:char(64);default:''" json:"PROTOCOL"`
	FloatingIP     string `gorm:"column:floating_ip;type:char(64);default:''" json:"FLOATING_IP"`
	FloatingIPPort int    `gorm:"column:floating_ip_port;type:int;default:null" json:"FLOATING_IP_PORT"`
	FixedIP        string `gorm:"column:fixed_ip;type:char(64);default:''" json:"FIXED_IP"`
	FixedIPPort    int    `gorm:"column:fixed_ip_port;type:int;default:null" json:"FIXED_IP_PORT"`
	VInterfaceID   int    `gorm:"column:port_id;type:int;default:null" json:"PORT_ID"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

func (NATRule) TableName() string {
	return "nat_rule"
}

type NATVMConnection struct {
	Base         `gorm:"embedded" json:",squash"`
	NATGatewayID int    `gorm:"column:nat_id;type:int;default:null" json:"NAT_ID"`
	VMID         int    `gorm:"column:vm_id;type:int;default:null" json:"VM_ID"`
	Domain       string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

func (NATVMConnection) TableName() string {
	return "nat_vm_connection"
}

type LB struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	Model          int    `gorm:"column:model;type:int;default:0" json:"MODEL"` // 1.Internal 2.External
	VIP            string `gorm:"column:vip;type:text;default:''" json:"VIP"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"EPC_ID"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"` // TODO delete in future
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	UID            string `gorm:"column:uid;type:char(64);default:''" json:"UID"`
}

func (LB) TableName() string {
	return "lb"
}

type LBListener struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	LBID           int    `gorm:"column:lb_id;type:int;default:0" json:"LB_ID"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	IPs            string `gorm:"column:ips;type:text;default:''" json:"IPS"`           // separated by ,
	SNATIPs        string `gorm:"column:snat_ips;type:text;default:''" json:"SNAT_IPS"` // separated by ,
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	Port           int    `gorm:"column:port;type:int;default:null" json:"PORT"`
	Protocol       string `gorm:"column:protocol;type:char(64);default:''" json:"PROTOCOL"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

func (LBListener) TableName() string {
	return "lb_listener"
}

type LBTargetServer struct {
	Base         `gorm:"embedded" json:",squash"`
	LBID         int    `gorm:"column:lb_id;type:int;default:0" json:"LB_ID"`
	LBListenerID int    `gorm:"column:lb_listener_id;type:int;default:0" json:"LB_LISTENER_ID"`
	VPCID        int    `gorm:"column:epc_id;type:int;default:0" json:"EPC_ID"`
	Type         int    `gorm:"column:type;type:int;default:0" json:"TYPE"` // 1.VM 2.IP
	IP           string `gorm:"column:ip;type:char(64);default:''" json:"IP"`
	VMID         int    `gorm:"column:vm_id;type:int;default:0" json:"VM_ID"`
	Port         int    `gorm:"column:port;type:int;default:null" json:"PORT"`
	Protocol     string `gorm:"column:protocol;type:char(64);default:''" json:"PROTOCOL"`
	Domain       string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

func (LBTargetServer) TableName() string {
	return "lb_target_server"
}

type LBVMConnection struct {
	Base   `gorm:"embedded" json:",squash"`
	LBID   int    `gorm:"column:lb_id;type:int;default:null" json:"LB_ID"`
	VMID   int    `gorm:"column:vm_id;type:int;default:null" json:"VM_ID"`
	Domain string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

func (LBVMConnection) TableName() string {
	return "lb_vm_connection"
}

type PeerConnection struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	LocalVPCID     int    `gorm:"column:local_epc_id;type:int;default:0" json:"LOCAL_EPC_ID"`
	RemoteVPCID    int    `gorm:"column:remote_epc_id;type:int;default:0" json:"REMOTE_EPC_ID"`
	LocalRegionID  int    `gorm:"column:local_region_id;type:int;default:0" json:"LOCAL_REGION_ID"`
	RemoteRegionID int    `gorm:"column:remote_region_id;type:int;default:0" json:"REMOTE_REGION_ID"`
	CreateMethod   int    `gorm:"column:create_method;type:int;default:0" json:"CREATE_METHOD"` // 0.learning 1.user_defined
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

type CEN struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	VPCIDs         string `gorm:"column:epc_ids;type:text;default:''" json:"EPC_IDS"` // separated by ,
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

func (CEN) TableName() string {
	return "cen"
}

type RDSInstance struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	State          int    `gorm:"column:state;type:tinyint(1);not null;default:0" json:"STATE"` // 0. Unknown 1. Running 2. Recovering
	Domain         string `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"EPC_ID"`
	Type           int    `gorm:"column:type;type:int;default:0" json:"TYPE"` // 0. Unknown 1. MySQL 2. SqlServer 3. PPAS 4. PostgreSQL 5. MariaDB
	Version        string `gorm:"column:version;type:char(64);default:''" json:"VERSION"`
	Series         int    `gorm:"column:series;type:tinyint(1);not null;default:0" json:"SERIES"` // 0. Unknown 1. basic 2. HA
	Model          int    `gorm:"column:model;type:tinyint(1);not null;default:0" json:"MODEL"`   // 0. Unknown 1. Primary 2. Readonly 3. Temporary 4. Disaster recovery 5. share
	UID            string `gorm:"column:uid;type:char(64);default:''" json:"UID"`
}

func (RDSInstance) TableName() string {
	return "rds_instance"
}

type RedisInstance struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label          string `gorm:"column:label;type:char(64);default:''" json:"LABEL"`
	State          int    `gorm:"column:state;type:tinyint(1);not null;default:0" json:"STATE"` // 0. Unknown 1. Running 2. Recovering
	Domain         string `gorm:"column:domain;type:char(64);default:''" json:"DOMAIN"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:0" json:"EPC_ID"`
	Version        string `gorm:"column:version;type:char(64);default:''" json:"VERSION"`
	InternalHost   string `gorm:"column:internal_host;type:varchar(128);default:''" json:"INTERNAL_HOST"`
	PublicHost     string `gorm:"column:public_host;type:varchar(128);default:''" json:"PUBLIC_HOST"`
	UID            string `gorm:"column:uid;type:char(64);default:''" json:"UID"`
}

type VIP struct {
	Base   `gorm:"embedded" json:",squash"`
	IP     string `gorm:"column:ip;type:char(64);default:''" json:"IP"`
	Domain string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	VTapID uint32 `gorm:"column:vtap_id;type:int;not null;default:0" json:"VTAP_ID"`
}

func (VIP) TableName() string {
	return "vip"
}

type PodCluster struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	ClusterName    string `gorm:"column:cluster_name;type:varchar(256);default:''" json:"CLUSTER_NAME"`
	Version        string `gorm:"column:version;type:varchar(256);default:''" json:"VERSION"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:null" json:"VPC_ID"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	SubDomain      string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

type PodNamespace struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	CloudTags      string `gorm:"column:cloud_tags;type:text;default:''" json:"CLOUD_TAGS"` // separated by ,
	PodClusterID   int    `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	SubDomain      string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

type PodNode struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Type           int    `gorm:"column:type;type:int;default:null" json:"TYPE"`               // 1: Master 2: Node
	ServerType     int    `gorm:"column:server_type;type:int;default:null" json:"SERVER_TYPE"` // 1: Host 2: VM
	State          int    `gorm:"column:state;type:int;default:1" json:"STATE"`                // 0: Exception 1: Normal
	IP             string `gorm:"column:ip;type:char(64);default:''" json:"IP"`
	VCPUNum        int    `gorm:"column:vcpu_num;type:int;default:0" json:"VCPU_NUM"`
	MemTotal       int    `gorm:"column:mem_total;type:int;default:0" json:"MEM_TOTAL"` // unit: M
	PodClusterID   int    `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	VPCID          int    `gorm:"column:epc_id;type:int;default:null" json:"VPC_ID"`
	SubDomain      string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

type PodIngress struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	PodNamespaceID int    `gorm:"column:pod_namespace_id;type:int;default:null" json:"POD_NAMESPACE_ID"`
	PodClusterID   int    `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	SubDomain      string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

type PodIngressRule struct {
	Base         `gorm:"embedded" json:",squash"`
	Name         string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Protocol     string `gorm:"column:protocol;type:char(64);default:''" json:"PROTOCOL"`
	Host         string `gorm:"column:host;type:text;default:''" json:"HOST"`
	PodIngressID int    `gorm:"column:pod_ingress_id;type:int;default:null" json:"POD_INGRESS_ID"`
	SubDomain    string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
}

type PodIngressRuleBackend struct {
	Base             `gorm:"embedded" json:",squash"`
	Path             string `gorm:"column:path;type:text;default:''" json:"PATH"`
	Port             int    `gorm:"column:port;type:int;default:null" json:"PORT"`
	PodServiceID     int    `gorm:"column:pod_service_id;type:int;default:null" json:"POD_SERVICE_ID"`
	PodIngressRuleID int    `gorm:"column:pod_ingress_rule_id;type:int;default:null" json:"POD_INGRESS_RULE_ID"`
	PodIngressID     int    `gorm:"column:pod_ingress_id;type:int;default:null" json:"POD_INGRESS_ID"`
	SubDomain        string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
}

type PodService struct {
	Base             `gorm:"embedded" json:",squash"`
	SoftDeleteBase   `gorm:"embedded" json:",squash"`
	Name             string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Label            string `gorm:"column:label;type:text;default:''" json:"LABEL"`           // separated by ,
	Annotation       string `gorm:"column:annotation;type:text;default:''" json:"ANNOTATION"` // separated by ,
	Alias            string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Type             int    `gorm:"column:type;type:int;default:null" json:"TYPE"`        // 1: ClusterIP 2: NodePort
	Selector         string `gorm:"column:selector;type:text;default:''" json:"SELECTOR"` // separated by ,
	ServiceClusterIP string `gorm:"column:service_cluster_ip;type:char(64);default:''" json:"SERVICE_CLUSTER_IP"`
	PodIngressID     int    `gorm:"column:pod_ingress_id;type:int;default:null" json:"POD_INGRESS_ID"`
	PodNamespaceID   int    `gorm:"column:pod_namespace_id;type:int;default:null" json:"POD_NAMESPACE_ID"`
	PodClusterID     int    `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	VPCID            int    `gorm:"column:epc_id;type:int;default:null" json:"VPC_ID"`
	AZ               string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region           string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	SubDomain        string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain           string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

type PodServicePort struct {
	Base         `gorm:"embedded" json:",squash"`
	Name         string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Protocol     string `gorm:"column:protocol;type:char(64);default:''" json:"PROTOCOL"`
	Port         int    `gorm:"column:port;type:int;default:null" json:"PORT"`
	TargetPort   int    `gorm:"column:target_port;type:int;default:null" json:"TARGET_PORT"`
	NodePort     int    `gorm:"column:node_port;type:int;default:null" json:"NODE_PORT"`
	PodServiceID int    `gorm:"column:pod_service_id;type:int;default:null" json:"POD_SERVICE_ID"`
	SubDomain    string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
}

type PodGroup struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Type           int    `gorm:"column:type;type:int;default:null" json:"TYPE"` // 1: Deployment 2: StatefulSet 3: ReplicationController
	PodNum         int    `gorm:"column:pod_num;type:int;default:1" json:"POD_NUM"`
	Label          string `gorm:"column:label;type:text;default:''" json:"LABEL"` // separated by ,
	PodNamespaceID int    `gorm:"column:pod_namespace_id;type:int;default:null" json:"POD_NAMESPACE_ID"`
	PodClusterID   int    `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	SubDomain      string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

type PodGroupPort struct {
	Base         `gorm:"embedded" json:",squash"`
	Name         string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Protocol     string `gorm:"column:protocol;type:char(64);default:''" json:"PROTOCOL"`
	Port         int    `gorm:"column:port;type:int;default:null" json:"PORT"`
	PodGroupID   int    `gorm:"column:pod_group_id;type:int;default:null" json:"POD_GROUP_ID"`
	PodServiceID int    `gorm:"column:pod_service_id;type:int;default:null" json:"POD_SERVICE_ID"`
	SubDomain    string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
}

type PodReplicaSet struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Name           string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Alias          string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	Label          string `gorm:"column:label;type:text;default:''" json:"LABEL"` // separated by ,
	PodNum         int    `gorm:"column:pod_num;type:int;default:1" json:"POD_NUM"`
	PodGroupID     int    `gorm:"column:pod_group_id;type:int;default:null" json:"POD_GROUP_ID"`
	PodNamespaceID int    `gorm:"column:pod_namespace_id;type:int;default:null" json:"POD_NAMESPACE_ID"`
	PodClusterID   int    `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	AZ             string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region         string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	SubDomain      string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}

func (PodReplicaSet) TableName() string {
	return "pod_rs"
}

type PrometheusTarget struct {
	Base           `gorm:"embedded" json:",squash"`
	SoftDeleteBase `gorm:"embedded" json:",squash"`
	Instance       string `gorm:"column:instance;type:varchar(255);default:''" json:"INSTANCE"`
	Job            string `gorm:"column:job;type:varchar(255);default:''" json:"JOB"`
	ScrapeURL      string `gorm:"column:scrape_url;type:varchar(2083);default:''" json:"SCRAPE_URL"`
	OtherLabels    string `gorm:"column:other_labels;type:text;default:''" json:"OTHER_LABELS"` // separated by ,
	SubDomain      string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain         string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
	PodClusterID   int    `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	CreateMethod   int    `gorm:"column:create_method;type:tinyint(1);default:1" json:"CREATE_METHOD"`
}

func (PrometheusTarget) TableName() string {
	return "prometheus_target"
}

type Pod struct {
	Base            `gorm:"embedded" json:",squash"`
	SoftDeleteBase  `gorm:"embedded" json:",squash"`
	Name            string `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Alias           string `gorm:"column:alias;type:char(64);default:''" json:"ALIAS"`
	State           int    `gorm:"column:state;type:int;not null" json:"STATE"`                    // 0.Exception 1.Running
	Label           string `gorm:"column:label;type:text;default:''" json:"LABEL"`                 // separated by ,
	Annotation      string `gorm:"column:annotation;type:text;default:''" json:"ANNOTATION"`       // separated by ,
	ENV             string `gorm:"column:env;type:text;default:''" json:"ENV"`                     // separated by ,
	ContainerIDs    string `gorm:"column:container_ids;type:text;default:''" json:"CONTAINER_IDS"` // separated by ,
	PodReplicaSetID int    `gorm:"column:pod_rs_id;type:int;default:null" json:"POD_RS_ID"`
	PodGroupID      int    `gorm:"column:pod_group_id;type:int;default:null" json:"POD_GROUP_ID"`
	PodNamespaceID  int    `gorm:"column:pod_namespace_id;type:int;default:null" json:"POD_NAMESPACE_ID"`
	PodNodeID       int    `gorm:"column:pod_node_id;type:int;default:null" json:"POD_NODE_ID"`
	PodClusterID    int    `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	VPCID           int    `gorm:"column:epc_id;type:int;default:null" json:"VPC_ID"`
	AZ              string `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region          string `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	SubDomain       string `gorm:"column:sub_domain;type:char(64);default:''" json:"SUB_DOMAIN"`
	Domain          string `gorm:"column:domain;type:char(64);not null" json:"DOMAIN"`
}
