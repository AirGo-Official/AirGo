package global

import (
	"AirGo/model"
	"AirGo/utils/websocket_plugin"
	"github.com/casbin/casbin/v2"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	Config     model.Config      //全局配置（本地yaml）
	VP         *viper.Viper      //
	LocalCache local_cache.Cache //本地kv cache

	Casbin             *casbin.CachedEnforcer          //casbin
	Server             model.Server                    //全局配置（数据库）
	Theme              model.Theme                     //全局主题配置
	Base64Captcha      *base64Captcha.Captcha          //Base64Captcha
	Base64CaptchaStore base64Captcha.Store             //Base64CaptchaStore
	EmailDialer        *gomail.Dialer                  //Email
	WsManager          *websocket_plugin.ClientManager //websocket 全局客户端管理
	Logrus             *logrus.Logger                  //日志
	RateLimit          model.RateLimitRule             //限流
	//Base64CaptchaDriver *base64Captcha.DriverDigit //Base64CaptchaDriver
)
