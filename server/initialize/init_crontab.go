package initialize

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/robfig/cron/v3"
)

var (
	customerService *admin_logic.CustomerService
	nodeService     *admin_logic.Node
)

type FuncItem struct {
	Func  func()
	Timer string
	Des   string
}

func InitCrontab() {
	var funcs = []FuncItem{
		FuncItem{CrontabSubExpirationCheck, "0 */1 * * * *", "检查用户流量有效期定时任务"},
		FuncItem{CrontabClearTraffic, "0 0 0 */10 * *", "清理数据库(traffic)定时任务"},
		FuncItem{CrontabCheckNodeStatus, "0 */1 * * * *", "检查节点状态定时任务"},
		FuncItem{CrontabUserTrafficReset, "1 0 0 * * *", "用户流量重置定时任务"},
	}
	global.Crontab = cron.New(cron.WithSeconds())
	for _, v := range funcs {
		global.Crontab.AddFunc(v.Timer, v.Func)
		global.Logrus.Info(v.Des)
	}
	global.Crontab.Start()
}

func CrontabSubExpirationCheck() {
	err := customerService.SubExpirationCheck()
	if err != nil {
		global.Logrus.Error("服务有效性定时检查任务 error:", err)
	}
}
func CrontabUserTrafficReset() {
	err := customerService.TrafficReset()
	if err != nil {
		global.Logrus.Error("订阅流量重置任务 error:", err)
	}
}
func CrontabClearTraffic() {
	err := nodeService.ClearNodeTraffic()
	if err != nil {
		global.Logrus.Error("清理数据库(traffic)定时任务 error:", err)
	}
	err = customerService.ClearCustomerServiceTrafficLog()
	if err != nil {
		global.Logrus.Error("清理数据库(traffic)定时任务 error:", err)
	}
}
func CrontabCheckNodeStatus() {
	if !global.Server.Notice.WhenNodeOffline {
		return
	}
	text := admin_logic.GetOfflineNodeStatus()
	if text == "" {
		return
	}
	var pm admin_logic.PushMessageService
	pm.UnifiedPushMessage(text)
}
