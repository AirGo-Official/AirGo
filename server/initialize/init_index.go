package initialize

import (
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/mojocn/base64Captcha"
	"github.com/panjf2000/ants/v2"
	"github.com/ppoonk/AirGo/database"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/router"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/service/user_logic"
	"github.com/ppoonk/AirGo/utils/logrus_plugin"
	queue "github.com/ppoonk/AirGo/utils/queue_plugin"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/yudeguang/ratelimit"
	"time"
)

// InitializeAll 初始化全部资源，注意顺序
func InitializeAll(startConfigPath string) {
	InitLogrus()               //logrus
	InitViper(startConfigPath) //初始化Viper

	router.InitRouter()  //注册路由
	database.StartGorm() //gorm连接数据库

	InitServer()     //加载全局系统配置
	InitCasbin()     //加载casbin
	InitTheme()      //加载全局主题
	InitLocalCache() //local cache

	InitGoroutinePool()     //初始化协程池
	InitBase64Captcha()     //Base64Captcha
	InitRatelimit()         //限流
	InitCrontab()           //定时任务
	InitQueue()             //队列
	InitTask()              //初始化一些任务
	router.ListenAndServe() //启动路由监听
}

// InitializeUpdate 升级核心时的初始化
func InitializeUpdate(startConfigPath string) {
	InitLogrus()               //logrus
	InitViper(startConfigPath) //初始化Viper
	router.InitRouter()        //注册路由
	database.StartGorm()       //gorm连接数据库
	InitServer()               //加载全局系统配置
	var funcs = []func() error{
		func() error {
			fmt.Println("升级数据库casbin_rule表")
			err := global.DB.Where("id > 0").Delete(&gormadapter.CasbinRule{}).Error
			if err != nil {
				return err
			}
			return database.InsertIntoCasbinRule()
		},
		func() error {
			fmt.Println("升级角色和菜单")
			//先删除role_and_menu
			err := global.DB.Where("role_id > 0").Delete(&model.RoleAndMenu{}).Error
			if err != nil {
				return err
			}
			//再删除菜单
			err = global.DB.Where("id > 0").Delete(&model.Menu{}).Error
			if err != nil {
				return err
			}
			//插入新的菜单
			err = database.InsertIntoMenu()
			if err != nil {
				return err
			}
			//插入新的role_and_menu
			return database.InsertIntoRoleAndMenu()
		},
	}
	for _, v := range funcs {
		err := v()
		if err != nil {
			fmt.Println("升级核心出错：", err.Error())
			return
		}
	}
	fmt.Println("升级核心完成")

}

// InitializeDB 仅加载数据库
func InitializeDB(startConfigPath string) {
	InitLogrus()               //logrus
	InitViper(startConfigPath) //初始化Viper
	database.StartGorm()       //gorm连接数据库
}
func InitLogrus() {
	global.Logrus = logrus_plugin.InitLogrus()
}
func InitServer() {
	//res, err := service.GetSetting()
	res, _, err := common_logic.CommonSqlFind[model.Server, string, model.Server]("id = 1")
	if err != nil {
		global.Logrus.Error("系统配置获取失败", err.Error())
		return
	}
	global.Server = res
}
func InitCasbin() {
	var casbinService admin_logic.Casbin
	res, err := casbinService.NewSyncedCachedEnforcer()
	if err != nil {
		panic(err)
	}
	global.Casbin = res
}
func InitTheme() {
	//res, err := service.GetThemeConfig()
	res, _, err := common_logic.CommonSqlFind[model.Theme, string, model.Theme]("id = 1")
	if err != nil {
		global.Logrus.Error("系统配置获取失败", err.Error())
		return
	}
	global.Theme = res
}
func InitLocalCache() {
	//初始化local cache配置
	var order *user_logic.Order
	global.LocalCache = local_cache.NewCache(
		local_cache.SetInternal(10*time.Second),      //设置哨兵时间间隔
		local_cache.SetDefaultExpire(10*time.Second), //设置默认的超时时间
		local_cache.SetCapture(order.OrderTimeout),
	)
}
func InitBase64Captcha() {
	// base64Captcha.DefaultMemStore 是默认的过期时间10分钟。也可以自己设定参数 base64Captcha.NewMemoryStore(GCLimitNumber, Expiration)
	global.Base64CaptchaStore = base64Captcha.DefaultMemStore
	driver := base64Captcha.NewDriverDigit(38, 120, 4, 0.2, 10)
	global.Base64Captcha = base64Captcha.NewCaptcha(driver, global.Base64CaptchaStore)
}
func InitRatelimit() {
	global.RateLimit.IPRole = ratelimit.NewRule()
	global.RateLimit.IPRole.AddRule(time.Second*60, int(global.Server.Security.RateLimitParams.IPRoleParam))
	global.RateLimit.VisitRole = ratelimit.NewRule()
	global.RateLimit.VisitRole.AddRule(time.Second*60, int(global.Server.Security.RateLimitParams.VisitParam))
}
func InitGoroutinePool() {
	global.GoroutinePool, _ = ants.NewPool(4, ants.WithPreAlloc(true))
}
func InitQueue() {
	global.Queue = queue.NewQueue()
	global.Queue.SetConditions(100)
}
func InitTask() {
	admin_logic.InitEmailSvc() //邮件
	//service.InitTgBotSvc() //tg bot
	user_logic.InitOrderSvc() //
}
