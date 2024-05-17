package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// DeleteTicket
// @Tags [admin api] ticket
// @Summary 删除工单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Ticket true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/ticket/deleteTicket [delete]
func DeleteTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlDelete[model.Ticket, model.Ticket](ticket)
	if err != nil {
		response.Fail("DeleteTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteTicket success", nil, ctx)

}

// UpdateTicket
// @Tags [admin api] ticket
// @Summary 更新工单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Ticket true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/ticket/updateTicket [post]
func UpdateTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Ticket](ticket)
	if err != nil {
		response.Fail("UpdateTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateTicket success", nil, ctx)

}

// GetTicketList
// @Tags [admin api] ticket
// @Summary 获取工单列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/ticket/getTicketList [post]
func GetTicketList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	data, total, err := service.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		response.Fail("GetTicketList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetTicketList success", model.CommonDataResp{
		Total: total,
		Data:  data,
	}, ctx)
}

// FirstTicket
// @Tags [admin api] ticket
// @Summary 获取工单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Ticket true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/ticket/firstTicket [post]
func FirstTicket(ctx *gin.Context) {
	var params model.Ticket
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	ticket, err := service.AdminTicketSvc.FirstTicket(&model.Ticket{ID: params.ID})

	if err != nil {
		response.Fail("FirstTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("FirstTicket success", ticket, ctx)
}

// SendTicketMessage
// @Tags [admin api] ticket
// @Summary 发送工单消息
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.TicketMessage true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/ticket/sendTicketMessage [post]
func SendTicketMessage(ctx *gin.Context) {
	var msg model.TicketMessage
	err := ctx.ShouldBind(&msg)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	msg.IsAdmin = true
	err = service.AdminTicketSvc.NewTicketMessage(&msg)
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("SendTicketMessage success", nil, ctx)
}
