package global

import (
	"github.com/casbin/casbin/v2"
	ants "github.com/panjf2000/ants/v2"
	"github.com/ppoonk/AirGo/model"
	queue "github.com/ppoonk/AirGo/utils/queue_plugin"
	"github.com/sirupsen/logrus"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"gorm.io/gorm"
)

// TODO 精简
var (
	Config        model.Config                 //全局配置（本地yaml）
	DB            *gorm.DB                     //数据库
	LocalCache    local_cache.Cache            //本地kv cache
	Casbin        *casbin.SyncedCachedEnforcer //casbin
	Server        model.Server                 //全局配置（数据库）
	Logrus        *logrus.Logger               //日志
	RateLimit     model.RateLimitRule          //限流器
	GoroutinePool *ants.Pool                   //goroutine池
	Queue         *queue.Queue                 //mini queue
)
