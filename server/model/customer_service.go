package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// 用户当前生效的服务，例如 订阅、订阅聚合 这种长期、持续性的服务
type CustomerService struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id"             gorm:"primaryKey"`
	// 关联用户
	UserID   int64  `json:"user_id"      gorm:"用户ID"`
	UserName string `json:"user_name"    gorm:"comment:用户名"`
	//
	ServiceStatus  bool       `json:"service_status"   gorm:"comment:服务是否有效"`
	ServiceStartAt *time.Time `json:"service_start_at" gorm:"comment:服务开始时间;default:null"`
	ServiceEndAt   *time.Time `json:"service_end_at"   gorm:"comment:服务结束时间;default:null"`
	// 续费参数
	IsRenew       bool   `json:"is_renew"         gorm:"comment:是否可续费"`
	RenewalAmount string `json:"renewal_amount"   gorm:"comment:续费金额"`
	//OriginalAmount string `json:"original_amount"  gorm:"comment:原始金额"`
	// 商品参数-基础参数
	GoodsID   int64  `json:"goods_id"      gorm:"comment:商品ID"`
	Subject   string `json:"subject"       gorm:"comment:商品标题"`
	Des       string `json:"des"           gorm:"comment:描述;type:text"`
	Price     string `json:"price"         gorm:"comment:金额，单位为元，精确到小数点后两位"`
	GoodsType string `json:"goods_type"    gorm:"comment:类型，general=普通商品 subscribe=订阅 recharge=充值"`
	Duration  int64  `json:"duration"      gorm:"comment:购买时长(单位：月)"`
	// 商品参数-订阅
	TotalBandwidth  int64     `json:"total_bandwidth"   gorm:"comment:总流量(Byte)"`
	NodeConnector   int64     `json:"node_connector"    gorm:"comment:可连接客户端数量"`
	NodeSpeedLimit  int64     `json:"node_speed_limit"  gorm:"comment:限速Mbps（Mbps）"`
	TrafficResetDay int64     `json:"traffic_reset_day" gorm:"comment:流量重置日"`
	SubStatus       bool      `json:"sub_status"        gorm:"comment:订阅状态"`
	SubUUID         uuid.UUID `json:"sub_uuid"          gorm:"comment:订阅UUID"`
	UsedUp          int64     `json:"used_up"           gorm:"comment:已用上行流量(Byte)"`
	UsedDown        int64     `json:"used_down"         gorm:"comment:已用下行流量(Byte)"`
}
type PushCustomerServiceRequest struct {
	CustomerServiceID int64  `json:"customer_service_id"`
	ToUserName        string `json:"to_user_name"`
}
