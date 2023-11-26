package model

import (
	"time"
)

type Node struct {
	ID        int64     `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//基础参数
	Remarks         string `json:"remarks"`                                  //别名
	Address         string `json:"address"`                                  //地址
	Port            int64  `json:"port"`                                     //端口
	NodeOrder       int64  `json:"node_order"`                               //节点排序
	Enabled         bool   `json:"enabled"`                                  //是否为激活节点
	NodeSpeedlimit  int64  `json:"node_speedlimit"`                          //节点限速/Mbps
	TrafficRate     int64  `json:"traffic_rate"`                             //倍率
	NodeType        string `json:"node_type"`                                //节点类型 vless,vmess,trojan,shadowsocks
	IsSharedNode    bool   `json:"is_shared_node"`                           //共享节点，不修改uuid和host
	V               string `json:"v"   gorm:"default:2"`                     //
	Scy             string `json:"scy"`                                      //加密方式 none,auto,chacha20-poly1305,aes-128-gcm,aes-256-gcm,2022-blake3-aes-128-gcm,2022-blake3-aes-256-gcm,2022-blake3-chacha20-poly1305
	ServerKey       string `json:"server_key"`                               //
	Aid             int64  `json:"aid"  gorm:"default:0"`                    //额外ID
	VlessFlow       string `json:"flow" gorm:"default:none"`                 //流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443
	VlessEncryption string `json:"encryption" gorm:"default:none"`           //加密方式 vless:none
	Network         string `json:"network" gorm:"default:ws"`                //传输协议 tcp,kcp,ws,h2,quic,grpc
	Type            string `json:"type"    gorm:"default:none"`              //伪装类型 ws,h2：无    tcp,kcp：none，http    mKCP,quic：none，srtp，utp，wechat-video，dtls，wireguard
	Host            string `json:"host"`                                     //伪装域名
	Path            string `json:"path"    gorm:"default:/"`                 //path(ws,h2)
	GrpcMode        string `json:"mode"    gorm:"default:multi"`             //grpc传输模式 gun，multi
	ServiceName     string `json:"service_name" gorm:"default:service_name"` //gRPC 的 ServiceName
	Security        string `json:"security" gorm:"default:none"`             //传输层安全类型 none,tls,reality
	Sni             string `json:"sni"`                                      //
	Fingerprint     string `json:"fp"`                                       //
	Alpn            string `json:"alpn"`                                     //
	AllowInsecure   bool   `json:"allowInsecure" gorm:"default:true"`        //tls 跳过证书验证
	Dest            string `json:"dest"`
	PrivateKey      string `json:"private_key"`
	PublicKey       string `json:"pbk"`
	ShortId         string `json:"sid"`
	SpiderX         string `json:"spx"`

	//共享节点额外需要的参数
	UUID string `json:"uuid"` //用户id
	//中转参数
	EnableTransfer  bool   `json:"enable_transfer" gorm:"default:false"` //是否启用中转
	TransferAddress string `json:"transfer_address"`                     //中转ip
	TransferPort    int64  `json:"transfer_port"`                        //中转port
	//上行/下行
	TotalUp   int64 `json:"total_up"        gorm:"-"` //Byte
	TotalDown int64 `json:"total_down"      gorm:"-"` //Byte
	//关联参数
	Goods       []Goods      `json:"goods"         gorm:"many2many:goods_and_nodes"`       //多对多,关联商品
	TrafficLogs []TrafficLog `json:"-"             gorm:"foreignKey:NodeID;references:ID"` //has many
	Access      []Access     `json:"access"        gorm:"many2many:node_and_access"`       //访问控制
	AccessIds   []int64      `json:"access_ids"    gorm:"-"`
}

// 节点状态
type NodeStatus struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Status      bool      `json:"status"`
	LastTime    time.Time `json:"last_time"`
	UserAmount  int64     `json:"user_amount"`
	TrafficRate int64     `json:"traffic_rate"` //节点倍率
	U           float64   `json:"u"`            //Byte
	D           float64   `json:"d"`            //Byte
	CPU         float64   `json:"cpu"`
	Mem         float64   `json:"mem"`
	Disk        float64   `json:"disk"`
	Uptime      int64     `json:"uptime"`
}

// 共享节点
type NodeShared struct {
	ID              int64     `json:"id" gorm:"primarykey"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Remarks         string    `json:"remarks"`                            //别名
	Address         string    `json:"address"`                            //地址
	Port            int64     `json:"port"`                               //端口
	NodeOrder       int64     `json:"node_order"`                         //节点排序
	Enabled         bool      `json:"enabled" gorm:"default:true"`        //是否为激活节点
	NodeSpeedlimit  int64     `json:"node_speedlimit"`                    //节点限速/Mbps
	TrafficRate     int64     `json:"traffic_rate"`                       //倍率
	NodeType        string    `json:"node_type"`                          //节点类型 vless,vmess,shadowsocks,hysteria
	IsSharedNode    bool      `json:"is_shared_node" gorm:"default:true"` //共享节点，不修改uuid和host
	V               string    `json:"v"`                                  //
	Scy             string    `json:"scy"`                                //加密方式 none,auto,chacha20-poly1305,aes-128-gcm,aes-256-gcm,2022-blake3-aes-128-gcm,2022-blake3-aes-256-gcm,2022-blake3-chacha20-poly1305
	ServerKey       string    `json:"server_key"`                         //
	Aid             int64     `json:"aid"`                                //额外ID
	VlessFlow       string    `json:"flow"`                               //流控 null,xtls-rprx-vision,xtls-rprx-vision-udp443
	VlessEncryption string    `json:"encryption"`                         //加密方式 none
	Network         string    `json:"network"`                            //传输协议 tcp,kcp,ws,h2,quic,grpc
	Type            string    `json:"type"`                               //伪装类型 ws,h2：无    tcp,kcp：none，http    mKCP,quic：none，srtp，utp，wechat-video，dtls，wireguard
	Host            string    `json:"host"`                               //伪装域名
	Path            string    `json:"path"`                               //path(ws,h2)
	GrpcMode        string    `json:"mode"`                               //grpc传输模式 gun，multi
	ServiceName     string    `json:"service_name"`                       //gRPC 的 ServiceName
	Security        string    `json:"security"`                           //传输层安全类型 none,tls,reality
	Sni             string    `json:"sni"`                                //
	Fingerprint     string    `json:"fp"`                                 //
	Alpn            string    `json:"alpn"`                               //
	AllowInsecure   bool      `json:"allowInsecure"`                      //tls 跳过证书验证
	Dest            string    `json:"dest"`
	PrivateKey      string    `json:"private_key"`
	PublicKey       string    `json:"pbk"`
	ShortId         string    `json:"sid"`
	SpiderX         string    `json:"spx"`

	//共享节点额外需要的参数
	UUID string `json:"uuid"`
	//中转参数
	EnableTransfer  bool   `json:"enable_transfer" gorm:"default:false"` //是否启用中转
	TransferAddress string `json:"transfer_address"`                     //中转ip
	TransferPort    int64  `json:"transfer_port"`                        //中转port
	//上行/下行
	TotalUp   int64 `json:"total_up"        gorm:"-"` //Byte
	TotalDown int64 `json:"total_down"      gorm:"-"` //Byte
	//关联参数
	//Goods       []Goods      `json:"goods"         gorm:"many2many:goods_and_nodes"`       //多对多,关联商品
	//TrafficLogs []TrafficLog `json:"-"             gorm:"foreignKey:NodeID;references:ID"` //has many
}

// 新建共享节点
type NodeSharedReq struct {
	Url string `json:"url"`
}

// 修改混淆
type SubHost struct {
	Host string `json:"host"`
}

// 查询节点 with total
type NodesWithTotal struct {
	NodeList []Node `json:"node_list"`
	Total    int64  `json:"total"`
}

// node access 多对多
type NodeAndAccess struct {
	NodeID   int64
	AccessID int64
}
