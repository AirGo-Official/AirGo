package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetBalanceStatementList
// @Tags [customer api] finance
// @Summary 获取余额流水
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/finance/getBalanceStatementList [post]
func GetBalanceStatementList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	params.FieldParamsList = append(params.FieldParamsList, model.FieldParamsItem{
		//Operator:       "", //前端只查分页，不传其他参数 不用填
		Field: "user_id",
		//FieldType:      "",
		Condition:      "=",
		ConditionValue: fmt.Sprintf("%d", uIDInt),
	})
	var data model.CommonDataResp
	list, total, err := service.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetBalanceStatementList error:"+err.Error(), nil, ctx)
		return
	}
	data.Data = list
	data.Total = total

	response.OK("GetBalanceStatementList success", data, ctx)
}

// GetCommissionStatementList
// @Tags [customer api] finance
// @Summary 获取佣金流水
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/finance/getCommissionStatementList [post]
func GetCommissionStatementList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	params.FieldParamsList = append(params.FieldParamsList, model.FieldParamsItem{
		//Operator:       "", //前端只查分页，不传其他参数 不用填
		Field: "user_id",
		//FieldType:      "",
		Condition:      "=",
		ConditionValue: fmt.Sprintf("%d", uIDInt),
	})
	var data model.CommonDataResp
	list, total, err := service.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetCommissionStatementList error:"+err.Error(), nil, ctx)
		return
	}
	data.Data = list
	data.Total = total
	response.OK("GetCommissionStatementList success", data, ctx)
}

// WithdrawToBalance
// @Tags [customer api] finance
// @Summary 提现
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/finance/withdrawToBalance [get]
func WithdrawToBalance(ctx *gin.Context) {
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	err := service.FinanceSvc.WithdrawToBalance(uIDInt)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("WithdrawToBalance error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("WithdrawToBalance success", nil, ctx)
}

// GetCommissionSummary
// @Tags [customer api] finance
// @Summary 获取佣金统计
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/finance/getCommissionSummary [get]
func GetCommissionSummary(ctx *gin.Context) {
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	data, err := service.FinanceSvc.GetCommissionSummary(uIDInt)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetCommissionSummary error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetCommissionSummary success", data, ctx)

}

// GetInvitationUserList
// @Tags [customer api] finance
// @Summary 获取邀请人数
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/finance/getInvitationUserList [post]
func GetInvitationUserList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	params.FieldParamsList = append(params.FieldParamsList, model.FieldParamsItem{
		//Operator:       "", //前端只查分页，不传其他参数 不用填
		Field: "referrer_user_id",
		//FieldType:      "",
		Condition:      "=",
		ConditionValue: fmt.Sprintf("%d", uIDInt),
	})
	var data model.CommonDataResp
	list, total, err := service.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetInvitationUserList error:"+err.Error(), nil, ctx)
		return
	}
	data.Data = list
	data.Total = total
	response.OK("GetInvitationUserList success", data, ctx)
}
