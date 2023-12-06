package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

func NewTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail("NewTicket error:"+err.Error(), nil, ctx)
		return
	}
	if ticket.Title == "" {
		response.Fail("NewTicket error: ticker is empty", nil, ctx)
		return
	}
	uID, _ := GetUserIDFromGinContext(ctx)
	ticket.UserID = uID
	err = service.CommonSqlCreate[model.Ticket](ticket)
	if err != nil {
		response.Fail("NewTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewTicket success", nil, ctx)

}
func DeleteTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail("DeleteTicket error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlDelete[model.Ticket, model.Ticket](ticket)
	if err != nil {
		response.Fail("DeleteTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteTicket success", nil, ctx)

}
func UpdateTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail("UpdateTicket error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Ticket](ticket)
	if err != nil {
		response.Fail("UpdateTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateTicket success", nil, ctx)

}
func UpdateUserTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail("UpdateTicket error:"+err.Error(), nil, ctx)
		return
	}
	//
	uID, _ := GetUserIDFromGinContext(ctx)
	ticketParams := model.Ticket{
		ID:     ticket.ID,
		UserID: uID,
	}
	res, _, _ := service.CommonSqlFirst[model.Ticket, model.Ticket, model.Ticket](ticketParams)
	if res.ID == 0 {
		response.Fail("UpdateTicket error", nil, ctx)
		return
	}
	res.Title = ticket.Title
	res.Details = ticket.Details
	res.Status = model.TicketProcessing
	if ticket.Status == model.TicketClosed {
		res.Status = model.TicketClosed
	}
	service.CommonSqlSave[model.Ticket](res)
	response.OK("UpdateTicket success", nil, ctx)
}
func GetTicketList(ctx *gin.Context) {
	var params model.FieldParamsReq
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail("GetTicketList error:"+err.Error(), nil, ctx)
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
func GetUserTicketList(ctx *gin.Context) {
	var params model.FieldParamsReq
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail("GetUserTicketList error:"+err.Error(), nil, ctx)
		return
	}
	uID, _ := GetUserIDFromGinContext(ctx)
	params.FieldParamsList = append(params.FieldParamsList, model.FieldParamsItem{
		Operator:       "AND",
		Field:          "user_id",
		FieldType:      "",
		Condition:      "=",
		ConditionValue: fmt.Sprintf("%d", uID),
	})
	data, total, err := service.CommonSqlFindWithFieldParams(&params)
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
		response.Fail("SendTicketMessage error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlCreate[model.TicketMessage](msg)
	if err != nil {
		response.Fail("SendTicketMessage error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("SendTicketMessage success", nil, ctx)
}
func GetTicketMessage(ctx *gin.Context) {
	var params model.Ticket
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail("GetTicketMessage error:"+err.Error(), nil, ctx)
		return
	}
	var msgParams = model.TicketMessage{
		TicketID: params.ID,
	}
	data, _, err := service.CommonSqlFind[model.TicketMessage, model.TicketMessage, []model.TicketMessage](msgParams)
	if err != nil {
		response.Fail("GetTicketMessage error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetTicketMessage success", data, ctx)
}
