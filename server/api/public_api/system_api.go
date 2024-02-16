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
		EnableRegister:          global.Server.Website.EnableRegister,
		EnableEmailCode:         global.Server.Website.EnableEmailCode,
		EnableLoginEmailCode:    global.Server.Website.EnableLoginEmailCode,
		BackendUrl:              global.Server.Website.BackendUrl,
		EnabledClockIn:          global.Server.Website.EnabledClockIn,
		AcceptableEmailSuffixes: global.Server.Website.AcceptableEmailSuffixes,
	}
	response.OK("GetPublicSetting success", ps, ctx)
}
