package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api/admin_api"
	middleware "github.com/ppoonk/AirGo/router/middleware"
)

func (g *GinRouter) InitAdminRouter(r *gin.RouterGroup) {
	adminRouter := r.Group("/admin")
	adminRouter.Use(middleware.ParseJwt(), middleware.Casbin())
	// user
	userAdminRouter := adminRouter.Group("/user")
	{
		userAdminRouter.POST("/newUser", admin_api.NewUser)         //新建用户
		userAdminRouter.POST("/getUserList", admin_api.GetUserlist) //获取用户列表
		userAdminRouter.POST("/updateUser", admin_api.UpdateUser)   //修改用户
		userAdminRouter.DELETE("/deleteUser", admin_api.DeleteUser) //删除用户
		userAdminRouter.POST("/userSummary", admin_api.UserSummary) //获取订单统计
	}
	// customerService
	customerServiceAdminRouter := adminRouter.Group("/customerService")
	{
		customerServiceAdminRouter.POST("/getCustomerServiceList", admin_api.GetCustomerServiceList)
		customerServiceAdminRouter.POST("/updateCustomerService", admin_api.UpdateCustomerService)
		customerServiceAdminRouter.DELETE("/deleteCustomerService", admin_api.DeleteCustomerService)
	}
	// menu
	menuAdminRouter := adminRouter.Group("/menu")
	{
		menuAdminRouter.POST("/newMenu", admin_api.NewMenu)              //新建动态路由
		menuAdminRouter.POST("/updateMenu", admin_api.UpdateMenu)        //修改动态路由
		menuAdminRouter.DELETE("/delMenu", admin_api.DelMenu)            //删除动态路由
		menuAdminRouter.GET("/getAllMenuList", admin_api.GetAllMenuList) //获取全部动态路由
	}
	//role
	roleAdminRouter := adminRouter.Group("/role")
	{
		roleAdminRouter.POST("/newRole", admin_api.NewRole)        //添加role
		roleAdminRouter.GET("/getRoleList", admin_api.GetRoleList) //获取role list
		roleAdminRouter.POST("/updateRole", admin_api.UpdateRole)  //更新role
		roleAdminRouter.DELETE("/delRole", admin_api.DelRole)      //删除role

		roleAdminRouter.GET("/getAllPolicy", admin_api.GetAllPolicy)    //获取全部权限
		roleAdminRouter.POST("/getPolicyByID", admin_api.GetPolicyByID) //获取用户权限
	}
	// server
	serverAdminRouter := adminRouter.Group("/server")
	{
		serverAdminRouter.POST("/updateThemeConfig", admin_api.UpdateThemeConfig) //设置主题
		serverAdminRouter.GET("/getSetting", admin_api.GetSetting)                //获取系统设置
		serverAdminRouter.POST("/updateSetting", admin_api.UpdateSetting)         //修改系统设置

		serverAdminRouter.GET("/getCurrentVersion", admin_api.GetCurrentVersion)
		serverAdminRouter.GET("/getLatestVersion", admin_api.GetLatestVersion)
		serverAdminRouter.GET("/updateLatestVersion", admin_api.UpdateLatestVersion)
	}
	// node
	nodeAdminRouter := adminRouter.Group("/node")
	{
		nodeAdminRouter.POST("/newNode", admin_api.NewNode)                               //新建节点
		nodeAdminRouter.POST("/getNodeList", admin_api.GetNodeList)                       //获取节点
		nodeAdminRouter.POST("/getNodeListWithTraffic", admin_api.GetNodeListWithTraffic) //获取节点 with Traffic
		nodeAdminRouter.POST("/updateNode", admin_api.UpdateNode)                         //更新节点
		nodeAdminRouter.DELETE("/deleteNode", admin_api.DeleteNode)                       //删除节点

		nodeAdminRouter.POST("/nodeSort", admin_api.NodeSort)        //节点排序
		nodeAdminRouter.GET("/createx25519", admin_api.Createx25519) // reality x25519

		nodeAdminRouter.POST("/newNodeShared", admin_api.NewNodeShared) //新增共享节点
		nodeAdminRouter.POST("/parseUrl", admin_api.ParseUrl)           //解析url

		nodeAdminRouter.GET("/getNodeServerStatus", admin_api.GetNodeServerStatus)

	}
	// shop
	shopAdminRouter := adminRouter.Group("/shop")
	{
		shopAdminRouter.POST("/newGoods", admin_api.NewGoods)         //新建商品
		shopAdminRouter.GET("/getGoodsList", admin_api.GetGoodsList)  //查询全部商品
		shopAdminRouter.POST("/updateGoods", admin_api.UpdateGoods)   //更新商品
		shopAdminRouter.DELETE("/deleteGoods", admin_api.DeleteGoods) //删除商品
		shopAdminRouter.POST("/goodsSort", admin_api.GoodsSort)       //排序
	}
	// order
	orderAdminRouter := adminRouter.Group("/order")
	{
		orderAdminRouter.POST("/getOrderList", admin_api.GetOrderList) //获取订单列表
		orderAdminRouter.POST("/orderSummary", admin_api.OrderSummary) //获取订单统计
		orderAdminRouter.POST("/updateOrder", admin_api.UpdateOrder)   //更新用户订单
	}
	// pay
	payAdminRouter := adminRouter.Group("/pay")
	{
		payAdminRouter.POST("/newPay", admin_api.NewPay)         //新建支付
		payAdminRouter.GET("/getPayList", admin_api.GetPayList)  //获取支付列表
		payAdminRouter.POST("/updatePay", admin_api.UpdatePay)   //修改支付
		payAdminRouter.DELETE("/deletePay", admin_api.DeletePay) //删除支付
	}
	//report
	reportAdminRouter := adminRouter.Group("/report")
	{
		reportAdminRouter.POST("/getTables", admin_api.GetTables)
		reportAdminRouter.POST("/getColumn", admin_api.GetColumn)
		reportAdminRouter.POST("/reportSubmit", admin_api.ReportSubmit)

	}
	//article
	articleAdminRouter := adminRouter.Group("/article")
	{
		articleAdminRouter.POST("/newArticle", admin_api.NewArticle)
		articleAdminRouter.POST("/getArticleList", admin_api.GetArticleList)
		articleAdminRouter.POST("/updateArticle", admin_api.UpdateArticle)
		articleAdminRouter.DELETE("/deleteArticle", admin_api.DeleteArticle)
	}
	//coupon
	couponAdminRouter := adminRouter.Group("/coupon")
	{
		couponAdminRouter.POST("/newCoupon", admin_api.NewCoupon)
		couponAdminRouter.POST("/getCouponList", admin_api.GetCouponList)
		couponAdminRouter.POST("/updateCoupon", admin_api.UpdateCoupon)
		couponAdminRouter.DELETE("/deleteCoupon", admin_api.DeleteCoupon)
	}
	//access
	accessAdminRouter := adminRouter.Group("/access")
	{
		accessAdminRouter.POST("/newAccess", admin_api.NewAccessRoutes)
		accessAdminRouter.POST("/getAccessList", admin_api.GetAccessRoutesList)
		accessAdminRouter.POST("/updateAccess", admin_api.UpdateAccessRoutes)
		accessAdminRouter.DELETE("/deleteAccess", admin_api.DeleteAccessRoutes)
	}
	//migration
	migrationAdminRouter := adminRouter.Group("/migration")
	{
		migrationAdminRouter.POST("/migrationData", admin_api.Migration)
	}
	//ticket
	ticketAdminRouter := adminRouter.Group("/ticket")
	{
		ticketAdminRouter.POST("/firstTicket", admin_api.FirstTicket)
		ticketAdminRouter.POST("/getTicketList", admin_api.GetTicketList)
		ticketAdminRouter.POST("/updateTicket", admin_api.UpdateTicket)
		ticketAdminRouter.POST("/sendTicketMessage", admin_api.SendTicketMessage)
		ticketAdminRouter.DELETE("/deleteTicket", admin_api.DeleteTicket)
	}
}
