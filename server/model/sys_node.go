package model

import "time"

// sspanel 响应 获取当前请求节点的节点设置
type SSNodeInfo struct {
	// NodeGroup      int64    `json:"node_group"`
	// NodeClass      int64    `json:"node_class"`
	//MuOnly         int64    `json:"mu_only"` //只启用单端口多用户
	NodeSpeedlimit int64  `json:"node_speedlimit"` //节点限速/Mbps
	TrafficRate    int64  `json:"traffic_rate"`    //倍率
	Sort           int64  `json:"sort"`            //类型
	Server         string `json:"server"`          //
	SSType         string `json:"type"`            //显示与隐藏
}

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
	Scy             string `json:"scy"`                                      //加密方式 vless需要填 "none"，不能留空。none,auto,chacha20-poly1305,aes-128-gcm,aes-256-gcm,2022-blake3-aes-128-gcm,2022-blake3-aes-256-gcm,2022-blake3-chacha20-poly1305
	ServerKey       string `json:"server_key"`                               //
	Aid             int64  `json:"aid" gorm:"default:0"`                     //额外ID
	VlessFlow       string `json:"flow"`                                     //流控 null,xtls-rprx-vision,xtls-rprx-vision-udp443
	VlessEncryption string `json:"encryption" gorm:"default:none"`           //加密方式 none
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

}

// vmess 格式
type Vmess struct {
	V            string `json:"v" `   //
	Name         string `json:"ps"`   //节点名
	Address      string `json:"add"`  //节点地址
	Port         string `json:"port"` //端口
	Uuid         string `json:"id"`   //用户UUID
	Aid          string `json:"aid"`  //额外ID
	Net          string `json:"net"`  //传输协议
	Disguisetype string `json:"type"` //伪装类型
	Host         string `json:"host"` //伪装域名
	Path         string `json:"path"` //
	Tls          string `json:"tls"`  //传输层安全
	Alpn         string `json:"alpn"`
	Fp           string `json:"fp"`
	Sni          string `json:"sni"`
}

// clash  yaml格式
type ClashYaml struct {
	Port               int64             `yaml:"port"`
	SocksPort          int64             `yaml:"socks-port"`
	RedirPort          int64             `yaml:"redir-port"`
	AllowLan           bool              `yaml:"allow-lan"`
	Mode               string            `yaml:"mode"`
	LogLevel           string            `yaml:"log-level"`
	ExternalController string            `yaml:"external-controller"`
	Secret             string            `yaml:"secret"`
	Proxies            []ClashProxy      `yaml:"proxies"`
	ProxyGroups        []ClashProxyGroup `yaml:"proxy-groups"`
	Rules              []string          `yaml:"rules"`
}
type ClashProxy struct {
	//基础参数
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Server  string `yaml:"server"`
	Port    int    `yaml:"port"`
	Uuid    string `yaml:"uuid"`
	Network string `yaml:"network"`
	Udp     bool   `yaml:"udp"`
	//vmess参数
	Alterid string `yaml:"alterId"` //0
	Cipher  string `yaml:"cipher"`  //auto
	//trojan 参数
	Password string `yaml:"password"`
	//vless流控
	Flow string `yaml:"flow"`

	Tls               bool     `yaml:"tls"`
	Sni               string   `yaml:"sni"`
	ClientFingerprint string   `yaml:"client-fingerprint"` // # Available: "chrome","firefox","safari","ios","random", currently only support TLS transport in TCP/GRPC/WS/HTTP for VLESS/Vmess and trojan.
	Alpn              []string `yaml:"alpn"`               //h2 http/1.1
	Servername        string   `yaml:"servername"`         //REALITY servername
	SkipCertVerify    bool     `yaml:"skip-cert-verify"`

	//WsPath    string    `yaml:"ws-path"`
	//WsHeaders WsHeaders `yaml:"ws-headers"`
	WsOpts      WsOpts      `yaml:"ws-opts"`
	RealityOpts RealityOpts `yaml:"reality-opts"`
	GrpcOpts    GrpcOpts    `yaml:"grpc-opts"`
	H2Opts      H2Opts      `yaml:"h2-opts"`
}

type WsOpts struct {
	Path                string            `yaml:"path"`
	Headers             map[string]string `yaml:"headers"`
	MaxEarlyData        int               `yaml:"max-early-data"`         //2048
	EarlyDataHeaderName string            `yaml:"early-data-header-name"` //Sec-WebSocket-Protocol
}
type WsHeaders struct {
	Host string `yaml:"Host"`
}

type RealityOpts struct {
	PublicKey string `yaml:"public-key"`
	ShortID   string `yaml:"short-id"`
}
type GrpcOpts struct {
	GrpcServiceName string `yaml:"grpc-service-name"` //grpc
}

type H2Opts struct {
	Host []string `yaml:"host"`
	Path string   `yaml:"path"`
}
type HttpOpts struct {
	Method  string                `yaml:"method"` //GET
	Path    map[string][]string   `yaml:"path"`
	Headers map[string]Connection `yaml:"headers"`
}
type Connection []string

type ClashProxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
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
	Scy             string `json:"scy"`                                      //加密方式 vless需要填 "none"，不能留空.none,auto,chacha20-poly1305,aes-128-gcm,aes-256-gcm,2022-blake3-aes-128-gcm,2022-blake3-aes-256-gcm,2022-blake3-chacha20-poly1305
	ServerKey       string `json:"server_key"`                               //
	Aid             int64  `json:"aid" gorm:"default:0"`                     //额外ID
	VlessFlow       string `json:"flow"`                                     //流控 null,xtls-rprx-vision,xtls-rprx-vision-udp443
	VlessEncryption string `json:"encryption" gorm:"default:none"`           //加密方式 none
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
