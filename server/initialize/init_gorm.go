package initialize

import (
	"AirGo/global"
	"AirGo/model"
	utils "AirGo/utils/encrypt_plugin"
	"errors"
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
// Author SliverHorn
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
		SkipDefaultTransaction: true, //关闭事务，将获得大约 30%+ 性能提升
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "gormv2_",
			SingularTable: true, //单数表名
		},
	}); err != nil {
		global.Logrus.Error("gorm.Open error:", err)
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
		DefaultStringSize:         191, // string 类型字段的默认长度
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		SkipDefaultTransaction: true, //关闭事务，将获得大约 30%+ 性能提升
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "gormv2_",
			SingularTable: true, //单数表名
		},
	}); err != nil {
		global.Logrus.Error("gorm.Open error:", err)
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
		model.User{},
		//动态路由表
		model.DynamicRoute{},
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
	)
	if err != nil {
		//os.Exit(0)
		global.Logrus.Error("table AutoMigrate error:", err.Error())
		return
	}
	global.Logrus.Info("table AutoMigrate success")
}

// 导入数据
func InsertInto(db *gorm.DB) error {
	uuid1 := uuid.NewV4()
	uuid2 := uuid.NewV4()
	expiedTime := time.Date(2099, 9, 9, 9, 9, 9, 0, time.Local)
	sysUserData := []model.User{
		{
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
			},
		},
		{
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
	if err := db.Create(&sysUserData).Error; err != nil {
		return errors.New("db.Create(&userData) Error")
	}
	//插入sys_dynamic-router_data表
	DynamicRouteData := []model.DynamicRoute{
		{ParentID: 0, Path: "/admin", Name: "admin", Component: "/layout/routerView/parent.vue", Meta: model.Meta{Title: "超级管理员", Icon: "iconfont icon-shouye_dongtaihui"}},     //id==1
		{ParentID: 1, Path: "/admin/menu", Name: "adminMenu", Component: "/admin/menu/index.vue", Meta: model.Meta{Title: "菜单管理", Icon: "iconfont icon-caidan"}},                //id==2
		{ParentID: 1, Path: "/admin/role", Name: "adminRole", Component: "/admin/role/index.vue", Meta: model.Meta{Title: "角色管理", Icon: "iconfont icon-icon-"}},                 //id==3
		{ParentID: 1, Path: "/admin/user", Name: "adminUser", Component: "/admin/user/index.vue", Meta: model.Meta{Title: "用户管理", Icon: "iconfont icon-gerenzhongxin"}},         //id==4
		{ParentID: 1, Path: "/admin/order", Name: "adminOrder", Component: "/admin/order/index.vue", Meta: model.Meta{Title: "订单管理", Icon: "iconfont icon--chaifenhang"}},       //id==5
		{ParentID: 1, Path: "/admin/node", Name: "adminNode", Component: "/admin/node/index.vue", Meta: model.Meta{Title: "节点管理", Icon: "iconfont icon-shuxingtu"}},             //id==6
		{ParentID: 1, Path: "/admin/shop", Name: "adminShop", Component: "/admin/shop/index.vue", Meta: model.Meta{Title: "商品管理", Icon: "iconfont icon-zhongduancanshuchaxun"}}, //id==7
		{ParentID: 1, Path: "/admin/system", Name: "system", Component: "/admin/system/index.vue", Meta: model.Meta{Title: "系统设置", Icon: "iconfont icon-xitongshezhi"}},         //id==8
		{ParentID: 1, Path: "/admin/article", Name: "article", Component: "/admin/article/index.vue", Meta: model.Meta{Title: "文章设置", Icon: "iconfont icon-huanjingxingqiu"}},   //id==9
		{ParentID: 1, Path: "/admin/coupon", Name: "coupon", Component: "/admin/coupon/index.vue", Meta: model.Meta{Title: "折扣码管理", Icon: "ele-ShoppingBag"}},                   //id==10

		{ParentID: 0, Path: "/home", Name: "home", Component: "/home/index.vue", Meta: model.Meta{Title: "首页", Icon: "iconfont icon-shouye"}},                           //11
		{ParentID: 0, Path: "/shop", Name: "shop", Component: "/shop/index.vue", Meta: model.Meta{Title: "商店", Icon: "iconfont icon-zidingyibuju"}},                     //12
		{ParentID: 0, Path: "/myOrder", Name: "myOrder", Component: "/myOrder/index.vue", Meta: model.Meta{Title: "我的订单", Icon: "iconfont icon--chaifenhang"}},          //13
		{ParentID: 0, Path: "/personal", Name: "personal", Component: "/personal/index.vue", Meta: model.Meta{Title: "个人信息", Icon: "iconfont icon-gerenzhongxin"}},      //14
		{ParentID: 0, Path: "/serverStatus", Name: "serverStatus", Component: "/serverStatus/index.vue", Meta: model.Meta{Title: "节点状态", Icon: "iconfont icon-putong"}}, //15
		{ParentID: 0, Path: "/gallery", Name: "gallery", Component: "/gallery/index.vue", Meta: model.Meta{Title: "无限图库", Icon: "iconfont icon-step"}},                  //16
		{ParentID: 0, Path: "/income", Name: "income", Component: "/income/index.vue", Meta: model.Meta{Title: "营收概览", Icon: "iconfont icon-xingqiu"}},                  //17
		{ParentID: 0, Path: "/isp", Name: "isp", Component: "/isp/index.vue", Meta: model.Meta{Title: "套餐监控", Icon: "iconfont icon-xingqiu"}},                           //18

		{ParentID: 0, Path: "/article/notice", Name: "notice", Component: "/article/index_notice.vue", Meta: model.Meta{Title: "公告", Icon: "ele-ChatLineSquare"}},   //19
		{ParentID: 0, Path: "/article/knowledge", Name: "knowledge", Component: "/article/index_knowledge.vue", Meta: model.Meta{Title: "知识库", Icon: "fa fa-book"}}, //20
	}
	if err := db.Create(&DynamicRouteData).Error; err != nil {
		return errors.New("sys_dynamic-router_data表数据初始化失败!")
	}
	//插入user_role表
	sysRoleData := []model.Role{
		{ID: 1, RoleName: "admin", Description: "超级管理员"},
		{ID: 2, RoleName: "普通用户", Description: "普通用户"},
	}
	if err := db.Create(&sysRoleData).Error; err != nil {
		return errors.New("user_role表数据初始化失败!")
	}
	//插入user_and_role表
	userAndRoleData := []model.UserAndRole{
		{UserID: 1, RoleID: 1},
		{UserID: 2, RoleID: 2},
	}
	if err := db.Create(&userAndRoleData).Error; err != nil {
		return errors.New("user_and_role_data表数据初始化失败!")
	}
	//插入role_and_menu
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
		{RoleID: 1, DynamicRouteID: 11}, //首页
		{RoleID: 1, DynamicRouteID: 12}, //商店
		{RoleID: 1, DynamicRouteID: 13}, //我的订单
		{RoleID: 1, DynamicRouteID: 14}, //个人信息
		{RoleID: 1, DynamicRouteID: 15}, //节点状态
		//{RoleID: 1, DynamicRouteID: 16}, //无限图库
		{RoleID: 1, DynamicRouteID: 17}, //营收概览
		{RoleID: 1, DynamicRouteID: 18}, //套餐监控
		{RoleID: 1, DynamicRouteID: 19}, //公告
		{RoleID: 1, DynamicRouteID: 20}, //知识库

		//普通用户的权限
		{RoleID: 2, DynamicRouteID: 11}, //首页
		{RoleID: 2, DynamicRouteID: 12}, //商店
		{RoleID: 2, DynamicRouteID: 13}, //我的订单
		{RoleID: 2, DynamicRouteID: 14}, //个人信息
		{RoleID: 2, DynamicRouteID: 15}, //节点状态
		//{RoleID: 2, DynamicRouteID: 16}, //无限图库
		{RoleID: 2, DynamicRouteID: 18}, //套餐监控
		{RoleID: 2, DynamicRouteID: 19}, //公告
		{RoleID: 2, DynamicRouteID: 20}, //知识库
	}
	if err := global.DB.Create(&roleAndMenuData).Error; err != nil {
		return errors.New("role_and_menu表数据初始化失败!")
	}
	//插入货物 goods
	goodsData := []model.Goods{
		{Subject: "10G|30天", TotalBandwidth: 10, ExpirationDate: 30, TotalAmount: "0.01", Des: text2},
		{Subject: "20G|180天", TotalBandwidth: 20, ExpirationDate: 180, TotalAmount: "0", Des: text2},
	}
	if err := global.DB.Create(&goodsData).Error; err != nil {
		return errors.New("goods表数据初始化失败!")
	}
	//插入node
	key1 := utils.RandomString(32)
	key2 := utils.RandomString(32)
	nodeData := []model.Node{
		{Remarks: "测试节点1", Address: "www.10010.com", Path: "/path", Port: 5566, NodeType: "vless", Enabled: true, ServerKey: key1},
		{Remarks: "测试节点2", Address: "www.10086.com", Path: "/path", Port: 5566, NodeType: "vless", Enabled: true, ServerKey: key2},
	}
	if err := global.DB.Create(&nodeData).Error; err != nil {
		return errors.New("node表数据初始化失败!")
	}
	//插入goods_and_nodes
	goodsAndNodesData := []model.GoodsAndNodes{
		{GoodsID: 1, NodeID: 1},
		{GoodsID: 1, NodeID: 2},
		{GoodsID: 2, NodeID: 1},
		{GoodsID: 2, NodeID: 2},
	}
	if err := global.DB.Create(&goodsAndNodesData).Error; err != nil {
		return errors.New("goods_and_nodes表数据初始化失败!")
	}
	// 插入casbin_rule
	casbinRuleData := []gormadapter.CasbinRule{
		// user
		{Ptype: "p", V0: "1", V1: apiPre + "/user/changeSubHost", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/user/changeUserPassword", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/user/resetSub", V2: "GET"},

		{Ptype: "p", V0: "1", V1: apiPre + "/user/getUserList", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/user/newUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/user/updateUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/user/deleteUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/user/findUser", V2: "POST"},

		// role
		{Ptype: "p", V0: "1", V1: apiPre + "/role/getRoleList", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/role/modifyRoleInfo", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/role/addRole", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/role/delRole", V2: "POST"},

		// menu
		{Ptype: "p", V0: "1", V1: apiPre + "/menu/getAllRouteList", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/menu/getAllRouteTree", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/menu/newDynamicRoute", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/menu/delDynamicRoute", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/menu/updateDynamicRoute", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/menu/findDynamicRoute", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPre + "/menu/getRouteList", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/menu/getRouteTree", V2: "GET"},

		//shop
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/preCreatePay", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/purchase", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/getAllEnabledGoods", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/getAllGoods", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/findGoods", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/newGoods", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/deleteGoods", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/updateGoods", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/shop/goodsSort", V2: "POST"},

		//node
		{Ptype: "p", V0: "1", V1: apiPre + "/node/getAllNode", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/node/newNode", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/node/deleteNode", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/node/updateNode", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/node/getTraffic", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/node/nodeSort", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPre + "/node/newNodeShared", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/node/getNodeSharedList", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/node/deleteNodeShared", V2: "POST"},

		//casbin
		{Ptype: "p", V0: "1", V1: apiPre + "/casbin/getPolicyByRoleIds", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/casbin/updateCasbinPolicy", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/casbin/getAllPolicy", V2: "GET"},

		//order
		{Ptype: "p", V0: "1", V1: apiPre + "/order/getOrderInfo", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/order/getAllOrder", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/order/getOrderByUserID", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/order/completedOrder", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/order/getMonthOrderStatistics", V2: "POST"},

		//system
		{Ptype: "p", V0: "1", V1: apiPre + "/system/updateThemeConfig", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/system/getSetting", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/system/updateSetting", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/system/createx25519", V2: "GET"},

		//upload
		{Ptype: "p", V0: "1", V1: apiPre + "/upload/newPictureUrl", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/upload/getPictureList", V2: "POST"},

		//ws
		{Ptype: "p", V0: "1", V1: apiPre + "/websocket/msg", V2: "GET"},

		//article
		{Ptype: "p", V0: "1", V1: apiPre + "/article/newArticle", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/article/deleteArticle", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/article/updateArticle", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/article/getArticle", V2: "POST"},

		//report
		{Ptype: "p", V0: "1", V1: apiPre + "/report/getDB", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/report/getTables", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/report/getColumn", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/report/reportSubmit", V2: "POST"},

		//coupon
		{Ptype: "p", V0: "1", V1: apiPre + "/coupon/newCoupon", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/coupon/deleteCoupon", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/coupon/updateCoupon", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/coupon/getCoupon", V2: "POST"},

		//isp
		{Ptype: "p", V0: "1", V1: apiPre + "/isp/sendCode", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/isp/ispLogin", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/isp/getMonitorByUserID", V2: "POST"},

		//pay
		{Ptype: "p", V0: "1", V1: apiPre + "/pay/getEnabledPayList", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/pay/getPayList", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPre + "/pay/newPay", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/pay/deletePay", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPre + "/pay/updatePay", V2: "POST"},

		//普通用户权限
		{Ptype: "p", V0: "2", V1: apiPre + "/user/changeUserPassword", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPre + "/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPre + "/user/resetSub", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPre + "/user/changeSubHost", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPre + "/menu/getRouteList", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPre + "/menu/getRouteTree", V2: "GET"},

		{Ptype: "p", V0: "2", V1: apiPre + "/order/getOrderInfo", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPre + "/order/getOrderByUserID", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPre + "/shop/preCreatePay", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPre + "/shop/purchase", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPre + "/shop/getAllEnabledGoods", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPre + "/shop/findGoods", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPre + "/websocket/msg", V2: "GET"},

		{Ptype: "p", V0: "2", V1: apiPre + "/upload/newPictureUrl", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPre + "/upload/getPictureList", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPre + "/article/getArticle", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPre + "/isp/sendCode", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPre + "/isp/ispLogin", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPre + "/isp/getMonitorByUserID", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPre + "/pay/getEnabledPayList", V2: "GET"},
	}
	if err := global.DB.Create(&casbinRuleData).Error; err != nil {
		return errors.New("casbin_rule表数据初始化失败!")
	}
	//主题配置
	themeData := model.Theme{
		ID: 1,
	}
	if err := global.DB.Create(&themeData).Error; err != nil {
		return errors.New("theme表数据初始化失败!")
	}

	//系统设置
	settingData := model.Server{
		ID: 1,
		Email: model.Email{
			EmailContent: text1,
			EmailSubject: "hello，我的宝！",
		},
	}
	if err := global.DB.Create(&settingData).Error; err != nil {
		return errors.New("server表数据初始化失败!")
	}
	//文章
	articleData := []model.Article{
		{Type: "home", Title: "首页自定义显示内容", Introduction: "首页自定义显示内容，可编辑，可显示与隐藏，不可删除！", Content: text3},
		{Type: "home", Title: "首页弹窗内容", Introduction: "首页弹窗，可编辑，可显示与隐藏，不可删除！", Content: text4},
	}
	if err := global.DB.Create(&articleData).Error; err != nil {
		return errors.New("article表数据初始化失败!")
	}
	return nil
}
