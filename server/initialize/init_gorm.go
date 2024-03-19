package initialize

import (
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/glebarez/sqlite"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/initialize_logic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// StartGorm 初始化数据库并产生数据库全局变量
func StartGorm() {
	var err error
	switch global.Config.SystemParams.DbType {
	case "mysql":
		global.DB, err = GormMysql()
	case "sqlite":
		global.DB, err = GormSqlite()
	default:
		global.DB, err = GormMysql()
	}
	if err != nil {
		panic("Database connection failed:" + err.Error())
	}
	if global.DB != nil {
		if !global.DB.Migrator().HasTable(&model.User{}) {
			global.Logrus.Info("Start creating database and initializing data...")
			RegisterTables() //创建table
			InsertInto()     //导入数据
		} else {
			RegisterTables() //AutoMigrate 自动迁移 schema
		}
	} else {
		panic("Database connection failed")
	}
}

// GormSqlite 初始化sqlite数据库
func GormSqlite() (*gorm.DB, error) {
	//db, err := sql.Open("sqlite", ":memory:")

	if db, err := gorm.Open(sqlite.Open(global.Config.Sqlite.Path), &gorm.Config{
		//SkipDefaultTransaction: true, //关闭事务
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //单数表名
		},
	}); err != nil {
		global.Logrus.Error("gorm open error:", err)
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db, nil
	}
}

// GormMysql 初始化Mysql数据库
func GormMysql() (*gorm.DB, error) {
	mysqlConfig := mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", global.Config.Mysql.Username, global.Config.Mysql.Password, global.Config.Mysql.Address, global.Config.Mysql.Port, global.Config.Mysql.Dbname, global.Config.Mysql.Config),
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		//SkipDefaultTransaction: true, //关闭事务
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}); err != nil {
		global.Logrus.Error("gorm open error:", err)
		return nil, err
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+global.Config.Mysql.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db, nil
	}
}
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

func InsertInto() {
	var err error
	defer global.Logrus.Error(err)
	var funcs = []func() error{
		initialize_logic.InsertIntoUser,
		initialize_logic.InsertIntoMenu,
		initialize_logic.InsertIntoRole,
		initialize_logic.InsertIntoUserAndRole,
		initialize_logic.InsertIntoRoleAndMenu,
		//		initialize_logic.InsertIntoGoods,
		//		initialize_logic.InsertIntoNode,
		//		initialize_logic.InsertIntoGoodsAndNodes,
		//		initialize_logic.InsertIntoCasbinRule,
		initialize_logic.InsertIntoCasbinRule,
		initialize_logic.InsertIntoTheme,
		initialize_logic.InsertIntoServer,
		initialize_logic.InsertIntoArticle,
		initialize_logic.InsertIntoAccess,
		initialize_logic.InsertIntoPay,
	}
	for _, v := range funcs {
		err = v()
		if err != nil {
			return
		}
	}
}
