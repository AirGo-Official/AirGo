package service

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/mail_plugin"
)

// 修改系统配置
func UpdateSetting(setting *model.Server) error {
	global.Server = *setting
	err := global.DB.Save(&global.Server).Error
	if err != nil {
		return err
	}
	//重新加载email
	global.GoroutinePool.Submit(func() {
		d := mail_plugin.InitEmailDialer(global.Server.Email.EmailHost, int(global.Server.Email.EmailPort), global.Server.Email.EmailFrom, global.Server.Email.EmailSecret)
		if d != nil {
			global.EmailDialer = d
		}
	})
	//重新加载tg bot
	if global.Server.Notice.BotToken != "" {
		global.GoroutinePool.Submit(func() {
			global.Logrus.Info("重新加载tg bot")
			//关闭
			TGBotCloseListen()
			//重启
			go TGBotStartListen()
		})
	} else {
		global.GoroutinePool.Submit(func() {
			global.Logrus.Info("停止 tg bot")
			TGBotCloseListen()
		})
	}
	return nil
}
