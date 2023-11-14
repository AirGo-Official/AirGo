package service

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/utils/mail_plugin"
	"reflect"
)

// 修改系统配置
func UpdateSetting(setting *model.Server) error {
	if !reflect.DeepEqual(setting.Security, model.Security{}) {
		global.Server.Security = setting.Security

	} else if !reflect.DeepEqual(setting.Email, model.Email{}) {
		global.Server.Email = setting.Email

	} else if !reflect.DeepEqual(setting.Subscribe, model.Subscribe{}) {
		global.Server.Subscribe = setting.Subscribe
	}

	err := global.DB.Save(&global.Server).Error
	if err != nil {
		return err
	}
	//重新加载email
	d := mail_plugin.InitEmailDialer(global.Server.Email.EmailHost, int(global.Server.Email.EmailPort), global.Server.Email.EmailFrom, global.Server.Email.EmailSecret)
	if d != nil {
		global.EmailDialer = d
	}
	return nil
}
