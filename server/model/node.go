package model

import (
	"time"
)

type Node struct {
	ID        int64     `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//基础参数
	Remarks        string  `json:"remarks"     gorm:"comment:别名"`
	Address        string  `json:"address"     gorm:"comment:地址"`
	Port           int64   `json:"port"        gorm:"comment:端口"`
	NodeOrder      int64   `json:"node_order"  gorm:"comment:节点排序"`
	Enabled        bool    `json:"enabled"     gorm:"comment:是否为激活节点"`
	NodeSpeedLimit int64   `json:"node_speed_limit" gorm:"comment:节点限速,Mbps"`
	TrafficRate    float64 `json:"traffic_rate"    gorm:"comment:倍率"`
	NodeType       string  `json:"node_type"       gorm:"comment:节点类型 normal transfer shared"`

	//一些协议参数
	Protocol        string `json:"protocol"   gorm:"comment:协议类型：vless,vmess,trojan,shadowsocks,hysteria2等"`
	V               string `json:"v"          gorm:"default:2;comment:v"`
	Scy             string `json:"scy"        gorm:"comment:加密方式 none,auto,chacha20-poly1305,aes-128-gcm,aes-256-gcm,2022-blake3-aes-128-gcm,2022-blake3-aes-256-gcm,2022-blake3-chacha20-poly1305"`
	ServerKey       string `json:"server_key" gorm:"comment:ss协议:server_key"`
	Aid             int64  `json:"aid"        gorm:"default:0;comment:vmess 额外ID"`
	VlessFlow       string `json:"flow"       gorm:"comment:流控 null,xtls-rprx-vision,xtls-rprx-vision-udp443"`
	VlessEncryption string `json:"encryption" gorm:"default:none;comment:加密方式 vless:none"`
	Network         string `json:"network"    gorm:"default:ws;comment:传输协议 tcp,kcp,ws,h2,quic,grpc"`
	Type            string `json:"type"       gorm:"comment:伪装类型 ws,h2：无;    tcp,kcp：none，http;    mKCP,quic：none，srtp，utp，wechat-video，dtls，wireguard"`
	Host            string `json:"host"       gorm:"comment:伪装域名"`
	Path            string `json:"path"       gorm:"default:/;comment:path"`
	GrpcMode        string `json:"mode"       gorm:"default:multi;comment:grpc传输模式 gun，multi"`
	ServiceName     string `json:"service_name" gorm:"default:service_name;comment:gRPC 的 ServiceName"`
	Security        string `json:"security"     gorm:"comment:传输层安全类型 none,tls,reality"`
	Sni             string `json:"sni"          gorm:"comment:sni"`
	Fingerprint     string `json:"fp"           gorm:"comment:fp"`
	Alpn            string `json:"alpn"         gorm:"comment:alpn"`
	AllowInsecure   bool   `json:"allowInsecure" gorm:"default:true;comment:allowInsecure"`
	Dest            string `json:"dest"          gorm:"comment:dest"`
	PrivateKey      string `json:"private_key"   gorm:"comment:private_key"`
	PublicKey       string `json:"pbk"           gorm:"comment:pbk"`
	ShortId         string `json:"sid"           gorm:"comment:sid"`
	SpiderX         string `json:"spx"           gorm:"comment:spx"`
	HyPorts         string `json:"hy_ports"      gorm:"comment:hy2端口跳跃"`
	HyUpMbps        int64  `json:"hy_up_mbps" gorm:"comment:Hysteria2服务器上行最大速率(Mbps)"`
	HyDownMbps      int64  `json:"hy_down_mbps" gorm:"comment:Hysteria2服务器下行最大速率(Mbps)"`
	HyObfs          string `json:"hy_obfs" gorm:"comment:Hysteria2协议混淆:salamander"`
	HyObfsPassword  string `json:"hy_obfs_password" gorm:"comment:Hysteria2协议混淆密码"`

	//中转参数
	TransferAddress string `json:"transfer_address" gorm:"comment:中转ip"`
	TransferPort    int64  `json:"transfer_port"    gorm:"comment:中转port"`
	TransferNodeID  int64  `json:"transfer_node_id" gorm:"comment:中转绑定的节点ID"`
	//已用上行/已用下行，统计节点流量时使用
	TotalUp   int64 `json:"total_up"        gorm:"-"` //Byte
	TotalDown int64 `json:"total_down"      gorm:"-"` //Byte
	//关联参数
	Goods       []Goods          `json:"goods"         gorm:"many2many:goods_and_nodes"`       //多对多,关联商品
	TrafficLogs []NodeTrafficLog `json:"-"             gorm:"foreignKey:NodeID;references:ID"` //has many
	Access      []Access         `json:"access"        gorm:"many2many:node_and_access"`       //访问控制
	//共享节点需要的uuid;订阅下发是实际的用户uuid
	UUID string `json:"uuid"            gorm:"comment:UUID"`
}

// 节点状态
type NodeStatus struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Status      bool      `json:"status"`
	LastTime    time.Time `json:"last_time"`
	UserAmount  int64     `json:"user_amount"`
	TrafficRate float64   `json:"traffic_rate"` //节点倍率
	U           float64   `json:"u"`            //Byte
	D           float64   `json:"d"`            //Byte
	CPU         float64   `json:"cpu"`
	Mem         float64   `json:"mem"`
	Disk        float64   `json:"disk"`
	Uptime      int64     `json:"uptime"`
}

// 新建共享节点
type NodeSharedReq struct {
	Url string `json:"url"`
}

// 修改混淆
type ChangeHostRequest struct {
	Host string `json:"host"`
}

// node access 多对多
type NodeAndAccess struct {
	NodeID   int64
	AccessID int64
}
