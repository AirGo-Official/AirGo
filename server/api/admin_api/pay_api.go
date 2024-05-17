package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetPayList
// @Tags [admin api] pay
// @Summary 获取全部支付列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/pay/getPayList [get]
func GetPayList(ctx *gin.Context) {
	list, _, err := service.CommonSqlFind[model.Pay, string, []model.Pay]("")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("GetPayList success", list, ctx)

}

// NewPay
// @Tags [admin api] pay
// @Summary 新建支付
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Pay true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/pay/newPay [post]
func NewPay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlCreate[model.Pay](receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("NewPay success", nil, ctx)

}

// DeletePay
// @Tags [admin api] pay
// @Summary 删除支付
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Pay true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/pay/deletePay [delete]
func DeletePay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlDelete[model.Pay, model.Pay](receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("DeletePay success", nil, ctx)

}

// UpdatePay
// @Tags [admin api] pay
// @Summary 修改支付
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Pay true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/pay/updatePay [post]
func UpdatePay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Pay](receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("UpdatePay success", nil, ctx)

}
