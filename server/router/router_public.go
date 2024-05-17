package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api/public_api"
	middleware "github.com/ppoonk/AirGo/router/middleware"
)

func (g *GinRouter) InitPublicRouter(r *gin.RouterGroup) {
	publicRouter := r.Group("/public")
	//airgo
	airgoRouter := publicRouter.Group("/airgo")
	{
		airgoRouter.GET("/node/getNodeInfo", public_api.AGGetNodeInfo)
		airgoRouter.POST("/node/reportNodeStatus", public_api.AGReportNodeStatus)
		airgoRouter.GET("/user/getUserlist", public_api.AGGetUserlist)
		airgoRouter.POST("/user/reportUserTraffic", public_api.AGReportUserTraffic)
		airgoRouter.POST("/user/reportNodeOnlineUsers", public_api.AGReportNodeOnlineUsers)
	}
	//shop
	shopRouter := publicRouter.Group("/shop").Use(middleware.RateLimitIP())
	{
		shopRouter.GET("/epayNotify", public_api.EpayNotify)      //易支付异步回调
		shopRouter.POST("/alipayNotify", public_api.AlipayNotify) //支付宝异步验证支付结果
	}
	// 订阅
	subRouter := publicRouter.Group("/sub").Use(middleware.RateLimitIP())
	{
		subRouter.GET("/:id", public_api.GetSub)
		subRouter.GET("/:id/:name", public_api.GetSub)
	}
	// user
	userRouter := publicRouter.Group("/user").Use(middleware.RateLimitIP())
	{
		userRouter.POST("/register", public_api.Register)                   //用户注册
		userRouter.POST("/login", public_api.Login)                         //用户登录
		userRouter.POST("/resetUserPassword", public_api.ResetUserPassword) //重置密码
	}
	//server
	serverRouter := publicRouter.Group("/server").Use(middleware.RateLimitIP())
	{
		serverRouter.GET("/getThemeConfig", public_api.GetThemeConfig)     //获取主题配置
		serverRouter.GET("/getPublicSetting", public_api.GetPublicSetting) //获取公共系统设置
	}

	//code
	codeRouter := publicRouter.Group("/code").Use(middleware.RateLimitIP())
	{
		codeRouter.POST("/getEmailCode", public_api.GetEmailCode)        //获取验证码
		codeRouter.GET("/getBase64Captcha", public_api.GetBase64Captcha) //获取base64Captcha
	}
	//article
	articleRouter := publicRouter.Group("/article").Use(middleware.RateLimitIP())
	{
		articleRouter.GET("/getDefaultArticleList", public_api.GetDefaultArticleList)
	}
}
