package global

import (
	"github.com/casbin/casbin/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mojocn/base64Captcha"
	ants "github.com/panjf2000/ants/v2"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/websocket_plugin"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	DB                 *gorm.DB                        //数据库
	Config             model.Config                    //全局配置（本地yaml）
	VP                 *viper.Viper                    //
	LocalCache         local_cache.Cache               //本地kv cache
	Casbin             *casbin.CachedEnforcer          //casbin
	Server             model.Server                    //全局配置（数据库）
	Theme              model.Theme                     //全局主题配置
	Base64Captcha      *base64Captcha.Captcha          //Base64Captcha
	Base64CaptchaStore base64Captcha.Store             //Base64CaptchaStore
	EmailDialer        *gomail.Dialer                  //Email
	WsManager          *websocket_plugin.ClientManager //websocket 全局客户端管理
	Logrus             *logrus.Logger                  //日志
	RateLimit          model.RateLimitRule             //限流器
	GoroutinePool      *ants.Pool                      //线程池
	Crontab            *cron.Cron                      //定时任务
	TGBot              *tgbotapi.BotAPI                //tg bot
	ContextGroup       *model.ContextGroup             //
)

const (
	CtxSetUserName = "UserName"
	CtxSetUserID   = "UserID"
	//
	UserRouteList   = "UserRouteList"
	AllRouteList    = "AllRouteList"
	AllEnabledGoods = "AllEnabledGoods"
	//
	UserRoleIds = "UserRoleIds"
	//
	NodeStatus           = "status"
	NodeStatusIsNotified = "notified"
	//tg 常量
	TGBotCtx    = "TGBotCtx"
	TGBotCancel = "TGBotCancel"

	//节点类型
	NodeTypeVmess       = "vmess"
	NodeTypeVless       = "vless"
	NodeTypeTrojan      = "trojan"
	NodeTypeHysteria    = "hysteria"
	NodeTypeShadowsocks = "shadowsocks"

	//传输协议
	NetworkWs   = "ws"
	NetworkTcp  = "tcp"
	NetworkKcp  = "kcp"
	NetworkGrpc = "grpc"
	NetworkQuic = "quic"
)
