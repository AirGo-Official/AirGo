package initialize

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/robfig/cron/v3"
	"strconv"
	"strings"
	"time"
)

// 秒级操作
// c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithLogger(
// 	cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

// 添加一个任务
//_, err := c.AddFunc("*/5 * * * * *", func() {
// 	fmt.Println("开始定时任务")
// })
//c.Start()
// 允许往正在执行的 cron 中添加任务
//c.AddFunc("@daily", func() { fmt.Println("Every day") })
// 检查上一个和下一个任务执行的时间
//inspect(c.Entries())
//c.Stop()  // 停止调度，但正在运行的作业不会被停止

func InitCrontab() {
	global.Crontab = cron.New(cron.WithSeconds())

	// 用户流量，有效期任务,默认2分钟
	_, err := global.Crontab.AddFunc("0 */2 * * * *", func() {
		err := service.UserExpiryCheck()
		if err != nil {
			global.Logrus.Error("用户流量有效期定时任务 error:", err)
		}
	})
	if err != nil {
		global.Logrus.Error("用户流量有效期定时任务 error:", err)
	} else {
		global.Logrus.Info("用户流量有效期定时任务")
	}
	// 定时清理数据库(traffic),默认10天
	_, err = global.Crontab.AddFunc("0 0 0 */10 * *", func() {
		y, m, _ := time.Now().Date()
		startTime := time.Date(y, m-2, 1, 0, 0, 0, 0, time.Local)
		err := global.DB.Where("created_at < ?", startTime).Delete(&model.TrafficLog{}).Error
		if err != nil {
			global.Logrus.Error("清理数据库(traffic)定时任务 error:", err)
		}
	})
	if err != nil {
		global.Logrus.Error("清理数据库(traffic)定时任务 error:", err)
	} else {
		global.Logrus.Info("清理数据库(traffic)定时任务")
	}
	// 定时检查节点状态并通知,默认2分钟
	_, err = global.Crontab.AddFunc("0 */2 * * * *", func() {
		if !global.Server.Notice.WhenNodeOffline {
			return
		}
		text := service.GetOfflineNodeStatus()
		if text == "" {
			return
		}
		if global.Server.Notice.TGAdmin == "" {
			return
		}
		tgIDs := strings.Fields(global.Server.Notice.TGAdmin)
		for _, v := range tgIDs {
			chatID, _ := strconv.ParseInt(v, 10, 64)
			service.TGBotSendMessage(chatID, text)
		}

	})
	if err != nil {
		global.Logrus.Error("检查节点状态并定时任务 error:", err)
	} else {
		global.Logrus.Info("检查节点状态并定时任务")
	}
	global.Crontab.Start()

}
