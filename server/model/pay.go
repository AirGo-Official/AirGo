package model

import "time"

type PreCreatePayToFrontend struct {
	AlipayInfo AlipayPreCreatePayToFrontend `json:"alipay_info"`
	EpayInfo   EpayPreCreatePayToFrontend   `json:"epay_info"`
}

type Pay struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id"       gorm:"primaryKey"`

	Name       string `json:"name"`         //别名
	PayType    string `json:"pay_type"`     //支付类型：alipay epay balance
	PayLogoUrl string `json:"pay_logo_url"` //logo url
	Status     bool   `json:"status"`       //true:启用

	AliPay AliPay `json:"alipay" gorm:"embedded"`
	Epay   Epay   `json:"epay" gorm:"embedded"`
}

type AliPay struct {
	AlipayAppID         string `json:"alipay_app_id"`
	AlipayNotifyURL     string `json:"alipay_notify_url"`
	AlipayAppPrivateKey string `json:"alipay_app_private_key" gorm:"type:text"`
	AlipayAliPublicKey  string `json:"alipay_ali_public_key"  gorm:"type:text"`
	AlipayEncryptKey    string `json:"alipay_encrypt_key"` //alipay 接口内容加密密钥
}

type AlipayPreCreatePayToFrontend struct {
	QRCode string `json:"qr_code"`
}

type Epay struct {
	EpayPid       int64  `json:"epay_pid"`        //商户ID
	EpayKey       string `json:"epay_key"`        //商户密钥
	EpayApiURL    string `json:"epay_api_url"`    //api地址
	EpayReturnURL string `json:"epay_return_url"` //页面跳转通知地址
	EpayNotifyURL string `json:"epay_notify_url"` //异步通知地址
	EpayType      string `json:"epay_type"`       //支付方式, alipay	支付宝 wxpay	微信支付 qqpay	QQ钱包 bank	网银支付
}

// 易支付支付预创建
type EpayPreCreatePay struct {
	Pid        int64  `json:"pid"`          //商户ID
	Type       string `json:"type"`         //支付方式
	OutTradeNo string `json:"out_trade_no"` //商户订单号
	NotifyUrl  string `json:"notify_url"`   //服务器异步通知地址
	ReturnUrl  string `json:"return_url"`   //页面跳转通知地址
	Name       string `json:"name"`         //商品名称,如超过127个字节会自动截取
	Money      string `json:"money"`        //商品金额,如：1.00	单位：元，最大2位小数
	Sign       string `json:"sign"`         //签名字符串，所有参数按照参数名ASCII码从小到大排序（a-z），sign、sign_type、和空值不参与签名！sign = md5 ( a=b&c=d&e=f + KEY ) （注意：+ 为各语言的拼接符，不是字符！），md5结果为小写。
	SignType   string `json:"sign_type"`    //签名类型
}

// 易支付支付预创建返回给前端
type EpayPreCreatePayToFrontend struct {
	EpayApiURL string `json:"epay_api_url"`

	EpayPreCreatePay EpayPreCreatePay `json:"epay_pre_create_pay"`
}

// 易支付支付预创建响应
type EpayPreCreatePayResponse struct {
	Code      int64  `json:"code"`      //返回状态码，1为成功，其它值为失败
	Msg       string `json:"msg"`       //返回信息，失败时返回原因
	TradeNo   string `json:"trade_no"`  //订单号
	Payurl    string `json:"payurl"`    //支付跳转url
	QRCode    string `json:"qrcode"`    //二维码链接
	Urlscheme string `json:"urlscheme"` //小程序跳转url
}

// 易支付支付结果通知
// 通知类型：服务器异步通知（notify_url）、页面跳转通知（return_url）
type EpayResultResponse struct {
	Pid         int64  `json:"pid"          form:"pid"`          //商户ID
	TradeNo     string `json:"trade_no"     form:"trade_no"`     //易支付订单号
	OutTradeNo  string `json:"out_trade_no" form:"out_trade_no"` //商户订单号
	Type        string `json:"type"         form:"type"`         //支付方式
	Name        string `json:"name"         form:"name"`         //商品名称
	Money       string `json:"money"        form:"money"`        //商品金额
	TradeStatus string `json:"trade_status" form:"trade_status"` //支付状态
	Param       string `json:"param"        form:"param"`        //业务扩展参数
	Sign        string `json:"sign"         form:"sign"`         //签名字符串，所有参数按照参数名ASCII码从小到大排序（a-z），sign、sign_type、和空值不参与签名！sign = md5 ( a=b&c=d&e=f + KEY ) （注意：+ 为各语言的拼接符，不是字符！），md5结果为小写。
	SignType    string `json:"sign_type"    form:"sign_type"`    //签名类型
}
