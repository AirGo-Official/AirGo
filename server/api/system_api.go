package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/response"
	"github.com/gin-gonic/gin"
)

// 主题配置
func GetThemeConfig(ctx *gin.Context) {
	response.OK("GetThemeConfig success", global.Theme, ctx)
}

// 更新主题
func UpdateThemeConfig(ctx *gin.Context) {
	var theme model.Theme
	err := ctx.ShouldBind(&theme)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateThemeConfig error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Theme](theme)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateThemeConfig error:"+err.Error(), nil, ctx)
		return
	}
	//重新加载主题设置
	global.Theme = theme
	response.OK("UpdateThemeConfig success", nil, ctx)
}

// 获取系统设置
func GetSetting(ctx *gin.Context) {
	res, _, err := service.CommonSqlFind[model.Server, string, model.Server]("id = 1")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetSetting error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetSetting success", res, ctx)
}

// 获取公共系统设置
func GetPublicSetting(ctx *gin.Context) {
	var ps = model.PublicSystem{
		EnableRegister:       global.Server.System.EnableRegister,
		EnableEmailCode:      global.Server.System.EnableEmailCode,
		EnableLoginEmailCode: global.Server.System.EnableLoginEmailCode,
		RebateRate:           global.Server.System.RebateRate,
		BackendUrl:           global.Server.System.BackendUrl,
		EnabledClockIn:       global.Server.System.EnabledClockIn,
	}
	response.OK("GetPublicSetting success", ps, ctx)
}

func UpdateSetting(ctx *gin.Context) {
	var setting model.Server
	err := ctx.ShouldBind(&setting)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateSetting error:"+err.Error(), nil, ctx)
		return
	}
	err = service.UpdateSetting(&setting)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateSetting error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateSetting success", nil, ctx)
}
