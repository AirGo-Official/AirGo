package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api/user_api"
	middleware "github.com/ppoonk/AirGo/router/middleware"
)

func (g *GinRouter) InitUserRouter(r *gin.RouterGroup) {
	customerRouter := r.Group("/customer")
	customerRouter.Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	//user
	userRouter := customerRouter.Group("/user")
	{
		userRouter.GET("/getUserInfo", user_api.GetUserInfo)                //获取自身信息
		userRouter.POST("/changeUserPassword", user_api.ChangeUserPassword) //修改密码
		userRouter.POST("/changeUserAvatar", user_api.ChangeUserAvatar)     //修改头像
		userRouter.GET("/clockIn", user_api.ClockIn)                        //打卡
		userRouter.POST("/setUserNotice", user_api.SetUserNotice)           //设置通知
	}
	// customer
	customerServiceRouter := customerRouter.Group("/customerService")
	{
		customerServiceRouter.GET("/getCustomerServiceList", user_api.GetCustomerServiceList)
		customerServiceRouter.POST("/resetSubscribeUUID", user_api.ResetSubscribeUUID)
		customerServiceRouter.POST("/pushCustomerService", user_api.PushCustomerService)
		customerServiceRouter.DELETE("/deleteCustomerService", user_api.DeleteCustomerService)
	}
	// menu
	menuRouter := customerRouter.Group("/menu")
	{
		menuRouter.GET("/getMenuList", user_api.GetMenuList) //获取当前角色动态路由
	}
	// shop
	shopRouter := customerRouter.Group("/shop")
	{
		shopRouter.POST("/purchase", user_api.Purchase)                      //支付
		shopRouter.GET("/getEnabledGoodsList", user_api.GetEnabledGoodsList) //查询全部已启用商品
	}
	// order
	orderRouter := customerRouter.Group("/order")
	{
		orderRouter.POST("/preCreateOrder", user_api.PreCreateOrder)           //交易预创建(提交订单)
		orderRouter.POST("/getOrderInfo", user_api.GetOrderInfo)               //获取订单详情(下单时）
		orderRouter.POST("/getOrderList", user_api.GetOrderList)               //获取订单
		orderRouter.POST("/getOrderInfoWaitPay", user_api.GetOrderInfoWaitPay) //获取待支付订单
	}
	// pay
	payRouter := customerRouter.Group("/pay")
	{
		payRouter.GET("/getEnabledPayList", user_api.GetEnabledPayList) //获取已激活支付列表
	}
	// article
	articleRouter := customerRouter.Group("/article")
	{
		articleRouter.POST("/getArticleList", user_api.GetArticleList)
	}
	//ticket
	ticketRouter := customerRouter.Group("/ticket")
	{
		ticketRouter.POST("/newTicket", user_api.NewTicket)
		ticketRouter.POST("/getUserTicketList", user_api.GetUserTicketList)
		ticketRouter.POST("/updateUserTicket", user_api.UpdateUserTicket)
		ticketRouter.POST("/sendTicketMessage", user_api.SendTicketMessage)
		ticketRouter.POST("/firstTicket", user_api.FirstTicket)
	}
	//traffic
	trafficRouter := customerRouter.Group("/traffic")
	{
		trafficRouter.POST("/getSubTrafficList", user_api.GetSubTrafficList)
	}
	//finance
	financeRouter := customerRouter.Group("/finance")
	{
		financeRouter.POST("/getBalanceStatementList", user_api.GetBalanceStatementList)
		financeRouter.POST("/getCommissionStatementList", user_api.GetCommissionStatementList)
		financeRouter.POST("/getInvitationUserList", user_api.GetInvitationUserList)
		financeRouter.GET("/withdrawToBalance", user_api.WithdrawToBalance)
		financeRouter.GET("/getCommissionSummary", user_api.GetCommissionSummary)

	}
}
