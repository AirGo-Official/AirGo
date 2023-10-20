package initialize

import (
	"AirGo/global"
	"AirGo/model"
)

func InitializeAll() {
	InitLogrus()            //logrus
	global.VP = InitViper() //初始化Viper
	global.DB = Gorm()      //gorm连接数据库
	if global.DB != nil {
		if !global.DB.Migrator().HasTable(&model.User{}) {
			global.Logrus.Info("未找到sys_user库表,开始建表并初始化数据...")
			RegisterTables()      //创建table
			InsertInto(global.DB) //导入数据
		} else {
			RegisterTables() //AutoMigrate 自动迁移 schema
		}
	} else {
		//os.Exit(1)
		panic("数据库连接失败")
	}
	InitServer()        //加载全局系统配置
	InitCasbin()        //加载casbin
	InitTheme()         //加载全局主题
	InitLocalCache()    //local cache
	InitBase64Captcha() //Base64Captcha
	InitCrontab()       //定时任务
	//InitAlipayClient()  //alipay
	InitEmailDialer() //gomail Dialer
	InitWebsocket()   //websocket
	InitRatelimit()   //限流
	InitRouter()      //初始总路由
}
