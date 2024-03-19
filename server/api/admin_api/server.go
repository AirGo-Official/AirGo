package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

// 更新主题
func UpdateThemeConfig(ctx *gin.Context) {
	var theme model.Theme
	err := ctx.ShouldBind(&theme)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlSave[model.Theme](theme)
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
	res, _, err := common_logic.CommonSqlFind[model.Server, string, model.Server]("id = 1")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetSetting error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetSetting success", res, ctx)
}

// 更新系统设置
func UpdateSetting(ctx *gin.Context) {
	var setting model.Server
	err := ctx.ShouldBind(&setting)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = systemService.UpdateSetting(&setting)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateSetting error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateSetting success", nil, ctx)
}

func GetCurrentVersion(ctx *gin.Context) {
	response.OK("GetCurrentVersion success", gin.H{"version": constant.V}, ctx)
}

func GetLatestVersion(ctx *gin.Context) {
	v, err := systemService.GetLatestVersion()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetLatestVersion error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetLatestVersion success", gin.H{"version": v}, ctx)
}

func UpdateLatestVersion(ctx *gin.Context) {
	api.SSE(ctx)
	err := systemService.DownloadLatestVersion(ctx)
	if err != nil {
		response.ResponseSSE("message error", err.Error(), ctx)
		return
	}
	response.ResponseSSE("success", "success", ctx)
}
