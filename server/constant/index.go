package constant

const (
	//工单状态
	TICKET_PROCESSING = "TICKET_PROCESSING" //工单进行中
	TICKET_CLOSED     = "TICKET_CLOSED"     //工单已关闭

	// gin.Context
	CTX_SET_USERNAME = "UserName"
	CTX_SET_USERID   = "UserID"

	//节点类型
	NODE_TYPE_VMESS       = "vmess"
	NODE_TYPE_VLESS       = "vless"
	NODE_TYPE_TROJAN      = "trojan"
	NODE_TYPE_HYSTERIA    = "hysteria"
	NODE_TYPE_SHADOWSOCKS = "shadowsocks"
	NODE_TYPE_TRANSFER    = "transfer"

	//传输协议
	NETWORK_WS   = "ws"
	NETWORK_TCP  = "tcp"
	NETWORK_KCP  = "kcp"
	NETWORK_GRPC = "grpc"
	NETWORK_QUIC = "quic"

	//邮件类型
	EMAIL_TYPE_USER_REGISTER = "EMAIL_TYPE_USER_REGISTER"
	EMAIL_TYPE_USER_RESETPWD = "EMAIL_TYPE_USER_RESETPWD"

	// 商品
	GOODS_TYPE_GENERAL   = "general"
	GOODS_TYPE_SUBSCRIBE = "subscribe"
	GOODS_TYPE_RECHARGE  = "recharge"
	// 发货
	DELIVER_TYPE_NONE   = "none"
	DELIVER_TYPE_AUTO   = "auto"
	DELIVER_TYPE_MANUAL = "manual"

	// 订单类型
	ORDER_TYPE_NEW     = "New"     //新购订单
	ORDER_TYPE_RENEW   = "Renew"   //续费订单
	ORDER_TYPE_RESTORE = "Restore" //恢复库存

	// 订单状态
	ORDER_STATUS_CREATED        = "CREATED"
	ORDER_STATUS_WAIT_BUYER_PAY = "WAIT_BUYER_PAY"
	ORDER_STATUS_TRADE_CLOSED   = "TRADE_CLOSED"
	ORDER_STATUS_TRADE_SUCCESS  = "TRADE_SUCCESS"
	ORDER_STATUS_TRADE_FINISHED = "TRADE_FINISHED"
	ORDER_STATUS_COMPLETED      = "COMPLETED "

	// 支付类型
	PAY_TYPE_ALIPAY  = "alipay"  // 支付宝alipay
	PAY_TYPE_EPAY    = "epay"    //易支付
	PAY_TYPE_BALANCE = "balance" //余额支付
)
