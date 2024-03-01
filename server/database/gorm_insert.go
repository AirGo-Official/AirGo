package database

import (
	"errors"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/router"
	utils "github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"strings"
)

func InsertInto() {
	var err error
	defer global.Logrus.Error(err)

	var funcs = []func() error{
		InsertIntoUser,
		InsertIntoMenu,
		InsertIntoRole,
		InsertIntoUserAndRole,
		InsertIntoRoleAndMenu,
		//InsertIntoGoods,
		//InsertIntoNode,
		//InsertIntoGoodsAndNodes,
		//InsertIntoCasbinRule,
		InsertIntoCasbinRule,
		InsertIntoTheme,
		InsertIntoServer,
		InsertIntoArticle,
		InsertIntoAccess,
		InsertIntoPay,
	}
	for _, v := range funcs {
		err = v()
		if err != nil {
			return
		}
	}
}

func InsertIntoUser() error {
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
			//ReferrerCode:   "",
			Balance: 0,
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
func InsertIntoMenu() error {
	DynamicRouteData := []model.Menu{
		{ID: 1, ParentID: 0, Remarks: "管理员", Path: "/admin", Name: "admin", Component: "/layout/routerView/parent.vue", Meta: model.Meta{Title: "message.router.admin", Icon: "iconfont icon-shouye_dongtaihui"}},
		{ID: 2, ParentID: 1, Remarks: "菜单", Path: "/admin/menu", Name: "adminMenu", Component: "/admin/menu/index.vue", Meta: model.Meta{Title: "message.router.adminMenu", Icon: "iconfont icon-caidan"}},
		{ID: 3, ParentID: 1, Remarks: "角色", Path: "/admin/role", Name: "adminRole", Component: "/admin/role/index.vue", Meta: model.Meta{Title: "message.router.adminRole", Icon: "iconfont icon-icon-"}},
		{ID: 4, ParentID: 1, Remarks: "用户", Path: "/admin/user", Name: "adminUser", Component: "/admin/user/index.vue", Meta: model.Meta{Title: "message.router.adminUser", Icon: "iconfont icon-gerenzhongxin"}},
		{ID: 5, ParentID: 1, Remarks: "订单", Path: "/admin/order", Name: "adminOrder", Component: "/admin/order/index.vue", Meta: model.Meta{Title: "message.router.adminOrder", Icon: "iconfont icon--chaifenhang"}},
		{ID: 6, ParentID: 1, Remarks: "节点", Path: "/admin/node", Name: "adminNode", Component: "/admin/node/index.vue", Meta: model.Meta{Title: "message.router.adminNode", Icon: "iconfont icon-shuxingtu"}},
		{ID: 7, ParentID: 1, Remarks: "商店", Path: "/admin/shop", Name: "adminShop", Component: "/admin/shop/index.vue", Meta: model.Meta{Title: "message.router.adminShop", Icon: "iconfont icon-zhongduancanshuchaxun"}},
		{ID: 8, ParentID: 1, Remarks: "系统", Path: "/admin/system", Name: "adminSystem", Component: "/admin/system/index.vue", Meta: model.Meta{Title: "message.router.adminSystem", Icon: "iconfont icon-xitongshezhi"}},
		{ID: 9, ParentID: 1, Remarks: "文章", Path: "/admin/article", Name: "adminArticle", Component: "/admin/article/index.vue", Meta: model.Meta{Title: "message.router.adminArticle", Icon: "iconfont icon-huanjingxingqiu"}},
		{ID: 10, ParentID: 1, Remarks: "工单管理", Path: "/admin/ticket", Name: "adminTicket", Component: "/admin/ticket/index.vue", Meta: model.Meta{Title: "message.router.adminTicket", Icon: "ele-DocumentRemove"}},
		{ID: 11, ParentID: 1, Remarks: "营收概览", Path: "/admin/income", Name: "adminIncome", Component: "/admin/income/index.vue", Meta: model.Meta{Title: "message.router.adminIncome", Icon: "iconfont icon-xingqiu"}},

		{ID: 12, ParentID: 0, Remarks: "首页", Path: "/home", Name: "home", Component: "/home/index.vue", Meta: model.Meta{Title: "message.router.home", Icon: "iconfont icon-shouye"}},
		{ID: 13, ParentID: 0, Remarks: "商店", Path: "/shop", Name: "shop", Component: "/shop/index.vue", Meta: model.Meta{Title: "message.router.shop", Icon: "iconfont icon-zidingyibuju"}},
		{ID: 14, ParentID: 0, Remarks: "我的订单", Path: "/myOrder", Name: "myOrder", Component: "/myOrder/index.vue", Meta: model.Meta{Title: "message.router.myOrder", Icon: "iconfont icon--chaifenhang"}},
		{ID: 15, ParentID: 0, Remarks: "个人信息", Path: "/personal", Name: "personal", Component: "/personal/index.vue", Meta: model.Meta{Title: "message.router.personal", Icon: "iconfont icon-gerenzhongxin"}},
		{ID: 16, ParentID: 0, Remarks: "文档", Path: "/documents", Name: "documents", Component: "/documents/index.vue", Meta: model.Meta{Title: "message.router.documents", Icon: "ele-ChatLineSquare"}},
		{ID: 17, ParentID: 0, Remarks: "工单", Path: "/ticket", Name: "ticket", Component: "/ticket/index.vue", Meta: model.Meta{Title: "message.router.ticket", Icon: "ele-DocumentRemove"}},
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
	}
	if err := global.DB.Create(&roleAndMenuData).Error; err != nil {
		return errors.New("role_and_menu表数据初始化失败!")
	}

	return nil
}
func InsertIntoGoods() error {
	goodsData := []model.Goods{
		{ID: 1, Subject: "10G|30天", TotalBandwidth: 10, Price: "0.00", Des: text2},
	}
	if err := global.DB.Create(&goodsData).Error; err != nil {
		return errors.New("goods表数据初始化失败!")
	}
	return nil
}
func InsertIntoNode() error {
	nodeData := []model.Node{
		{ID: 1, Remarks: "测试节点1", Address: "www.10010.com", Path: "/path", Port: 5566, NodeType: constant.NODE_TYPE_NORMAL, Protocol: constant.NODE_PROTOCOL_VLESS, Enabled: true, TrafficRate: 1},
		{ID: 2, Remarks: "测试节点2", Address: "www.10086.com", Path: "/path", Port: 5566, NodeType: constant.NODE_TYPE_NORMAL, Protocol: constant.NODE_PROTOCOL_VLESS, Enabled: true, TrafficRate: 1},
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
	routes := router.Router.Routes()
	//fmt.Println("routes:", routes)

	casbinRuleData := []gormadapter.CasbinRule{}
	adminCasbinRuleData := []gormadapter.CasbinRule{}
	userCasbinRuleData := []gormadapter.CasbinRule{}

	for _, v := range routes {
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
		Website: model.Website{
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
		{ID: 1, Type: "home", Title: "首页自定义显示内容", Introduction: "首页自定义显示内容，可编辑，不可删除！", Content: defaultHtml, Status: true},
		{ID: 2, Type: "dialog", Title: "首页弹窗内容", Introduction: "首页弹窗，可编辑，可显示与隐藏，不可删除！", Content: defaultDialog, Status: true},
	}
	if err := global.DB.Create(&articleData).Error; err != nil {
		return errors.New("article表数据初始化失败!")
	}
	return nil
}
func InsertIntoAccess() error {
	accessData := []model.Access{
		{
			Name:  "禁用流量消耗器",
			Route: rule1,
		},
		{
			Name:  "禁用一些敏感网站和测速网站",
			Route: rule2,
		}}
	if err := global.DB.Create(&accessData).Error; err != nil {
		return errors.New("access表数据初始化失败!")
	}
	return nil
}
func InsertIntoPay() error {
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
