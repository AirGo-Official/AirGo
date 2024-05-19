package constant

const (
	// 限流
	ERROR_IP_LIMIT   = "Current IP requests are frequent."
	ERROR_USER_LIMIT = "Current user requests are frequent."

	//token
	ERROR_NO_TOKEN_IN_THE_REQUEST = "There is no token in the request."
	// 请求参数解析错误
	ERROR_REQUEST_PARAMETER_PARSING_ERROR = "Request parameter parsing error."

	// 商品
	ERROR_INVALID_GOODS_TYPE   = "Invalid goods type" // 无效商品类型
	ERROR_GOODS_NOT_SALE       = "The goods is not for sale."
	ERROR_STOCK_OF_GOODS_EMPTY = "The stock of goods are empty."
	ERROR_GOODS_EXCEEDED_QUOTA = "The product you purchased has exceeded the purchase limit."

	// 订单
	ERROR_DUPLICATE_ORDER      = "Please do not place duplicate orders."
	ERROR_ORDER_TIMED_OUT      = "Order has timed out."
	ERROR_INVALID_ORDER_TYPE   = "Invalid order type."
	ERROR_INVALID_ORDER_PARAMS = "Invalid order params."
	// 支付
	ERROR_INVALID_PAY_TYPE      = "Invalid pay type."      // 无效支付类型
	ERROR_BALANCE_IS_NOT_ENOUGH = "Balance is not enough." // 余额不足

	// 折扣码
	ERROR_COUPON_NOT_EXIST       = "The coupon does not exist."
	ERROR_COUPON_QUERY_ERROR     = "Coupon query error."
	ERROR_COUPON_NOT_APPLICABLE  = "Coupon not applicable."
	ERROR_COUPON_HAS_EXPIRED     = "The coupon has expired."
	ERROR_COUPON_COUNT_EXHAUSTED = "Coupon count exhausted."

	// 用户服务
	ERROR_CUSTOMER_SERVICE_NO_RENEWAL = "Customer service is not allowed to be renewed."

	//节点
	ERROR_INVALID_NODE_TYPE = "Invalid node type."

	//finance
	ERROR_COMMISSION_IS_NOT_ENOUGH = "Insufficient commission"

	//未开启开服务
	ERROR_SERVICE_NOT_ENABLED = "Service not enabled"
)
