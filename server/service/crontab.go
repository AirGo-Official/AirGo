package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AirGo-Official/AirGo/constant"
	"github.com/AirGo-Official/AirGo/global"
	"github.com/AirGo-Official/AirGo/model"
	"github.com/AirGo-Official/AirGo/utils/format_plugin"
	"github.com/robfig/cron/v3"
)

type Crontab struct {
	Crontab  *cron.Cron
	TasksMap map[string]cron.EntryID
	Lock     sync.RWMutex
}

type Task struct {
	Func  func()
	Timer string
	Name  string
	Des   string
}

var CrontabSvc *Crontab

func InitCrontab() {
	CrontabSvc = newCrontab()
	var ts = []Task{
		{subExpirationCheck, "0 */1 * * * *", "subExpirationCheck", "已启动任务：检测用户服务有效期"},                     //运行间隔：每分钟
		{clearTraffic, "0 0 0 */10 * *", "clearTraffic", "已启动任务：清理数据库流量记录"},                                //运行间隔：每10天
		{nodeOffline, "0 */10 * * * *", "nodeOffline", "已启动任务：检查节点离线状态"},                                   //运行间隔：每10分钟
		{userTrafficReset, "1 0 0 * * *", "userTrafficReset", "已启动任务：用户订阅服务流量重置"},                          //运行间隔：每天00:00:01
		{customerServiceAlmostExpired, "0 0 */10 * * *", "customerServiceAlmostExpired", "已启动任务：用户服务到期提醒"}, //运行间隔：每10小时
	}
	CrontabSvc.AddTask(ts...)
	CrontabSvc.Start()
}
func newCrontab() *Crontab {
	return &Crontab{
		Crontab:  cron.New(cron.WithSeconds()),
		TasksMap: make(map[string]cron.EntryID),
	}
}

func (c *Crontab) Start() {
	c.Crontab.Start()
}
func (c *Crontab) Stop() {
	c.Crontab.Stop()
}
func (c *Crontab) AddTask(t ...Task) {
	if len(t) < 0 {
		c.logErr(errors.New("illegal parameter"))
		return
	}
	for _, v := range t {
		id, err := c.Crontab.AddFunc(v.Timer, v.Func)
		if err != nil {
			c.logErr(err)
			continue
		}
		c.Lock.Lock()
		c.TasksMap[v.Name] = id
		c.Lock.Unlock()
		c.log(v.Des)
	}
}
func (c *Crontab) Remove(name string) {
	c.Lock.RLock()
	id, ok := c.TasksMap[name]
	c.Lock.RUnlock()
	if !ok {
		return
	}
	c.Crontab.Remove(id)
	c.Lock.Lock()
	delete(c.TasksMap, name)
	c.Lock.Unlock()
}
func (c *Crontab) log(msg string) {
	global.Logrus.Info(msg)
}
func (c *Crontab) logErr(err error) {
	global.Logrus.Error(err)
}

func subExpirationCheck() {
	err := AdminCustomerServiceSvc.SubExpirationCheck()
	if err != nil {
		global.Logrus.Error("检测用户服务有效期 error:", err)
	}
}
func customerServiceAlmostExpired() {
	list, err := AdminCustomerServiceSvc.GetCustomerServiceListAlmostExpired()
	if err != nil {
		return
	}
	messageArray := make([]MessageInfo, 0)
	f := func(msg MessageInfo, v model.CustomerService) MessageInfo {
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
		user, err := AdminUserSvc.FirstUser(&model.User{ID: v.UserID})
		if err != nil {
			continue
		}
		if !user.WhenServiceAlmostExpired {
			continue
		}
		msg := MessageInfo{
			UserID:      v.UserID,
			MessageType: MESSAGE_TYPE_USER,
			User:        user,
			Message: strings.Join([]string{
				"【服务到期提醒】",
				fmt.Sprintf("时间：%s\n", time.Now().Format("2006-01-02 15:04:05")),
			}, "\n"),
		}
		messageArray = append(messageArray, f(msg, v))
	}
	for _, v := range messageArray {
		PushMessageSvc.PushMessage(&v)
	}
}
func userTrafficReset() {
	err := AdminCustomerServiceSvc.TrafficReset()
	if err != nil {
		global.Logrus.Error("用户订阅服务流量重置 error:", err)
	}
}
func clearTraffic() {
	err := AdminNodeSvc.ClearNodeTraffic()
	if err != nil {
		global.Logrus.Error("清理数据库流量记录 error:", err)
	}
	err = AdminCustomerServiceSvc.ClearCustomerServiceTrafficLog()
	if err != nil {
		global.Logrus.Error("清理数据库流量记录 error:", err)
	}
}
func nodeOffline() {
	if !global.Server.Notice.WhenNodeOffline {
		return
	}
	var text = []string{
		"【节点离线通知】",
		"时间：" + time.Now().Format("2006-01-02 15:04:05"),
	}
	list := AdminNodeSvc.GetNodesStatus()
	if len(*list) == 0 {
		return
	}
	count := 0
	for _, v := range *list {
		if !v.Status {
			text = append(text, fmt.Sprintf("id: %d name: %s", v.ID, v.Name))
			count++
		}
	}
	if count > 0 {
		for k, _ := range global.Server.Notice.AdminIDCache {
			PushMessageSvc.PushMessage(&MessageInfo{
				MessageType: MESSAGE_TYPE_ADMIN,
				UserID:      k,
				Message:     strings.Join(text, "\n"),
			})
		}
	}
}
