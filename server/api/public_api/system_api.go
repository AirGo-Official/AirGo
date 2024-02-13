package public_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/response"
)

// 主题配置
func GetThemeConfig(ctx *gin.Context) {
	response.OK("GetThemeConfig success", global.Theme, ctx)
}

// 获取公共系统设置
func GetPublicSetting(ctx *gin.Context) {
	var ps = model.PublicSystem{
		EnableRegister:          global.Server.Subscribe.EnableRegister,
		EnableEmailCode:         global.Server.Subscribe.EnableEmailCode,
		EnableLoginEmailCode:    global.Server.Subscribe.EnableLoginEmailCode,
		BackendUrl:              global.Server.Subscribe.BackendUrl,
		EnabledClockIn:          global.Server.Subscribe.EnabledClockIn,
		AcceptableEmailSuffixes: global.Server.Subscribe.AcceptableEmailSuffixes,
	}
	response.OK("GetPublicSetting success", ps, ctx)
}
