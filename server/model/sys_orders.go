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

type Orders struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt *time.Time `json:"-"        gorm:"index"`
	ID       int64  `json:"id"       gorm:"primary_key"`
	UserID   int64  `json:"user_id"   gorm:"comment:外键"` //外键
	UserName string `json:"user_name" gorm:"comment:订单拥有者"`
	User     User   `json:"-"     gorm:"foreignKey:UserID;references:ID"` //订单拥有者

	OutTradeNo      string `json:"out_trade_no" gorm:"comment:商户订单号"`
	GoodsID         int64  `json:"goods_id"     gorm:"comment:商品id"`
	Subject         string `json:"subject"      gorm:"comment:商品的标题/交易标题/订单标题/订单关键字"`
	Price           string `json:"price"        gorm:"comment:商品的价格"`
	PayID           int64  `json:"pay_id"       gorm:"comment:支付方式id"`
	PayType         string `json:"pay_type"     gorm:"comment:支付方式，alipay,epay"`
	CouponID        int64  `json:"coupon_id"        gorm:"comment:折扣码id"`
	CouponName      string `json:"coupon_name"      gorm:"comment:折扣码"`
	CouponAmount    string `json:"coupon_amount"    gorm:"comment:折扣金额"`
	DeductionAmount string `json:"deduction_amount" gorm:"comment:旧套餐折扣金额"`
	RemainAmount    string `json:"remain_amount"    gorm:"comment:余额折扣金额"`
	TradeStatus     string `json:"trade_status"     gorm:"default:Created;comment:交易状态 1、WAIT_BUYER_PAY（交易创建，等待买家付款）；2、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）；3、TRADE_SUCCESS（交易支付成功）； 4、TRADE_FINISHED（交易结束，不可退款）；5、Completed（手动完成订单）；6、Created（订单已创建）"`
	TotalAmount     string `json:"total_amount"     gorm:"comment:订单金额"`
	ReceiptAmount   string `json:"receipt_amount"   gorm:"comment:实收金额"`
	BuyerPayAmount  string `json:"buyer_pay_amount" gorm:"comment:付款金额"`
	//第三方支付参数
	PayInfo      PreCreatePayToFrontend `json:"pay_info"         gorm:"-"` //支付信息，epay，alipay等"
	TradeNo      string                 `json:"trade_no"         gorm:"comment:第三方交易号"`
	BuyerLogonId string                 `json:"buyer_logon_id"   gorm:"comment:买家第三方账号"`
}

// 分页订单数据，带总数
type OrdersWithTotal struct {
	OrderList []Orders `json:"order_list"`
	Total     int64    `json:"total"`
}

// 订单收入统计
type OrderStatistics struct {
	Total       int64   `json:"total"`
	TotalAmount float64 `json:"total_amount"`
}

type OrdersHeader struct {
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	OutTradeNo     string `json:"out_trade_no"`
	GoodsID        string `json:"goods_id"`
	Subject        string `json:"subject"`
	Price          string `json:"price"`
	PayType        string `json:"pay_type"`
	TradeNo        string `json:"trade_no"`
	BuyerLogonId   string `json:"buyer_logon_id"`
	TradeStatus    string `json:"trade_status"`
	TotalAmount    string `json:"total_amount"`
	ReceiptAmount  string `json:"receipt_amount"`
	BuyerPayAmount string `json:"buyer_pay_amount"`
}

var OrdersHeaderItem = OrdersHeader{
	CreatedAt:  "创建日期",
	UpdatedAt:  "更新日期",
	ID:         "ID",
	UserID:     "用户ID",
	UserName:   "用户名",
	OutTradeNo: "系统订单号",
	GoodsID:    "商品ID",
	Subject:    "商品标题",
	//Price:"价格",
	PayType:        "支付类型",
	TradeNo:        "支付宝订单号",
	BuyerLogonId:   "买家支付宝账号",
	TradeStatus:    "交易状态",
	TotalAmount:    "订单金额",
	ReceiptAmount:  "实收金额",
	BuyerPayAmount: "付款金额",
}
