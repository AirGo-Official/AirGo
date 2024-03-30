package initialize

import (
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/utils/format_plugin"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

var (
	customerService *admin_logic.CustomerService
	nodeService     *admin_logic.Node
	userService     *admin_logic.User
)

type FuncItem struct {
	Func  func()
	Timer string
	Des   string
}

func InitCrontab() {
	var funcs = []FuncItem{
		FuncItem{CrontabSubExpirationCheck, "0 */1 * * * *", "已启动任务：检测用户服务有效期"},           //运行间隔：每分钟
		FuncItem{CrontabClearTraffic, "0 0 0 */10 * *", "已启动任务：清理数据库流量记录"},                //运行间隔：每10天
		FuncItem{CrontabNodeOffline, "0 */10 * * * *", "已启动任务：检查节点离线状态"},                  //运行间隔：每10分钟
		FuncItem{CrontabUserTrafficReset, "1 0 0 * * *", "已启动任务：用户订阅服务流量重置"},              //运行间隔：每天00:00:01
		FuncItem{CrontabCustomerServiceAlmostExpired, "0 0 */10 * * *", "已启动任务：用户服务到期提醒"}, //运行间隔：每10小时
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
func CrontabCustomerServiceAlmostExpired() {
	list, err := customerService.GetCustomerServiceListAlmostExpired()
	if err != nil {
		return
	}
	messageMap := make(map[int64]admin_logic.MessageInfo, 0)
	f := func(msg admin_logic.MessageInfo, v model.CustomerService) admin_logic.MessageInfo {
		msg.Message += strings.Join([]string{
			"-------------------------------------",
			fmt.Sprintf("服务名称：%s", v.Subject),
			fmt.Sprintf("服务结束时间：%s", v.ServiceEndAt.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("是否可续费：%v\n", v.IsRenew),
		}, "\n")
		if v.IsRenew {
			msg.Message += fmt.Sprintf("续费金额：%s\n", v.RenewalAmount)
		}
		if v.GoodsType == constant.GOODS_TYPE_SUBSCRIBE {
			msg.Message += strings.Join([]string{
				fmt.Sprintf("订阅状态：%v", v.SubStatus),
				fmt.Sprintf("总流量：%s GB", format_plugin.ByteToGB(v.TotalBandwidth)),
				fmt.Sprintf("已用上行：%s GB", format_plugin.ByteToGB(v.UsedUp)),
				fmt.Sprintf("已用下行：%s GB\n", format_plugin.ByteToGB(v.UsedDown)),
			}, "\n")
		}
		return msg
	}

	for _, v := range *list {
		if msg, ok := messageMap[v.UserID]; ok {
			messageMap[v.UserID] = f(msg, v)
		} else {
			user, err := userService.FirstUser(&model.User{ID: v.UserID})
			if err != nil {
				continue
			}
			if !user.WhenServiceAlmostExpired {
				continue
			}
			msg = admin_logic.MessageInfo{
				UserID:      v.UserID,
				MessageType: admin_logic.MESSAGE_TYPE_USER,
				User:        user,
				Message: strings.Join([]string{
					"【服务到期提醒】",
					fmt.Sprintf("时间：%s\n", time.Now().Format("2006-01-02 15:04:05")),
				}, "\n"),
			}
			messageMap[v.UserID] = f(msg, v)
		}
	}
	for _, v := range messageMap {
		admin_logic.PushMessageSvc.PushMessage(&v)
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
			MessageType: admin_logic.MESSAGE_TYPE_ADMIN,
			UserID:      k,
			Message:     strings.Join(text, "\n"),
		})
	}
}
