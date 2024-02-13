package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

// 获取全部支付列表
func GetPayList(ctx *gin.Context) {
	list, _, err := common_logic.CommonSqlFind[model.Pay, string, []model.Pay]("")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("GetPayList success", list, ctx)

}

// 新建支付
func NewPay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlCreate[model.Pay](receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("NewPay success", nil, ctx)

}

// 删除支付
func DeletePay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlDelete[model.Pay, model.Pay](receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("DeletePay success", nil, ctx)

}

// 修改支付
func UpdatePay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlSave[model.Pay](receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("UpdatePay success", nil, ctx)

}
