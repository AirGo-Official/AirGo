package model

import (
	"time"
)

const (
	GoodsTypeGeneral   = "general"
	GoodsTypeSubscribe = "subscribe"
	GoodsTypeRecharge  = "recharge"
	DeliverTypeNone    = "none"
	DeliverTypeAuto    = "auto"
	DeliverTypeManual  = "manual"
	//查询
	AllEnabledGoods          = "AllEnabledGoods"
	AllEnabledGoodsGeneral   = "AllEnabledGoodsGeneral"
	AllEnabledGoodsSubscribe = "AllEnabledGoodsSubscribe"
	AllEnabledGoodsRecharge  = "AllEnabledGoodsRecharge"
)

type Goods struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id" gorm:"primary_key"`
	//基础参数
	GoodsOrder  int64  `json:"goods_order"   gorm:"comment:排序"`
	Subject     string `json:"subject"       gorm:"comment:商品标题"`
	Des         string `json:"des"           gorm:"comment:描述;type:text"`
	TotalAmount string `json:"total_amount"  gorm:"comment:金额，单位为元，精确到小数点后两位"`
	Status      bool   `json:"status"        gorm:"default:true;comment:是否显示"`
	GoodsType   string `json:"goods_type"    gorm:"comment:类型，general=普通商品 subscribe=订阅 recharge=充值"`
	DeliverType string `json:"deliver_type"  gorm:"comment:发货类型，none=不发货，auto=自动发货，manual=手动发货"`
	DeliverText string `json:"deliver_text"  gorm:"comment:发货内容;type:text"`
	//关联参数
	CheckedNodes []int64  `json:"checked_nodes" gorm:"-"` //前端套餐编辑时选中的节点
	Nodes        []Node   `json:"nodes"         gorm:"many2many:goods_and_nodes"`
	Coupon       []Coupon `json:"coupon"        gorm:"many2many:goods_and_coupon"`
	//订阅参数
	TotalBandwidth     int64  `json:"total_bandwidth"      gorm:"comment:总流量"`
	ExpirationDate     int64  `json:"expiration_date"      gorm:"comment:有效期"`
	NodeConnector      int64  `json:"node_connector"       gorm:"comment:可连接客户端数量"`
	TrafficResetMethod string `json:"traffic_reset_method" gorm:"comment:流量重置方式,Stack 叠加,NotStack 不叠加"`
	ResetDay           int64  `json:"reset_day"            gorm:"comment:流量重置日"`
	//充值参数
	RechargeAmount string `json:"recharge_amount" gorm:"comment:充值金额"`
}

// 商品和节点 多对多 表
type GoodsAndNodes struct {
	GoodsID int64 `gorm:"column:goods_id"`
	NodeID  int64 `gorm:"column:node_id"`
}

// 商品和折扣码 多对多 表
type GoodsAndCoupon struct {
	GoodsID  int64 `gorm:"column:goods_id"`
	CouponID int64 `gorm:"column:coupon_id"`
}
