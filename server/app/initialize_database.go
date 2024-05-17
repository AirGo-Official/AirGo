package app

import (
	"errors"
	"fmt"
	"strings"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	utils "github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DataBase struct {
	routes gin.RoutesInfo
}

func NewDataBase() *DataBase {
	return &DataBase{
		routes: nil,
	}
}
func (d *DataBase) ConnectDatabase() {
	var err error
	var dialector gorm.Dialector

	switch global.Config.SystemParams.DbType {
	case "mysql", "mariadb":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
			global.Config.Mysql.Username,
			global.Config.Mysql.Password,
			global.Config.Mysql.Address,
			global.Config.Mysql.Port,
			global.Config.Mysql.Dbname,
			global.Config.Mysql.Config)
		mysqlConfig := mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         191,
			SkipInitializeWithVersion: false,
		}
		dialector = mysql.New(mysqlConfig)
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(global.Config.Sqlite.Path)
	default:
		panic("Illegal database type")
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		//SkipDefaultTransaction: true, //关闭事务
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("Database connection failed:" + err.Error())
	}
	switch global.Config.SystemParams.DbType {
	case "mysql", "mariadb":
		db.InstanceSet("gorm:table_options", "ENGINE="+global.Config.Mysql.Engine)
	default:
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
	global.DB = db

	if !global.DB.Migrator().HasTable(&model.User{}) {
		global.Logrus.Info("Start creating database and initializing data...")
		d.RegisterTables() //创建table
		d.InsertInto()     //导入数据
	} else {
		d.RegisterTables() //AutoMigrate 自动迁移 schema
	}
}

func (d *DataBase) RegisterTables() {
	err := global.DB.AutoMigrate(
		model.User{},                // 用户表
		model.CustomerService{},     //用户服务
		model.Menu{},                //动态路由表
		model.Role{},                //角色表
		gormadapter.CasbinRule{},    //casbin
		model.Goods{},               //商品
		model.Order{},               //订单
		model.NodeTrafficLog{},      //流量统计表
		model.Theme{},               //主题
		model.Server{},              //系统设置参数
		model.Article{},             //文章
		model.Coupon{},              //折扣
		model.Node{},                //节点
		model.Pay{},                 //支付
		model.Access{},              //访问控制
		model.Ticket{},              //工单
		model.TicketMessage{},       //工单消息
		model.UserTrafficLog{},      //用户流量记录
		model.BalanceStatement{},    //用户余额明细
		model.CommissionStatement{}, //邀请佣金明细
	)
	if err != nil {
		global.Logrus.Error("table AutoMigrate error:", err.Error())
		panic(err)
	}
	global.Logrus.Info("table AutoMigrate success")
}

func (d *DataBase) InsertInto() {
	var err error
	defer global.Logrus.Error(err)
	var funcs = []func() error{
		d.InsertIntoUser,
		d.InsertIntoMenu,
		d.InsertIntoRole,
		d.InsertIntoUserAndRole,
		d.InsertIntoRoleAndMenu,
		//		d.InsertIntoGoods,
		//		d.InsertIntoNode,
		//		d.InsertIntoGoodsAndNodes,
		d.InsertIntoCasbinRule,
		d.InsertIntoTheme,
		d.InsertIntoServer,
		d.InsertIntoArticle,
		d.InsertIntoAccess,
		d.InsertIntoPay,
	}
	for _, v := range funcs {
		err = v()
		if err != nil {
			return
		}
	}
}

func (d *DataBase) InsertIntoUser() error {
	sysUserData := []model.User{
		{
			//CreatedAt: time.Time{},
			//UpdatedAt: time.Time{},
			//DeletedAt: nil,

			ID:             1,
			UserName:       global.Config.SystemParams.AdminEmail,
			Password:       utils.BcryptEncode(global.Config.SystemParams.AdminPassword),
			NickName:       "admin",
			Avatar:         "https://api.multiavatar.com/admin.svg",
			Enable:         true,
			InvitationCode: utils.RandomString(8),
			Balance:        0,
			//TgID:           0,
			//RoleGroup:      nil,
			//Orders:         nil,
		},
	}
	if err := global.DB.Create(&sysUserData).Error; err != nil {
		return errors.New("db.Create(&userData) Error")
	}
	return nil
}
func (d *DataBase) InsertIntoMenu() error {
	DynamicRouteData := []model.Menu{
		{ID: 1, ParentID: 0, Remarks: "管理员", Path: "/admin", Name: "admin", Component: "/layout/routerView/parent.vue", Meta: model.Meta{Title: "message.router.admin", Icon: "ri-admin-line"}},
		{ID: 2, ParentID: 1, Remarks: "菜单", Path: "/admin/menu", Name: "adminMenu", Component: "/admin/menu/index.vue", Meta: model.Meta{Title: "message.router.adminMenu", Icon: "ri-menu-add-line"}},
		{ID: 3, ParentID: 1, Remarks: "角色", Path: "/admin/role", Name: "adminRole", Component: "/admin/role/index.vue", Meta: model.Meta{Title: "message.router.adminRole", Icon: "ri-file-shield-2-line"}},
		{ID: 4, ParentID: 1, Remarks: "用户", Path: "/admin/user", Name: "adminUser", Component: "/admin/user/index.vue", Meta: model.Meta{Title: "message.router.adminUser", Icon: "ri-user-search-line"}},
		{ID: 5, ParentID: 1, Remarks: "订单", Path: "/admin/order", Name: "adminOrder", Component: "/admin/order/index.vue", Meta: model.Meta{Title: "message.router.adminOrder", Icon: "ri-receipt-line"}},
		{ID: 6, ParentID: 1, Remarks: "节点", Path: "/admin/node", Name: "adminNode", Component: "/admin/node/index.vue", Meta: model.Meta{Title: "message.router.adminNode", Icon: "ri-cloud-line"}},
		{ID: 7, ParentID: 1, Remarks: "商店", Path: "/admin/shop", Name: "adminShop", Component: "/admin/shop/index.vue", Meta: model.Meta{Title: "message.router.adminShop", Icon: "ri-store-3-line"}},
		{ID: 8, ParentID: 1, Remarks: "系统", Path: "/admin/system", Name: "adminSystem", Component: "/admin/system/index.vue", Meta: model.Meta{Title: "message.router.adminSystem", Icon: "ri-settings-4-line"}},
		{ID: 9, ParentID: 1, Remarks: "文章", Path: "/admin/article", Name: "adminArticle", Component: "/admin/article/index.vue", Meta: model.Meta{Title: "message.router.adminArticle", Icon: "ri-file-settings-line"}},
		{ID: 10, ParentID: 1, Remarks: "工单管理", Path: "/admin/ticket", Name: "adminTicket", Component: "/admin/ticket/index.vue", Meta: model.Meta{Title: "message.router.adminTicket", Icon: "ri-customer-service-2-line"}},
		{ID: 11, ParentID: 1, Remarks: "营收概览", Path: "/admin/income", Name: "adminIncome", Component: "/admin/income/index.vue", Meta: model.Meta{Title: "message.router.adminIncome", Icon: "ri-bar-chart-box-line"}},

		{ID: 12, ParentID: 0, Remarks: "首页", Path: "/home", Name: "home", Component: "/home/index.vue", Meta: model.Meta{Title: "message.router.home", Icon: "ri-home-3-line"}},
		{ID: 13, ParentID: 0, Remarks: "商店", Path: "/shop", Name: "shop", Component: "/shop/index.vue", Meta: model.Meta{Title: "message.router.shop", Icon: "ri-store-3-line"}},
		{ID: 14, ParentID: 0, Remarks: "我的订单", Path: "/myOrder", Name: "myOrder", Component: "/myOrder/index.vue", Meta: model.Meta{Title: "message.router.myOrder", Icon: "ri-bill-line"}},
		{ID: 15, ParentID: 0, Remarks: "个人信息", Path: "/personal", Name: "personal", Component: "/personal/index.vue", Meta: model.Meta{Title: "message.router.personal", Icon: "ri-user-line"}},
		{ID: 16, ParentID: 0, Remarks: "文档", Path: "/documents", Name: "documents", Component: "/documents/index.vue", Meta: model.Meta{Title: "message.router.documents", Icon: "ri-article-line"}},
		{ID: 17, ParentID: 0, Remarks: "工单", Path: "/ticket", Name: "ticket", Component: "/ticket/index.vue", Meta: model.Meta{Title: "message.router.ticket", Icon: "ri-customer-service-2-line"}},
		{ID: 18, ParentID: 0, Remarks: "财务中心", Path: "/finance", Name: "finance", Component: "/finance/index.vue", Meta: model.Meta{Title: "message.router.finance", Icon: "ri-wallet-line"}},
	}
	if err := global.DB.Create(&DynamicRouteData).Error; err != nil {
		return errors.New("sys_dynamic-router_data表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoRole() error {
	sysRoleData := []model.Role{
		{ID: 1, RoleName: "admin", Description: "超级管理员"},
		{ID: 2, RoleName: "普通用户", Description: "普通用户"},
	}
	if err := global.DB.Create(&sysRoleData).Error; err != nil {
		return errors.New("user_role表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoUserAndRole() error {
	userAndRoleData := []model.UserAndRole{
		{UserID: 1, RoleID: 1},
	}
	if err := global.DB.Create(&userAndRoleData).Error; err != nil {
		return errors.New("user_and_role_data表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoRoleAndMenu() error {
	roleAndMenuData := []model.RoleAndMenu{
		//管理员的权限
		{RoleID: 1, MenuID: 1},  //超级管理员
		{RoleID: 1, MenuID: 2},  //菜单管理
		{RoleID: 1, MenuID: 3},  //角色管理
		{RoleID: 1, MenuID: 4},  //用户管理
		{RoleID: 1, MenuID: 5},  //订单管理
		{RoleID: 1, MenuID: 6},  //节点管理
		{RoleID: 1, MenuID: 7},  //商品管理
		{RoleID: 1, MenuID: 8},  //系统设置
		{RoleID: 1, MenuID: 9},  //文章设置
		{RoleID: 1, MenuID: 10}, //折扣码管理
		{RoleID: 1, MenuID: 11}, //访问控制
		{RoleID: 1, MenuID: 12}, //数据迁移
		{RoleID: 1, MenuID: 13}, //工单管理
		{RoleID: 1, MenuID: 14},
		{RoleID: 1, MenuID: 15},
		{RoleID: 1, MenuID: 16},
		{RoleID: 1, MenuID: 17},
		{RoleID: 1, MenuID: 18},
		//普通用户的权限
		{RoleID: 2, MenuID: 12},
		{RoleID: 2, MenuID: 13},
		{RoleID: 2, MenuID: 14},
		{RoleID: 2, MenuID: 15},
		{RoleID: 2, MenuID: 16},
		{RoleID: 2, MenuID: 17},
		{RoleID: 2, MenuID: 18},
	}
	if err := global.DB.Create(&roleAndMenuData).Error; err != nil {
		return errors.New("role_and_menu表数据初始化失败!")
	}

	return nil
}
func (d *DataBase) InsertIntoGoods() error {
	goodsData := []model.Goods{
		{ID: 1, Subject: "10G|30天", TotalBandwidth: 10, Price: "0.00", Des: constant.Text2},
	}
	if err := global.DB.Create(&goodsData).Error; err != nil {
		return errors.New("goods表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoNode() error {
	nodeData := []model.Node{
		{ID: 1, Remarks: "测试节点1", Address: "www.10010.com", Path: "/path", Port: 5566, NodeType: constant.NODE_TYPE_NORMAL, Protocol: constant.NODE_PROTOCOL_VLESS, Enabled: true, TrafficRate: 1},
		{ID: 2, Remarks: "测试节点2", Address: "www.10086.com", Path: "/path", Port: 5566, NodeType: constant.NODE_TYPE_NORMAL, Protocol: constant.NODE_PROTOCOL_VLESS, Enabled: true, TrafficRate: 1},
	}
	if err := global.DB.Create(&nodeData).Error; err != nil {
		return errors.New("node表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoGoodsAndNodes() error {
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
func (d *DataBase) InsertIntoCasbinRule() error {
	casbinRuleData := []gormadapter.CasbinRule{}
	adminCasbinRuleData := []gormadapter.CasbinRule{}
	userCasbinRuleData := []gormadapter.CasbinRule{}

	for _, v := range d.routes {
		if strings.Index(v.Path, "public") != -1 || strings.Index(v.Path, "airgo") != -1 {
			continue
		} else if strings.Index(v.Path, "admin") != -1 || strings.Index(v.Path, "swagger") != -1 {
			adminCasbinRuleData = append(adminCasbinRuleData, gormadapter.CasbinRule{Ptype: "p", V0: "1", V1: v.Path, V2: v.Method})
		} else {
			userCasbinRuleData = append(userCasbinRuleData, gormadapter.CasbinRule{Ptype: "p", V0: "1", V1: v.Path, V2: v.Method})
			userCasbinRuleData = append(userCasbinRuleData, gormadapter.CasbinRule{Ptype: "p", V0: "2", V1: v.Path, V2: v.Method})
		}

	}
	casbinRuleData = append(casbinRuleData, adminCasbinRuleData...)
	casbinRuleData = append(casbinRuleData, userCasbinRuleData...)
	//编号
	var i uint = 1
	for k, _ := range casbinRuleData {
		casbinRuleData[k].ID = i
		i++
	}
	if err := global.DB.Create(&casbinRuleData).Error; err != nil {
		return errors.New("casbin_rule表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoTheme() error {
	themeData := model.Theme{
		ID: 1,
	}
	if err := global.DB.Create(&themeData).Error; err != nil {
		return errors.New("theme表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoServer() error {
	settingData := model.Server{
		ID: 1,
		Email: model.Email{
			EmailContent: constant.Text1,
			EmailSubject: "Hello AirGo!",
		},
		Website: model.Website{
			EnableRegister:          true,
			AcceptableEmailSuffixes: "@qq.com\n@foxmail.com\n@gmail.com\n@163.com\n@126.com\n@yeah.net",
			EnableBase64Captcha:     true,
			EnableAssetsApi:         true,
		},
		Finance: model.Finance{
			EnableInvitationCommission: false,
			CommissionRate:             0.1,
			WithdrawThreshold:          50,
			EnableLottery:              false,
			Jackpot: model.Jackpot{
				{0.01, 6},
				{0.02, 5},
				{0.03, 4},
				{0.04, 3},
				{0.05, 2},
				{0.06, 1},
			},
		},
		Security: model.Security{
			Captcha: model.Captcha{
				KeyLong:            6,
				ImgWidth:           240,
				ImgHeight:          80,
				OpenCaptcha:        2,
				OpenCaptchaTimeOut: 300,
			},
			JWT: model.JWT{
				SigningKey:  "AirGo",
				ExpiresTime: "30d",
				BufferTime:  "1d",
				Issuer:      "AirGo",
			},
			RateLimitParams: model.RateLimitParams{
				IPRoleParam: 600,
				VisitParam:  600,
			},
		},
		Subscribe: model.Subscribe{
			SubName:   "AirGo",
			TEK:       "airgo",
			SurgeRule: constant.DEFAULT_SURGE_RULE,
			ClashRule: constant.DEFAULT_CLASH_RULE,
		},
	}
	if err := global.DB.Create(&settingData).Error; err != nil {
		return errors.New("server表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoArticle() error {
	articleData := []model.Article{
		{ID: 1, Type: "home", Title: "首页自定义显示内容", Introduction: "首页自定义显示内容，可编辑，不可删除！", Content: constant.DefaultHtml, Status: true},
		{ID: 2, Type: "dialog", Title: "首页弹窗内容", Introduction: "首页弹窗，可编辑，可显示与隐藏，不可删除！", Content: constant.DefaultDialog, Status: true},
	}
	if err := global.DB.Create(&articleData).Error; err != nil {
		return errors.New("article表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoAccess() error {
	accessData := []model.Access{
		{
			Name:  "禁用流量消耗器",
			Route: constant.Rule1,
		},
		{
			Name:  "禁用一些敏感网站和测速网站",
			Route: constant.Rule2,
		}}
	if err := global.DB.Create(&accessData).Error; err != nil {
		return errors.New("access表数据初始化失败!")
	}
	return nil
}
func (d *DataBase) InsertIntoPay() error {
	payData := model.Pay{
		Name:       "Balance payment",
		PayType:    constant.PAY_TYPE_BALANCE,
		PayLogoUrl: "https://telegraph-image.pages.dev/file/a57a72c5572277ff6b48f.jpg",
		Status:     true,
		AliPay:     model.AliPay{},
		Epay:       model.Epay{},
	}
	if err := global.DB.Create(&payData).Error; err != nil {
		return errors.New("pay表数据初始化失败!")
	}
	return nil
}

func (d *DataBase) DefaultForRoleMenuCasbin() error {
	fmt.Println("初始化数据库 casbin")
	err := global.DB.Where("id > 0").Delete(&gormadapter.CasbinRule{}).Error
	if err != nil {
		return err
	}
	err = d.InsertIntoCasbinRule()
	if err != nil {
		return err
	}
	fmt.Println("初始化数据库 角色和菜单")
	err = global.DB.Where("role_id > 0").Delete(&model.RoleAndMenu{}).Error //先删除role_and_menu
	if err != nil {
		return err
	}
	err = global.DB.Where("id > 0").Delete(&model.Menu{}).Error //再删除菜单
	if err != nil {
		return err
	}
	err = d.InsertIntoMenu() //插入新的菜单
	if err != nil {
		return err
	}
	return d.InsertIntoRoleAndMenu() //插入新的role_and_menu

}
