package database

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
)

// RegisterTables 注册数据库表
func RegisterTables() {
	err := global.DB.AutoMigrate(
		model.User{},             // 用户表
		model.CustomerService{},  //用户服务
		model.Menu{},             //动态路由表
		model.Role{},             //角色表
		gormadapter.CasbinRule{}, //casbin
		model.Goods{},            //商品
		model.Order{},            //订单
		model.NodeTrafficLog{},   //流量统计表
		model.Theme{},            //主题
		model.Server{},           //系统设置参数
		model.Article{},          //文章
		model.Coupon{},           //折扣
		model.Node{},             //节点
		model.Pay{},              //支付
		model.Access{},           //访问控制
		model.Ticket{},           //工单
		model.TicketMessage{},    //工单消息
		model.UserTrafficLog{},   //用户流量记录
	)
	if err != nil {
		global.Logrus.Error("table AutoMigrate error:", err.Error())
		panic(err)
	}
	global.Logrus.Info("table AutoMigrate success")
}
