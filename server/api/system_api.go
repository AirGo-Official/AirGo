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
	response.OK("主题获取成功", global.Theme, ctx)

}

// 更新主题
func UpdateThemeConfig(ctx *gin.Context) {
	var theme model.Theme
	err := ctx.ShouldBind(&theme)
	if err != nil {
		global.Logrus.Error("主题设置参数错误:", err)
		response.Fail("主题设置参数错误"+err.Error(), nil, ctx)
		return
	}
	//err = service.UpdateThemeConfig(&theme)
	err = service.CommonSqlSave[model.Theme](theme)
	if err != nil {
		global.Logrus.Error("设置主题 error:", err)
		response.Fail("主题设置错误"+err.Error(), nil, ctx)
		return
	}
	//重新加载主题设置
	global.Theme = theme
	response.OK("主题设置成功", nil, ctx)

}

// 获取系统设置
func GetSetting(ctx *gin.Context) {
	res, _, err := service.CommonSqlFind[model.Server, string, model.Server]("id = 1")

	if err != nil {
		global.Logrus.Error("系统设置获取错误:", err.Error())
		response.Fail("系统设置获取错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("系统设置获取成功", res, ctx)

}

// 获取公共系统设置
func GetPublicSetting(ctx *gin.Context) {

	//res, err := service.GetPublicSetting()

	var ps = model.PublicSystem{
		EnableRegister:       global.Server.System.EnableRegister,
		EnableEmailCode:      global.Server.System.EnableEmailCode,
		EnableLoginEmailCode: global.Server.System.EnableLoginEmailCode,
		RebateRate:           global.Server.System.RebateRate,
		BackendUrl:           global.Server.System.BackendUrl,
	}

	response.OK("系统设置获取成功", ps, ctx)
}

func UpdateSetting(ctx *gin.Context) {
	var setting model.Server
	err := ctx.ShouldBind(&setting)
	if err != nil {
		global.Logrus.Error("更改系统设置参数错误:", err.Error())
		response.Fail("更改系统设置参数错误"+err.Error(), nil, ctx)
		return
	}
	err = service.UpdateSetting(&setting)
	if err != nil {
		global.Logrus.Error("更改系统设置错误:", err.Error())
		response.Fail("更改系统设置错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("更改系统设置成功", nil, ctx)

}
