package public_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetThemeConfig
// @Tags [public api] server
// @Summary 获取主题
// @Produce json
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/public/server/getThemeConfig [get]
func GetThemeConfig(ctx *gin.Context) {
	if cache, ok := global.LocalCache.Get(constant.CACHE_THEME); ok {
		response.OK("GetThemeConfig success", cache, ctx)
		return
	}
	theme, _, err := service.CommonSqlFirst[model.Theme, model.Theme, model.Theme](model.Theme{ID: 1})
	if err != nil {
		response.Fail("GetThemeConfig error: "+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.SetNoExpire(constant.CACHE_THEME, theme)
	response.OK("GetThemeConfig success", theme, ctx)

}

// GetPublicSetting
// @Tags [public api] server
// @Summary 获取公共系统设置
// @Produce json
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/public/server/getPublicSetting [get]
func GetPublicSetting(ctx *gin.Context) {
	var ps = model.PublicSystem{
		EnableRegister:          global.Server.Website.EnableRegister,
		AcceptableEmailSuffixes: global.Server.Website.AcceptableEmailSuffixes,
		EnableBase64Captcha:     global.Server.Website.EnableBase64Captcha,
		EnableEmailCode:         global.Server.Website.EnableEmailCode,
		EnableLoginEmailCode:    global.Server.Website.EnableLoginEmailCode,
		BackendUrl:              global.Server.Subscribe.BackendUrl,
		CommissionRate:          global.Server.Finance.CommissionRate,
		WithdrawThreshold:       global.Server.Finance.WithdrawThreshold,
		EnableLottery:           global.Server.Finance.EnableLottery,
		Jackpot:                 global.Server.Finance.Jackpot,
		SubName:                 global.Server.Subscribe.SubName,
	}
	response.OK("GetPublicSetting success", ps, ctx)
}
