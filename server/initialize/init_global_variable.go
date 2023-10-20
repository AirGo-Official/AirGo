package initialize

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/logrus_plugin"
	"AirGo/utils/mail_plugin"
	"AirGo/utils/time_plugin"
	"AirGo/utils/websocket_plugin"
	"github.com/mojocn/base64Captcha"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/yudeguang/ratelimit"
	"time"
)

func InitBase64Captcha() {
	// base64Captcha.DefaultMemStore 是默认的过期时间10分钟。也可以自己设定参数 base64Captcha.NewMemoryStore(GCLimitNumber, Expiration)
	global.Base64CaptchaStore = base64Captcha.DefaultMemStore
	driver := base64Captcha.NewDriverDigit(38, 120, 4, 0.2, 10)
	global.Base64Captcha = base64Captcha.NewCaptcha(driver, global.Base64CaptchaStore)
}

func InitLogrus() {
	global.Logrus = logrus_plugin.InitLogrus()
}

func InitTheme() {
	//res, err := service.GetThemeConfig()
	res, _, err := service.CommonSqlFind[model.Theme, string, model.Theme]("id = 1")
	if err != nil {
		global.Logrus.Error("系统配置获取失败", err.Error())
		return
	}
	global.Theme = res
}

// 系统配置
func InitServer() {
	//res, err := service.GetSetting()
	res, _, err := service.CommonSqlFind[model.Server, string, model.Server]("id = 1")
	if err != nil {
		global.Logrus.Error("系统配置获取失败", err.Error())
		return
	}
	global.Server = res
}
func InitWebsocket() {
	global.WsManager = websocket_plugin.NewManager()
	global.WsManager.NewClientManager()
}

func InitRatelimit() {
	global.RateLimit.IPRole = ratelimit.NewRule()
	global.RateLimit.IPRole.AddRule(time.Second*60, int(global.Server.RateLimitParams.IPRoleParam))
	global.RateLimit.VisitRole = ratelimit.NewRule()
	global.RateLimit.VisitRole.AddRule(time.Second*60, int(global.Server.RateLimitParams.VisitParam))
}
func InitEmailDialer() {
	d := mail_plugin.InitEmailDialer(global.Server.Email.EmailHost, int(global.Server.Email.EmailPort), global.Server.Email.EmailFrom, global.Server.Email.EmailSecret)
	if d != nil {
		global.EmailDialer = d
	}
}
func InitLocalCache() {
	//判断有没有设置时间
	dr := time.Hour
	if global.Server.JWT.ExpiresTime != "" {
		dr, _ = time_plugin.ParseDuration(global.Server.JWT.ExpiresTime)
	}
	//初始化local cache配置
	global.LocalCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr), //设置默认的超时时间
	)
}
func InitCasbin() {
	global.Casbin = service.Casbin()
}
