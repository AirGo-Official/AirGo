package initialize

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func InitializeAll() {
	InitLogrus()            //logrus
	global.VP = InitViper() //初始化Viper
	global.DB = Gorm()      //gorm连接数据库
	if global.DB != nil {
		if !global.DB.Migrator().HasTable(&model.User{}) {
			global.Logrus.Info("未找到sys_user库表,开始建表并初始化数据...")
			RegisterTables() //创建table
			InsertInto()     //导入数据
		} else {
			RegisterTables() //AutoMigrate 自动迁移 schema
		}
	} else {
		panic("数据库连接失败")
	}
	InitServer()        //加载全局系统配置
	InitCasbin()        //加载casbin
	InitTheme()         //加载全局主题
	InitLocalCache()    //local cache
	InitBase64Captcha() //Base64Captcha
	InitCrontab()       //定时任务
	InitEmailDialer()   //gomail Dialer
	InitWebsocket()     //websocket
	InitRatelimit()     //限流
	InitGoroutinePool() //初始化线程池
	InitRouter()        //初始总路由，放在最后
}
func InitializeResetAdmin() {
	global.VP = InitViper()
	global.DB = Gorm()
	service.ResetAdminPassword()
}
func InitializeUpdate() {
	global.VP = InitViper() //初始化Viper
	global.DB = Gorm()      //gorm连接数据库
	InitServer()            //加载全局系统配置
	//升级数据库casbin_rule表
	err := global.DB.Where("id > 0").Delete(&gormadapter.CasbinRule{}).Error
	if err != nil {
		global.Logrus.Error(err.Error())
		return
	}
	//插入新的数据
	InsertIntoCasbinRule()

	//升级数据库dynamic_route表
	err = global.DB.Where("id > 0").Delete(&model.DynamicRoute{}).Error
	if err != nil {
		global.Logrus.Error(err.Error())
		return
	}
	//插入新的数据
	InsertIntoDynamicRoute()
}
