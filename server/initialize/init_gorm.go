package initialize

import (
	"errors"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	utils "github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"gorm.io/driver/sqlite"
	"time"

	//"go-admin/initialize"
	//github.com/satori/go.uuid
	gormadapter "github.com/casbin/gorm-adapter/v3"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	apiPre = "/api"
	// 默认邮件验证码样式
	text1 = `<div >
  <p >欢迎使用，请及时输入验证码，区分大小写</p>
  <span style="font-size:30px">emailcode</span>
</div>`
	// 商品默认描述
	text2 = `<h3 style="color:#00BFFF">究竟什么样的终点，才配得上这一路的颠沛流离---管泽元</h3>
<h3 style="color:#DDA0DD">世界聚焦于你---管泽元</h3>`
	text3 = `## 软件下载
- [v2rayNG for Android](https://ghproxy.com/https://github.com/2dust/v2rayNG/releases/latest/download/v2rayNG_1.8.9_arm64-v8a.apk)
- [Clash Meta for Android](https://ghproxy.com/https://github.com/MetaCubeX/ClashMetaForAndroid/releases/latest/download/cmfa-2.8.9-meta-arm64-v8a-release.apk)
- [v2rayN for Windows](https://ghproxy.com/https://github.com/2dust/v2rayN/releases/latest/download/v2rayN-With-Core.zip)

## 加入我们
-[QQ](http://)
-[Telegram](http://)`
	text4 = `# 欢迎使用！最新活动如下：
1. 免费注册体验！注册送流量
2. 邀请返利

**更多套餐请查看商店**`
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() *gorm.DB {

	switch global.Config.SystemParams.DbType {
	case "mysql":
		return GormMysql()
	case "sqlite":
		return GormSqlite()
	default:
		return GormMysql()
	}
}

// 初始化sqlite数据库
func GormSqlite() *gorm.DB {

	if db, err := gorm.Open(sqlite.Open(global.Config.Sqlite.Path), &gorm.Config{
		SkipDefaultTransaction: true, //关闭事务
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //单数表名
		},
	}); err != nil {
		global.Logrus.Error("gorm open error:", err)
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db
	}
}

// 初始化Mysql数据库
func GormMysql() *gorm.DB {
	mysqlConfig := mysql.Config{
		DSN:                       global.Config.Mysql.Username + ":" + global.Config.Mysql.Password + "@tcp(" + global.Config.Mysql.Address + ":" + global.Config.Mysql.Port + ")/" + global.Config.Mysql.Dbname + "?" + global.Config.Mysql.Config,
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		SkipDefaultTransaction: true, //关闭事务
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}); err != nil {
		global.Logrus.Error("gorm open error:", err)
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+global.Config.Mysql.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db
	}
}

// RegisterTables 注册数据库表
func RegisterTables() {
	err := global.DB.AutoMigrate(
		// 用户表
		model.User{}, //动态路由表
		// model.DynamicRoute{},
		//角色表
		model.Role{},
		//casbin
		gormadapter.CasbinRule{},
		//商品
		model.Goods{},
		//订单
		model.Orders{},
		//流量统计表
		model.TrafficLog{},
		//主题
		model.Theme{},
		//系统设置参数
		model.Server{},
		//图库
		model.Gallery{},
		//文章
		model.Article{},
		//折扣
		model.Coupon{},
		//isp
		model.ISP{},
		//节点
		model.Node{},
		//共享节点
		model.NodeShared{},
		//支付
		model.Pay{},
		//访问控制
		model.Access{},
		//工单
		model.Ticket{},
		model.TicketMessage{},
	)
	if err != nil {
		global.Logrus.Error("table AutoMigrate error:", err.Error())
		panic(err)
	}
	global.Logrus.Info("table AutoMigrate success")
}
func InsertInto() {
	var err error
	defer global.Logrus.Error(err)

	var funcs = []func() error{
		InsertIntoUser,
		InsertIntoDynamicRoute,
		InsertIntoRole,
		InsertIntoUserAndRole,
		InsertIntoRoleAndMenu,
		InsertIntoGoods,
		//InsertIntoNode,
		//InsertIntoGoodsAndNodes,
		InsertIntoCasbinRule,
		InsertIntoTheme,
		InsertIntoServer,
		InsertIntoArticle,
		InsertIntoAccess,
	}
	for _, v := range funcs {
		err = v()
		if err != nil {
			return
		}
	}
}

func InsertIntoUser() error {
	uuid1 := uuid.NewV4()
	uuid2 := uuid.NewV4()
	expiedTime := time.Date(2099, 9, 9, 9, 9, 9, 0, time.Local)
	sysUserData := []model.User{
		{
			ID:             1,
			UUID:           uuid1,
			UserName:       global.Config.SystemParams.AdminEmail,
			Password:       utils.BcryptEncode(global.Config.SystemParams.AdminPassword),
			NickName:       "admin",
			InvitationCode: utils.RandomString(8),
			SubscribeInfo: model.SubscribeInfo{
				GoodsID:      1,
				SubscribeUrl: utils.RandomString(8), //随机字符串订阅url
				SubStatus:    true,
				ExpiredAt:    &expiedTime,
				T:            10737418240,
			},
		},
		{
			ID:             2,
			UUID:           uuid2,
			UserName:       "123@oicq.com",
			Password:       utils.BcryptEncode("123456"),
			NickName:       "测试123",
			InvitationCode: utils.RandomString(8),
			SubscribeInfo: model.SubscribeInfo{
				GoodsID:      1,
				SubscribeUrl: utils.RandomString(8), //随机字符串订阅url
				SubStatus:    true,
				ExpiredAt:    &expiedTime,
			},
		},
	}
	if err := global.DB.Create(&sysUserData).Error; err != nil {
		return errors.New("db.Create(&userData) Error")
	}
	return nil
}
func InsertIntoDynamicRoute() error {
	DynamicRouteData := []model.DynamicRoute{
		{ID: 1, ParentID: 0, Path: "/admin", Name: "admin", Component: "/layout/routerView/parent.vue", Meta: model.Meta{Title: "超级管理员", Icon: "iconfont icon-shouye_dongtaihui"}},
		{ID: 2, ParentID: 1, Path: "/admin/menu", Name: "adminMenu", Component: "/admin/menu/index.vue", Meta: model.Meta{Title: "菜单", Icon: "iconfont icon-caidan"}},
		{ID: 3, ParentID: 1, Path: "/admin/role", Name: "adminRole", Component: "/admin/role/index.vue", Meta: model.Meta{Title: "角色", Icon: "iconfont icon-icon-"}},
		{ID: 4, ParentID: 1, Path: "/admin/user", Name: "adminUser", Component: "/admin/user/index.vue", Meta: model.Meta{Title: "用户", Icon: "iconfont icon-gerenzhongxin"}},
		{ID: 5, ParentID: 1, Path: "/admin/order", Name: "adminOrder", Component: "/admin/order/index.vue", Meta: model.Meta{Title: "订单", Icon: "iconfont icon--chaifenhang"}},
		{ID: 6, ParentID: 1, Path: "/admin/node", Name: "adminNode", Component: "/admin/node/index.vue", Meta: model.Meta{Title: "节点", Icon: "iconfont icon-shuxingtu"}},
		{ID: 7, ParentID: 1, Path: "/admin/shop", Name: "adminShop", Component: "/admin/shop/index.vue", Meta: model.Meta{Title: "商品", Icon: "iconfont icon-zhongduancanshuchaxun"}},
		{ID: 8, ParentID: 1, Path: "/admin/system", Name: "adminSystem", Component: "/admin/system/index.vue", Meta: model.Meta{Title: "系统", Icon: "iconfont icon-xitongshezhi"}},
		{ID: 9, ParentID: 1, Path: "/admin/article", Name: "adminArticle", Component: "/admin/article/index.vue", Meta: model.Meta{Title: "文章", Icon: "iconfont icon-huanjingxingqiu"}},
		{ID: 10, ParentID: 1, Path: "/admin/coupon", Name: "adminCoupon", Component: "/admin/coupon/index.vue", Meta: model.Meta{Title: "折扣码", Icon: "ele-ShoppingBag"}},
		{ID: 11, ParentID: 1, Path: "/admin/access", Name: "adminAccess", Component: "/admin/access/index.vue", Meta: model.Meta{Title: "访问控制", Icon: "ele-ChromeFilled"}},
		{ID: 12, ParentID: 1, Path: "/admin/migration", Name: "adminMigration", Component: "/admin/migration/index.vue", Meta: model.Meta{Title: "数据迁移", Icon: "fa fa-database"}},
		{ID: 13, ParentID: 1, Path: "/admin/ticket", Name: "adminTicket", Component: "/admin/ticket/index.vue", Meta: model.Meta{Title: "工单管理", Icon: "fa fa-file-o"}},
		{ID: 14, ParentID: 1, Path: "/admin/income", Name: "adminIncome", Component: "/admin/income/index.vue", Meta: model.Meta{Title: "营收概览", Icon: "iconfont icon-xingqiu"}},

		{ID: 15, ParentID: 0, Path: "/gallery", Name: "gallery", Component: "/gallery/index.vue", Meta: model.Meta{Title: "图库", Icon: "iconfont icon-step"}},
		{ID: 16, ParentID: 0, Path: "/isp", Name: "isp", Component: "/isp/index.vue", Meta: model.Meta{Title: "套餐监控", Icon: "iconfont icon-tongzhi1"}},

		{ID: 17, ParentID: 0, Path: "/home", Name: "home", Component: "/home/index.vue", Meta: model.Meta{Title: "首页", Icon: "iconfont icon-shouye"}},
		{ID: 18, ParentID: 0, Path: "/shop", Name: "shop", Component: "/shop/index.vue", Meta: model.Meta{Title: "商店", Icon: "iconfont icon-zidingyibuju"}},
		{ID: 19, ParentID: 0, Path: "/myOrder", Name: "myOrder", Component: "/myOrder/index.vue", Meta: model.Meta{Title: "我的订单", Icon: "iconfont icon--chaifenhang"}},
		{ID: 20, ParentID: 0, Path: "/personal", Name: "personal", Component: "/personal/index.vue", Meta: model.Meta{Title: "个人信息", Icon: "iconfont icon-gerenzhongxin"}},
		{ID: 21, ParentID: 0, Path: "/serverStatus", Name: "serverStatus", Component: "/serverStatus/index.vue", Meta: model.Meta{Title: "节点状态", Icon: "iconfont icon-putong"}},
		{ID: 22, ParentID: 0, Path: "/article/notice", Name: "notice", Component: "/article/index_notice.vue", Meta: model.Meta{Title: "公告", Icon: "ele-ChatLineSquare"}},
		{ID: 23, ParentID: 0, Path: "/article/knowledge", Name: "knowledge", Component: "/article/index_knowledge.vue", Meta: model.Meta{Title: "知识库", Icon: "fa fa-book"}},
		{ID: 24, ParentID: 0, Path: "/ticket", Name: "ticket", Component: "/ticket/index.vue", Meta: model.Meta{Title: "工单", Icon: "fa fa-file-o"}},
	}
	if err := global.DB.Create(&DynamicRouteData).Error; err != nil {
		return errors.New("sys_dynamic-router_data表数据初始化失败!")
	}
	return nil
}
func InsertIntoRole() error {
	sysRoleData := []model.Role{
		{ID: 1, RoleName: "admin", Description: "超级管理员"},
		{ID: 2, RoleName: "普通用户", Description: "普通用户"},
	}
	if err := global.DB.Create(&sysRoleData).Error; err != nil {
		return errors.New("user_role表数据初始化失败!")
	}
	return nil
}
func InsertIntoUserAndRole() error {
	userAndRoleData := []model.UserAndRole{
		{UserID: 1, RoleID: 1},
		{UserID: 2, RoleID: 2},
	}
	if err := global.DB.Create(&userAndRoleData).Error; err != nil {
		return errors.New("user_and_role_data表数据初始化失败!")
	}
	return nil
}
func InsertIntoRoleAndMenu() error {
	roleAndMenuData := []model.RoleAndMenu{
		//管理员的权限
		{RoleID: 1, DynamicRouteID: 1},  //超级管理员
		{RoleID: 1, DynamicRouteID: 2},  //菜单管理
		{RoleID: 1, DynamicRouteID: 3},  //角色管理
		{RoleID: 1, DynamicRouteID: 4},  //用户管理
		{RoleID: 1, DynamicRouteID: 5},  //订单管理
		{RoleID: 1, DynamicRouteID: 6},  //节点管理
		{RoleID: 1, DynamicRouteID: 7},  //商品管理
		{RoleID: 1, DynamicRouteID: 8},  //系统设置
		{RoleID: 1, DynamicRouteID: 9},  //文章设置
		{RoleID: 1, DynamicRouteID: 10}, //折扣码管理
		{RoleID: 1, DynamicRouteID: 11}, //访问控制
		{RoleID: 1, DynamicRouteID: 12}, //数据迁移
		{RoleID: 1, DynamicRouteID: 13}, //工单管理
		{RoleID: 1, DynamicRouteID: 14},

		//{RoleID: 1, DynamicRouteID: 15},
		//{RoleID: 1, DynamicRouteID: 16},

		{RoleID: 1, DynamicRouteID: 17},
		{RoleID: 1, DynamicRouteID: 18},
		{RoleID: 1, DynamicRouteID: 19},
		{RoleID: 1, DynamicRouteID: 20},
		{RoleID: 1, DynamicRouteID: 21},
		{RoleID: 1, DynamicRouteID: 22},
		{RoleID: 1, DynamicRouteID: 23},
		{RoleID: 1, DynamicRouteID: 24},

		//普通用户的权限
		{RoleID: 2, DynamicRouteID: 17},
		{RoleID: 2, DynamicRouteID: 18},
		{RoleID: 2, DynamicRouteID: 19},
		{RoleID: 2, DynamicRouteID: 20},
		{RoleID: 2, DynamicRouteID: 21},
		{RoleID: 2, DynamicRouteID: 22},
		{RoleID: 2, DynamicRouteID: 23},
		{RoleID: 2, DynamicRouteID: 24},
	}
	if err := global.DB.Create(&roleAndMenuData).Error; err != nil {
		return errors.New("role_and_menu表数据初始化失败!")
	}

	return nil
}
func InsertIntoGoods() error {
	goodsData := []model.Goods{
		{ID: 1, Subject: "10G|30天", TotalBandwidth: 10, ExpirationDate: 30, TotalAmount: "0.00", Des: text2, TrafficResetMethod: "NotStack"},
	}
	if err := global.DB.Create(&goodsData).Error; err != nil {
		return errors.New("goods表数据初始化失败!")
	}
	return nil
}
func InsertIntoNode() error {
	key1 := utils.RandomString(32)
	key2 := utils.RandomString(32)
	nodeData := []model.Node{
		{ID: 1, Remarks: "测试节点1", Address: "www.10010.com", Path: "/path", Port: 5566, NodeType: "vless", Enabled: true, ServerKey: key1, TrafficRate: 1},
		{ID: 2, Remarks: "测试节点2", Address: "www.10086.com", Path: "/path", Port: 5566, NodeType: "vless", Enabled: true, ServerKey: key2, TrafficRate: 1},
	}
	if err := global.DB.Create(&nodeData).Error; err != nil {
		return errors.New("node表数据初始化失败!")
	}
	return nil
}
func InsertIntoGoodsAndNodes() error {
	goodsAndNodesData := []model.GoodsAndNodes{
		{GoodsID: 1, NodeID: 1},
		{GoodsID: 1, NodeID: 2},
		{GoodsID: 2, NodeID: 1},
		{GoodsID: 2, NodeID: 2},
	}
	if err := global.DB.Create(&goodsAndNodesData).Error; err != nil {
		return errors.New("goods_and_nodes表数据初始化失败!")
	}
	return nil
}
func InsertIntoCasbinRule() error {
	casbinRuleData := []gormadapter.CasbinRule{
		// user
		{ID: 1, Ptype: "p", V0: "1", V1: apiPre + "/user/changeSubHost", V2: "POST"},
		{ID: 2, Ptype: "p", V0: "1", V1: apiPre + "/user/getUserInfo", V2: "GET"},
		{ID: 3, Ptype: "p", V0: "1", V1: apiPre + "/user/changeUserPassword", V2: "POST"},
		{ID: 4, Ptype: "p", V0: "1", V1: apiPre + "/user/resetSub", V2: "GET"},
		{ID: 5, Ptype: "p", V0: "1", V1: apiPre + "/user/clockin", V2: "GET"},

		{ID: 6, Ptype: "p", V0: "1", V1: apiPre + "/user/getUserList", V2: "POST"},
		{ID: 7, Ptype: "p", V0: "1", V1: apiPre + "/user/newUser", V2: "POST"},
		{ID: 8, Ptype: "p", V0: "1", V1: apiPre + "/user/updateUser", V2: "POST"},
		{ID: 9, Ptype: "p", V0: "1", V1: apiPre + "/user/deleteUser", V2: "POST"},
		{ID: 10, Ptype: "p", V0: "1", V1: apiPre + "/user/findUser", V2: "POST"},

		// role
		{ID: 11, Ptype: "p", V0: "1", V1: apiPre + "/role/getRoleList", V2: "POST"},
		{ID: 12, Ptype: "p", V0: "1", V1: apiPre + "/role/modifyRoleInfo", V2: "POST"},
		{ID: 13, Ptype: "p", V0: "1", V1: apiPre + "/role/addRole", V2: "POST"},
		{ID: 14, Ptype: "p", V0: "1", V1: apiPre + "/role/delRole", V2: "POST"},

		// menu
		{ID: 15, Ptype: "p", V0: "1", V1: apiPre + "/menu/getAllRouteList", V2: "GET"},
		{ID: 16, Ptype: "p", V0: "1", V1: apiPre + "/menu/getAllRouteTree", V2: "GET"},
		{ID: 17, Ptype: "p", V0: "1", V1: apiPre + "/menu/newDynamicRoute", V2: "POST"},
		{ID: 18, Ptype: "p", V0: "1", V1: apiPre + "/menu/delDynamicRoute", V2: "POST"},
		{ID: 19, Ptype: "p", V0: "1", V1: apiPre + "/menu/updateDynamicRoute", V2: "POST"},
		{ID: 20, Ptype: "p", V0: "1", V1: apiPre + "/menu/findDynamicRoute", V2: "POST"},

		{ID: 21, Ptype: "p", V0: "1", V1: apiPre + "/menu/getRouteList", V2: "GET"},
		{ID: 22, Ptype: "p", V0: "1", V1: apiPre + "/menu/getRouteTree", V2: "GET"},

		//shop
		{ID: 23, Ptype: "p", V0: "1", V1: apiPre + "/shop/preCreatePay", V2: "POST"},
		{ID: 24, Ptype: "p", V0: "1", V1: apiPre + "/shop/purchase", V2: "POST"},
		{ID: 25, Ptype: "p", V0: "1", V1: apiPre + "/shop/getAllEnabledGoods", V2: "GET"},
		{ID: 26, Ptype: "p", V0: "1", V1: apiPre + "/shop/getAllGoods", V2: "GET"},
		{ID: 27, Ptype: "p", V0: "1", V1: apiPre + "/shop/findGoods", V2: "POST"},
		{ID: 28, Ptype: "p", V0: "1", V1: apiPre + "/shop/newGoods", V2: "POST"},
		{ID: 29, Ptype: "p", V0: "1", V1: apiPre + "/shop/deleteGoods", V2: "POST"},
		{ID: 30, Ptype: "p", V0: "1", V1: apiPre + "/shop/updateGoods", V2: "POST"},
		{ID: 31, Ptype: "p", V0: "1", V1: apiPre + "/shop/goodsSort", V2: "POST"},

		//node
		{ID: 32, Ptype: "p", V0: "1", V1: apiPre + "/node/getAllNode", V2: "GET"},
		{ID: 33, Ptype: "p", V0: "1", V1: apiPre + "/node/newNode", V2: "POST"},
		{ID: 34, Ptype: "p", V0: "1", V1: apiPre + "/node/deleteNode", V2: "POST"},
		{ID: 35, Ptype: "p", V0: "1", V1: apiPre + "/node/updateNode", V2: "POST"},
		{ID: 36, Ptype: "p", V0: "1", V1: apiPre + "/node/getTraffic", V2: "POST"},
		{ID: 37, Ptype: "p", V0: "1", V1: apiPre + "/node/nodeSort", V2: "POST"},
		{ID: 38, Ptype: "p", V0: "1", V1: apiPre + "/node/createx25519", V2: "GET"},

		{ID: 39, Ptype: "p", V0: "1", V1: apiPre + "/node/newNodeShared", V2: "POST"},
		{ID: 40, Ptype: "p", V0: "1", V1: apiPre + "/node/getNodeSharedList", V2: "GET"},
		{ID: 41, Ptype: "p", V0: "1", V1: apiPre + "/node/deleteNodeShared", V2: "POST"},

		//casbin
		{ID: 42, Ptype: "p", V0: "1", V1: apiPre + "/casbin/getPolicyByRoleIds", V2: "POST"},
		{ID: 43, Ptype: "p", V0: "1", V1: apiPre + "/casbin/updateCasbinPolicy", V2: "POST"},
		{ID: 44, Ptype: "p", V0: "1", V1: apiPre + "/casbin/getAllPolicy", V2: "GET"},

		//order
		{ID: 45, Ptype: "p", V0: "1", V1: apiPre + "/order/getOrderInfo", V2: "POST"},
		{ID: 46, Ptype: "p", V0: "1", V1: apiPre + "/order/getAllOrder", V2: "POST"},
		{ID: 47, Ptype: "p", V0: "1", V1: apiPre + "/order/getOrderByUserID", V2: "POST"},
		{ID: 48, Ptype: "p", V0: "1", V1: apiPre + "/order/completedOrder", V2: "POST"},
		{ID: 49, Ptype: "p", V0: "1", V1: apiPre + "/order/getMonthOrderStatistics", V2: "POST"},

		//server
		{ID: 50, Ptype: "p", V0: "1", V1: apiPre + "/server/updateThemeConfig", V2: "POST"},
		{ID: 51, Ptype: "p", V0: "1", V1: apiPre + "/server/getSetting", V2: "GET"},
		{ID: 52, Ptype: "p", V0: "1", V1: apiPre + "/server/updateSetting", V2: "POST"},

		//upload
		{ID: 53, Ptype: "p", V0: "1", V1: apiPre + "/upload/newPictureUrl", V2: "POST"},
		{ID: 54, Ptype: "p", V0: "1", V1: apiPre + "/upload/getPictureList", V2: "POST"},

		//ws
		{ID: 55, Ptype: "p", V0: "1", V1: apiPre + "/websocket/msg", V2: "GET"},

		//article
		{ID: 56, Ptype: "p", V0: "1", V1: apiPre + "/article/newArticle", V2: "POST"},
		{ID: 57, Ptype: "p", V0: "1", V1: apiPre + "/article/deleteArticle", V2: "POST"},
		{ID: 58, Ptype: "p", V0: "1", V1: apiPre + "/article/updateArticle", V2: "POST"},
		{ID: 59, Ptype: "p", V0: "1", V1: apiPre + "/article/getArticle", V2: "POST"},

		//report
		{ID: 60, Ptype: "p", V0: "1", V1: apiPre + "/report/getDB", V2: "GET"},
		{ID: 61, Ptype: "p", V0: "1", V1: apiPre + "/report/getTables", V2: "POST"},
		{ID: 62, Ptype: "p", V0: "1", V1: apiPre + "/report/getColumn", V2: "POST"},
		{ID: 63, Ptype: "p", V0: "1", V1: apiPre + "/report/reportSubmit", V2: "POST"},

		//coupon
		{ID: 64, Ptype: "p", V0: "1", V1: apiPre + "/coupon/newCoupon", V2: "POST"},
		{ID: 65, Ptype: "p", V0: "1", V1: apiPre + "/coupon/deleteCoupon", V2: "POST"},
		{ID: 66, Ptype: "p", V0: "1", V1: apiPre + "/coupon/updateCoupon", V2: "POST"},
		{ID: 67, Ptype: "p", V0: "1", V1: apiPre + "/coupon/getCoupon", V2: "POST"},

		//isp
		{ID: 68, Ptype: "p", V0: "1", V1: apiPre + "/isp/sendCode", V2: "POST"},
		{ID: 69, Ptype: "p", V0: "1", V1: apiPre + "/isp/ispLogin", V2: "POST"},
		{ID: 70, Ptype: "p", V0: "1", V1: apiPre + "/isp/getMonitorByUserID", V2: "POST"},

		//pay
		{ID: 71, Ptype: "p", V0: "1", V1: apiPre + "/pay/getEnabledPayList", V2: "GET"},
		{ID: 72, Ptype: "p", V0: "1", V1: apiPre + "/pay/getPayList", V2: "GET"},
		{ID: 73, Ptype: "p", V0: "1", V1: apiPre + "/pay/newPay", V2: "POST"},
		{ID: 74, Ptype: "p", V0: "1", V1: apiPre + "/pay/deletePay", V2: "POST"},
		{ID: 75, Ptype: "p", V0: "1", V1: apiPre + "/pay/updatePay", V2: "POST"},

		//access
		{ID: 76, Ptype: "p", V0: "1", V1: apiPre + "/access/newRoutes", V2: "POST"},
		{ID: 77, Ptype: "p", V0: "1", V1: apiPre + "/access/updateRoutes", V2: "POST"},
		{ID: 78, Ptype: "p", V0: "1", V1: apiPre + "/access/deleteRoutes", V2: "POST"},
		{ID: 79, Ptype: "p", V0: "1", V1: apiPre + "/access/getRoutesList", V2: "POST"},

		//migration
		{ID: 80, Ptype: "p", V0: "1", V1: apiPre + "/migration/fromOther", V2: "POST"},

		//ticket
		{ID: 81, Ptype: "p", V0: "1", V1: apiPre + "/ticket/newTicket", V2: "POST"},
		{ID: 82, Ptype: "p", V0: "1", V1: apiPre + "/ticket/deleteTicket", V2: "POST"},
		{ID: 83, Ptype: "p", V0: "1", V1: apiPre + "/ticket/updateTicket", V2: "POST"},
		{ID: 84, Ptype: "p", V0: "1", V1: apiPre + "/ticket/updateUserTicket", V2: "POST"},
		{ID: 85, Ptype: "p", V0: "1", V1: apiPre + "/ticket/getUserTicketList", V2: "POST"},
		{ID: 86, Ptype: "p", V0: "1", V1: apiPre + "/ticket/getTicketList", V2: "POST"},
		{ID: 87, Ptype: "p", V0: "1", V1: apiPre + "/ticket/sendTicketMessage", V2: "POST"},
		{ID: 88, Ptype: "p", V0: "1", V1: apiPre + "/ticket/getTicketMessage", V2: "POST"},

		//普通用户权限
		{ID: 89, Ptype: "p", V0: "2", V1: apiPre + "/user/changeUserPassword", V2: "POST"},
		{ID: 90, Ptype: "p", V0: "2", V1: apiPre + "/user/getUserInfo", V2: "GET"},
		{ID: 91, Ptype: "p", V0: "2", V1: apiPre + "/user/resetSub", V2: "GET"},
		{ID: 92, Ptype: "p", V0: "2", V1: apiPre + "/user/changeSubHost", V2: "POST"},
		{ID: 93, Ptype: "p", V0: "2", V1: apiPre + "/user/clockin", V2: "GET"},

		{ID: 94, Ptype: "p", V0: "2", V1: apiPre + "/menu/getRouteList", V2: "GET"},
		{ID: 95, Ptype: "p", V0: "2", V1: apiPre + "/menu/getRouteTree", V2: "GET"},

		{ID: 96, Ptype: "p", V0: "2", V1: apiPre + "/order/getOrderInfo", V2: "POST"},
		{ID: 97, Ptype: "p", V0: "2", V1: apiPre + "/order/getOrderByUserID", V2: "POST"},

		{ID: 98, Ptype: "p", V0: "2", V1: apiPre + "/shop/preCreatePay", V2: "POST"},
		{ID: 99, Ptype: "p", V0: "2", V1: apiPre + "/shop/purchase", V2: "POST"},
		{ID: 100, Ptype: "p", V0: "2", V1: apiPre + "/shop/getAllEnabledGoods", V2: "GET"},
		{ID: 101, Ptype: "p", V0: "2", V1: apiPre + "/shop/findGoods", V2: "POST"},

		{ID: 102, Ptype: "p", V0: "2", V1: apiPre + "/websocket/msg", V2: "GET"},

		{ID: 103, Ptype: "p", V0: "2", V1: apiPre + "/upload/newPictureUrl", V2: "POST"},
		{ID: 104, Ptype: "p", V0: "2", V1: apiPre + "/upload/getPictureList", V2: "POST"},

		{ID: 105, Ptype: "p", V0: "2", V1: apiPre + "/article/getArticle", V2: "POST"},

		{ID: 106, Ptype: "p", V0: "2", V1: apiPre + "/isp/sendCode", V2: "POST"},
		{ID: 107, Ptype: "p", V0: "2", V1: apiPre + "/isp/ispLogin", V2: "POST"},
		{ID: 108, Ptype: "p", V0: "2", V1: apiPre + "/isp/getMonitorByUserID", V2: "POST"},

		{ID: 109, Ptype: "p", V0: "2", V1: apiPre + "/pay/getEnabledPayList", V2: "GET"},

		{ID: 110, Ptype: "p", V0: "2", V1: apiPre + "/ticket/newTicket", V2: "POST"},
		{ID: 111, Ptype: "p", V0: "2", V1: apiPre + "/ticket/getUserTicketList", V2: "POST"},
		{ID: 112, Ptype: "p", V0: "2", V1: apiPre + "/ticket/updateUserTicket", V2: "POST"},
		{ID: 113, Ptype: "p", V0: "2", V1: apiPre + "/ticket/sendTicketMessage", V2: "POST"},
		{ID: 114, Ptype: "p", V0: "2", V1: apiPre + "/ticket/getTicketMessage", V2: "POST"},
	}
	if err := global.DB.Create(&casbinRuleData).Error; err != nil {
		return errors.New("casbin_rule表数据初始化失败!")
	}
	return nil
}
func InsertIntoTheme() error {
	themeData := model.Theme{
		ID: 1,
	}
	if err := global.DB.Create(&themeData).Error; err != nil {
		return errors.New("theme表数据初始化失败!")
	}
	return nil
}
func InsertIntoServer() error {
	settingData := model.Server{
		ID: 1,
		Email: model.Email{
			EmailContent: text1,
			EmailSubject: "hello，我的宝！",
		},
		Subscribe: model.Subscribe{
			AcceptableEmailSuffixes: "@qq.com\n@foxmail.com\n@gmail.com\n@163.com\n@126.com\n@yeah.net",
		},
	}
	if err := global.DB.Create(&settingData).Error; err != nil {
		return errors.New("server表数据初始化失败!")
	}
	return nil
}
func InsertIntoArticle() error {
	articleData := []model.Article{
		{ID: 1, Type: "home", Title: "首页自定义显示内容", Introduction: "首页自定义显示内容，可编辑，可显示与隐藏，不可删除！", Content: text3, Status: true},
		{ID: 2, Type: "home", Title: "首页弹窗内容", Introduction: "首页弹窗，可编辑，可显示与隐藏，不可删除！", Content: text4, Status: true},
	}
	if err := global.DB.Create(&articleData).Error; err != nil {
		return errors.New("article表数据初始化失败!")
	}
	return nil
}
func InsertIntoAccess() error {
	accessData := model.Access{
		ID:    1,
		Name:  "禁用流量消耗器",
		Route: "api.vv1234.cn\nshua.leyz.top\nllss.atewm.cn\nsiriling.github.io\nshidahuilang.github.io\nfu-c-k.github.io\ndb.laomoe.com\nloss.98cat.cn\nnet.ljxnet.cn",
	}
	if err := global.DB.Create(&accessData).Error; err != nil {
		return errors.New("access表数据初始化失败!")
	}
	return nil
}
