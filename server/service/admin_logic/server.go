package admin_logic

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type System struct{}

// 修改系统配置
func (s *System) UpdateSetting(setting *model.Server) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&setting).Error
	})
	if err != nil {
		return err
	}
	global.Server = *setting
	//重新加载email
	global.GoroutinePool.Submit(func() {
		InitEmailSvc()
	})
	//重新加载tg bot
	if global.Server.Notice.BotToken != "" {
		global.GoroutinePool.Submit(func() {
			global.Logrus.Info("重新加载tg bot")
			//关闭
			TgBotSvc.TGBotCloseListen()
			//重启
			TgBotSvc.TGBotStart()
		})
	} else {
		global.GoroutinePool.Submit(func() {
			global.Logrus.Info("停止 tg bot")
			//关闭
			TgBotSvc.TGBotCloseListen()
		})
	}
	return nil
}
