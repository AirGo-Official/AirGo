package app

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/panjf2000/ants/v2"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/router"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/file_plugin"
	queue "github.com/ppoonk/AirGo/utils/queue_plugin"
	"github.com/sirupsen/logrus"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"github.com/yudeguang/ratelimit"
	"io"
	"os"
	"path"
	"time"
)

type App struct {
	router   *router.GinRouter
	dataBase *DataBase
}

func NewApp() *App {
	return &App{
		router:   router.NewGinRouter(),
		dataBase: NewDataBase(),
	}
}

func (a *App) InitConfig(configPath string) {
	v := viper.New()
	v.SetConfigFile(path.Join(configPath)) //config路径
	v.SetConfigType("yaml")                //设置文件的类型
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.Config); err != nil { //解析到全局配置
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

func (a *App) InitLogrus() {
	logger := logrus.New()
	src, _ := file_plugin.SetOutputFile()
	if global.Config.SystemParams.Mode == "dev" {
		logger.SetReportCaller(true)                //在输出日志中添加文件名和方法信息
		logger.Out = io.MultiWriter(src, os.Stdout) //同时打印到控制台及日志里
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.Out = src
		logger.SetLevel(logrus.InfoLevel)
	}
	//设置日志格式
	//logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	global.Logrus = logger
}

func (a *App) InitLocalCache() {
	global.LocalCache = local_cache.NewCache(
		local_cache.SetInternal(10*time.Second),               //设置哨兵时间间隔
		local_cache.SetDefaultExpire(10*time.Second),          //设置默认的超时时间
		local_cache.SetCapture(service.OrderSvc.OrderTimeout), //默认只有订单超时一个任务
	)
}

func (a *App) InitRouter() {
	a.router.InitRouter()
	a.dataBase.routes = a.router.Router.Routes()
}

func (a *App) ConnectDatabase() {
	a.dataBase.ConnectDatabase()
}

func (a *App) InitGlobalVariable() {
	// 全局配置参数
	settings, _, err := service.CommonSqlFind[model.Server, string, model.Server]("id = 1")
	if err != nil {
		global.Logrus.Error("系统配置获取失败", err.Error())
		panic(err)
	}
	global.Server = settings

	// casbin
	res, err := service.AdminCasbinSvc.NewSyncedCachedEnforcer()
	if err != nil {
		panic(err)
	}
	global.Casbin = res

	// GoroutinePool
	global.GoroutinePool, _ = ants.NewPool(4, ants.WithPreAlloc(true))
	global.GoroutinePool.Running()

	// RateLimit
	global.RateLimit.IPRole = ratelimit.NewRule()
	global.RateLimit.IPRole.AddRule(time.Second*60, int(global.Server.Security.RateLimitParams.IPRoleParam))
	global.RateLimit.VisitRole = ratelimit.NewRule()
	global.RateLimit.VisitRole.AddRule(time.Second*60, int(global.Server.Security.RateLimitParams.VisitParam))

	// queue
	global.Queue = queue.NewQueue()
	global.Queue.SetConditions(100)

	// 处理管理员通知
	service.AdminServerSvc.AdminAccountHandler()
}

func (a *App) InitTasks() {
	service.InitBase64Captcha()  //Base64Captcha
	service.InitCrontab()        //定时任务
	service.InitEmailSvc()       //邮件
	service.InitTgBotSvc()       //tg bot
	service.InitOrderSvc()       //订单
	service.InitPushMessageSvc() //推送消息
	service.InitNodeBackendSvc() //处理节点后端并发请求
}

func (a *App) Start() {
	a.router.Start()
}

func (a *App) Update() {
	err := a.dataBase.DefaultForRoleMenuCasbin()
	if err != nil {
		fmt.Println(err)
	}
	err = service.AdminServerSvc.ChangeDataForUpdate()
	if err != nil {
		fmt.Println(err)
	}
}
