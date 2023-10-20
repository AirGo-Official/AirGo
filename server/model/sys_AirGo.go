package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type AGNodeStatus struct {
	ID int64 `json:"id"`
	AGNodeStatusItem
}

// NodeStatus Node status
type AGNodeStatusItem struct {
	CPU    float64
	Mem    float64
	Disk   float64
	Uptime uint64
}

type AGNodeInfo struct {
	ID             int64  `json:"id"`
	NodeSpeedlimit int64  `json:"node_speedlimit"` //节点限速/Mbps
	TrafficRate    int64  `json:"traffic_rate"`    //倍率
	NodeType       string `json:"node_type"`       //节点类型 vless,vmess,trojan
	Remarks        string `json:"remarks"`         //别名
	Address        string `json:"address"`         //地址
	Port           int64  `json:"port"`            //端口

	//vmess参数
	Scy       string `json:"scy"`
	ServerKey string `json:"server_key"`
	Aid       int64  `json:"aid"`
	//vless参数
	VlessFlow string `json:"flow"` //流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443

	//传输参数
	Network     string `json:"network"`      //传输协议 tcp,kcp,ws,h2,quic,grpc
	Type        string `json:"type"`         //伪装类型 ws,h2：无    tcp,kcp：none，http    quic：none，srtp，utp，wechat-video，dtls，wireguard
	Host        string `json:"host"`         //伪装域名
	Path        string `json:"path"`         //path
	GrpcMode    string `json:"mode"`         //grpc传输模式 gun，multi
	ServiceName string `json:"service_name"` //

	//传输层安全
	Security    string `json:"security"` //传输层安全类型 none,tls,reality
	Sni         string `json:"sni"`      //
	Fingerprint string `json:"fp"`       //
	Alpn        string `json:"alpn"`     //
	Dest        string `json:"dest"`
	PrivateKey  string `json:"private_key"`
	PublicKey   string `json:"pbk"`
	ShortId     string `json:"sid"`
	SpiderX     string `json:"spx"`
}
type AGUserInfo struct {
	ID       int64     `json:"id"`
	UUID     uuid.UUID `json:"uuid"`
	Passwd   string    `json:"passwd"`
	UserName string    `json:"user_name"`
}

type AGUserTraffic struct {
	ID          int64               `json:"id"`
	UserTraffic []AGUserTrafficItem `json:"user_traffic"`
}

type AGUserTrafficItem struct {
	UID      int64
	Email    string
	Upload   int64
	Download int64
}

type AGREALITYConfig struct {
	Dest             string
	ProxyProtocolVer uint64
	ServerNames      []string
	PrivateKey       string
	MinClientVer     string
	MaxClientVer     string
	MaxTimeDiff      uint64
	ShortIds         []string
}

type AGREALITYx25519 struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// 数据库 traffic_log 流量统计表
type TrafficLog struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"      gorm:"primary_key"`
	NodeID    int64     `json:"node_id" gorm:"comment:节点ID"`
	U         int64     `json:"u"       gorm:"comment:上行流量 bit"`
	D         int64     `json:"d"       gorm:"comment:下行流量 bit"`
}
