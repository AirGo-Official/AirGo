package model

import "time"

const (
	OrderCreated        = "Created"        //订单已创建
	OrderCompleted      = "Completed"      //手动完成订单
	OrderWAIT_BUYER_PAY = "WAIT_BUYER_PAY" //交易创建，等待买家付款
	OrderTRADE_CLOSED   = "TRADE_CLOSED"   //未付款交易超时关闭，或支付完成后全额退款
	OrderTRADE_SUCCESS  = "TRADE_SUCCESS"  //交易支付成功
	OrderTRADE_FINISHED = "TRADE_FINISHED" //交易结束，不可退款

)

type Order struct {
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"-"                gorm:"index"`
	ID             int64      `json:"id"               gorm:"primaryKey"`
	OrderType      string     `json:"order_type"       gorm:"comment:订单类型:New=新购入;Renew=续费;Restore=订单补偿"`
	OrderRemarks   string     `json:"order_remarks"    gorm:"comment:订单备注"`
	TradeStatus    string     `json:"trade_status"     gorm:"comment:交易状态 1、WAIT_BUYER_PAY（交易创建，等待买家付款）；2、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）；3、TRADE_SUCCESS（交易支付成功）； 4、TRADE_FINISHED（交易结束，不可退款）；5、COMPLETED（手动完成订单）；6、CREATED（订单已创建）"`
	OutTradeNo     string     `json:"out_trade_no"     gorm:"comment:商户订单号"`
	OriginalAmount string     `json:"original_amount"  gorm:"comment:原始金额"`
	TotalAmount    string     `json:"total_amount"     gorm:"comment:订单金额"`
	BuyerPayAmount string     `json:"buyer_pay_amount" gorm:"comment:付款金额"`
	CouponAmount   string     `json:"coupon_amount"    gorm:"comment:折扣码折扣金额"`
	BalanceAmount  string     `json:"balance_amount"   gorm:"comment:余额折扣金额"`
	// 关联用户
	UserID   int64  `json:"user_id"          gorm:"comment:外键"` //外键
	UserName string `json:"user_name"        gorm:"comment:订单拥有者"`
	User     User   `json:"-"                gorm:"foreignKey:UserID;references:ID"` //订单拥有者
	// 商品参数
	GoodsID     int64  `json:"goods_id"      gorm:"comment:商品id"`
	Des         string `json:"des"           gorm:"comment:描述;type:text"`
	GoodsType   string `json:"goods_type"    gorm:"comment:类型，general=普通商品 subscribe=订阅 recharge=充值"`
	DeliverType string `json:"deliver_type"  gorm:"comment:发货类型，none=不发货，auto=自动发货，manual=手动发货"`
	DeliverText string `json:"deliver_text"  gorm:"comment:发货内容;type:text"`
	Subject     string `json:"subject"       gorm:"comment:商品的标题/交易标题/订单标题/订单关键字"`
	Price       string `json:"price"         gorm:"comment:商品的单价(默认月付价格)"`
	Duration    int64  `json:"duration"      gorm:"comment:购买时长(单位：月)"`
	// 服务参数
	CustomerServiceID int64 `json:"customer_service_id" gorm:"comment:客户服务ID"`
	// 支付参数
	PayID        int64                  `json:"pay_id"           gorm:"comment:支付方式id"`
	PayType      string                 `json:"pay_type"         gorm:"comment:支付方式，alipay,epay"`
	PayInfo      PreCreatePayToFrontend `json:"pay_info"         gorm:"-"` //支付信息，epay，alipay等"
	TradeNo      string                 `json:"trade_no"         gorm:"comment:第三方交易号"`
	BuyerLogonId string                 `json:"buyer_logon_id"   gorm:"comment:买家第三方账号"`
	CouponID     int64                  `json:"coupon_id"        gorm:"comment:折扣码id"`
	CouponName   string                 `json:"coupon_name"      gorm:"comment:折扣码"`
}

type PreCreateOrderRequest struct {
	OrderType     string `json:"order_type"`
	UserID        int64  `json:"user_id"`
	GoodsID       int64  `json:"goods_id"`
	UserServiceID int64  `json:"user_service_id"`
	CouponName    string `json:"coupon_name"`
}

// 订单收入统计
type OrderStatistics struct {
	Total       int64   `json:"total"`
	TotalAmount float64 `json:"total_amount"`
}
