package constant

const (
	CACHE_GENERAL_DURATION = 3 //数据一致性要求较低的数据缓存的时间，单位s

	CACHE_DEFAULT_ARTICLE = "CACHE_DEFAULT_ARTICLE" //默认文章

	// 订单
	CACHE_SUBMIT_ORDER_TIMEOUT         = 5                               //time.Minute
	CACHE_SUBMIT_ORDER_BY_USERID       = "SUBMIT_ORDER_BY_USERID:"       //用户提交的订单
	CACHE_SUBMIT_ORDER_BY_ORDERID      = "SUBMIT_ORDER_BY_ORDERID:"      //用户提交的订单
	CACHE_USERID_AND_GOODSID           = "USERID_AND_GOODSID:"           //幂等性检查
	CACHE_USERID_AND_CUSTOMERSERVICEID = "USERID_AND_CUSTOMERSERVICEID:" //幂等性检查

	// 商品
	CACHE_GOODS_BY_ID                 = "CACHE_GOODS_BY_ID:"
	CACHE_ALL_ENABLED_GOODS_GENERAL   = "AllEnabledGoodsGeneral"
	CACHE_ALL_ENABLED_GOODS_SUBSCRIBE = "AllEnabledGoodsSubscribe"
	CACHE_ALL_ENABLED_GOODS_RECHARGE  = "AllEnabledGoodsRecharge"

	// 折扣码
	CACHE_COUPON_BY_NAME = "CouponByName:"
	CACHE_COUPON_BY_ID   = "CouponByID:"

	// 角色
	CACHE_USER_ROLEIDS_BY_USERID = "UserRoleIDsByUserID:"
	CACHE_USER_ROLES_BY_USERID   = "UserRolesByUserID:"

	//节点
	CACHE_NODE_STATUS_BY_NODEID             = "NodeStatusByNodeID:"           //节点状态 todo 优化
	CACHE_NODE_STATUS_IS_NOTIFIED_BY_NODEID = "NodeStatusIsNotifiedByNodeID:" // todo 优化

	// 邮箱验证码
	CACHE_USER_REGISTER_EMAIL_CODE_BY_USERNAME  = "UserRegisterEmailCodeByUserName:" //注册验证码
	CACHE_USER_RESET_PWD_EMAIL_CODE_BY_USERNAME = "UserResetPwdEmailCodeByUserName:" //重置密码验证码

	//用户
	CACHE_USER_TOKEN_BY_ID      = "UserTokenByID:"     //用户缓存的token
	CACHE_USER_MENU_LIST_BY_ID  = "UserRouteListByID:" // 用户菜单（路由）
	CACHE_USER_IS_CLOCKIN_BY_ID = "UserIsClockInByID:"
)
