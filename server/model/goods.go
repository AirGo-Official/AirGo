package model

import (
	"time"
)

type Goods struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id" gorm:"primaryKey"`
	//基础参数
	GoodsOrder  int64  `json:"goods_order"   gorm:"comment:排序"`
	CoverImage  string `json:"cover_image"   gorm:"comment:封面图片"`
	Subject     string `json:"subject"       gorm:"comment:商品标题"`
	Des         string `json:"des"           gorm:"comment:描述;type:text"`
	Price       string `json:"price"         gorm:"comment:金额，单位为元，精确到小数点后两位"` //todo 前端修改
	IsShow      bool   `json:"is_show"       gorm:"comment:是否显示"`
	IsSale      bool   `json:"is_sale"       gorm:"comment:是否售卖，是否上架"`
	IsRenew     bool   `json:"is_renew"      gorm:"comment:是否可续费"`
	Quota       int64  `json:"quota"         gorm:"comment:限购数量"`
	Stock       int64  `json:"stock"         gorm:"comment:库存"`
	GoodsType   string `json:"goods_type"    gorm:"comment:类型，general=普通商品 subscribe=订阅 recharge=充值"`
	DeliverType string `json:"deliver_type"  gorm:"comment:发货类型，none=不发货，auto=自动发货，manual=手动发货"`
	DeliverText string `json:"deliver_text"  gorm:"comment:发货内容;type:text"`
	//订阅参数
	Price3Month            string `json:"price_3_month"`
	Price6Month            string `json:"price_6_month"`
	Price12Month           string `json:"price_12_month"`
	PriceUnlimitedDuration string `json:"price_unlimited_duration"`
	EnableTrafficReset     bool   `json:"enable_traffic_reset"`
	TotalBandwidth         int64  `json:"total_bandwidth"      gorm:"comment:总流量,单位GB"`
	NodeConnector          int64  `json:"node_connector"       gorm:"comment:可连接客户端数量"`
	NodeSpeedLimit         int64  `json:"node_speed_limit"     gorm:"default:0;comment:限速Mbps（Mbps）"`
	//充值参数
	RechargeAmount string `json:"recharge_amount" gorm:"comment:充值金额"`
	//关联
	Nodes  []Node   `json:"nodes"         gorm:"many2many:goods_and_nodes"`
	Coupon []Coupon `json:"coupon"        gorm:"many2many:goods_and_coupon"`
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
