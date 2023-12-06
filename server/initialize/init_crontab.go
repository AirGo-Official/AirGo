package initialize

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/robfig/cron/v3"
	"time"
)

type FuncItem struct {
	Func  func()
	Timer string
	Des   string
}

func InitCrontab() {
	var funcs = []FuncItem{
		FuncItem{CrontabUserExpiration, "0 */2 * * * *", "检查用户流量有效期定时任务"},
		FuncItem{CrontabClearTraffic, "0 0 0 */10 * *", "清理数据库(traffic)定时任务"},
		FuncItem{CrontabCheckNodeStatus, "0 */2 * * * *", "检查节点状态定时任务"},
		FuncItem{CrontabUserTrafficReset, "1 0 0 * * *", "用户流量重置定时任务"},
	}
	global.Crontab = cron.New(cron.WithSeconds())
	for _, v := range funcs {
		global.Crontab.AddFunc(v.Timer, v.Func)
		global.Logrus.Info(v.Des)
	}
	global.Crontab.Start()
}

func CrontabUserExpiration() {
	err := service.UserExpiryCheck()
	if err != nil {
		global.Logrus.Error("用户流量有效期定时任务 error:", err)
	}
}
func CrontabClearTraffic() {
	y, m, _ := time.Now().Date()
	startTime := time.Date(y, m-2, 1, 0, 0, 0, 0, time.Local)
	err := global.DB.Where("created_at < ?", startTime).Delete(&model.TrafficLog{}).Error
	if err != nil {
		global.Logrus.Error("清理数据库(traffic)定时任务 error:", err)
	}
}
func CrontabCheckNodeStatus() {
	if !global.Server.Notice.WhenNodeOffline {
		return
	}
	text := service.GetOfflineNodeStatus()
	if text == "" {
		return
	}
	service.UnifiedPushMessage(text)
}

func CrontabUserTrafficReset() {
	err := service.UserTrafficReset()
	if err != nil {
		global.Logrus.Error("用户流量重置任务 error:", err)
	}

}
func CrontabOnlineUsers() {

}
