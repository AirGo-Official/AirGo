package global

import (
	"github.com/casbin/casbin/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mojocn/base64Captcha"
	ants "github.com/panjf2000/ants/v2"
	"github.com/ppoonk/AirGo/model"
	queue "github.com/ppoonk/AirGo/utils/queue_plugin"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"gorm.io/gorm"
)

var (
	//Viper              *viper.Viper
	Config             model.Config                 //全局配置（本地yaml）
	DB                 *gorm.DB                     //数据库
	LocalCache         local_cache.Cache            //本地kv cache
	Casbin             *casbin.SyncedCachedEnforcer //casbin
	Server             model.Server                 //全局配置（数据库）
	Theme              model.Theme                  //全局主题配置
	Base64Captcha      *base64Captcha.Captcha       //Base64Captcha
	Base64CaptchaStore base64Captcha.Store          //Base64CaptchaStore
	Logrus             *logrus.Logger               //日志
	RateLimit          model.RateLimitRule          //限流器
	GoroutinePool      *ants.Pool                   //goroutine池
	Crontab            *cron.Cron                   //定时任务
	TGBot              *tgbotapi.BotAPI             //tg bot
	Queue              *queue.Queue                 //mini queue
)
