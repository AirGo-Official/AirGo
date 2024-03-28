package initialize

import (
	"fmt"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
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
		FuncItem{CrontabSubExpirationCheck, "0 */1 * * * *", "已启动任务：检测用户服务有效期"},
		FuncItem{CrontabClearTraffic, "0 0 0 */10 * *", "已启动任务：清理数据库流量记录"},
		FuncItem{CrontabNodeOffline, "0 */5 * * * *", "已启动任务：检查节点离线状态"},
		FuncItem{CrontabUserTrafficReset, "1 0 0 * * *", "已启动任务：用户订阅服务流量重置"},
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
		global.Logrus.Error("检测用户服务有效期 error:", err)
	}
}
func CrontabUserTrafficReset() {
	err := customerService.TrafficReset()
	if err != nil {
		global.Logrus.Error("用户订阅服务流量重置 error:", err)
	}
}
func CrontabClearTraffic() {
	err := nodeService.ClearNodeTraffic()
	if err != nil {
		global.Logrus.Error("清理数据库流量记录 error:", err)
	}
	err = customerService.ClearCustomerServiceTrafficLog()
	if err != nil {
		global.Logrus.Error("清理数据库流量记录 error:", err)
	}
}
func CrontabNodeOffline() {
	if !global.Server.Notice.WhenNodeOffline {
		return
	}
	var text = []string{
		"【节点离线通知】",
		"时间：" + time.Now().Format("2006-01-02 15:04:05"),
	}
	list := nodeService.GetNodesStatus()
	if len(*list) == 0 {
		return
	}
	for _, v := range *list {
		if !v.Status {
			text = append(text, fmt.Sprintf("id: %d name: %s", v.ID, v.Name))
		}
	}
	for k, _ := range global.Server.Notice.AdminIDCache {
		admin_logic.PushMessageSvc.PushMessage(&admin_logic.MessageInfo{
			UserID:  k,
			Message: strings.Join(text, "\n"),
		})
	}
}
