package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

func NewTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	if ticket.Title == "" {
		response.Fail("NewTicket error: ticker is empty", nil, ctx)
		return
	}
	uID, _ := api.GetUserIDFromGinContext(ctx)
	ticket.UserID = uID
	err = ticketService.NewTicket(&ticket)
	if err != nil {
		response.Fail("NewTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewTicket success", nil, ctx)

}
func UpdateUserTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uID, _ := api.GetUserIDFromGinContext(ctx)
	ticket.UserID = uID
	err = ticketService.UpdateUserTicket(&ticket)
	if err != nil {
		response.Fail("UpdateTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateTicket success", nil, ctx)
}
func GetUserTicketList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uID, _ := api.GetUserIDFromGinContext(ctx)
	params.FieldParamsList = append(params.FieldParamsList, model.FieldParamsItem{
		Operator:       "AND",
		Field:          "user_id",
		FieldType:      "",
		Condition:      "=",
		ConditionValue: fmt.Sprintf("%d", uID),
	})
	data, total, err := common_logic.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		response.Fail("GetUserTicketList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetUserTicketList success", model.CommonDataResp{
		Total: total,
		Data:  data,
	}, ctx)
}
func SendTicketMessage(ctx *gin.Context) {
	var msg model.TicketMessage
	err := ctx.ShouldBind(&msg)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	msg.IsAdmin = false
	uID, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR, nil, ctx)
		return
	}
	_, err = ticketService.FirstTicket(&model.Ticket{ID: msg.TicketID, UserID: uID})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	err = ticketService.NewTicketMessage(&msg)
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("SEND_TICKET_MESSAGE success", nil, ctx)
}
func FirstTicket(ctx *gin.Context) {
	var params model.Ticket
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uID, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR, nil, ctx)
		return
	}
	ticket, err := ticketService.FirstTicket(&model.Ticket{ID: params.ID, UserID: uID})

	if err != nil {
		response.Fail("FirstTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("FirstTicket success", ticket, ctx)
}
