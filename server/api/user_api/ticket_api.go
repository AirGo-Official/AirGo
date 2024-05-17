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
	"strings"
)

// NewTicket
// @Tags [customer api] ticket
// @Summary 新建工单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Ticket true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/ticket/newTicket [post]
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
	err = service.TicketSvc.NewTicket(&ticket)
	if err != nil {
		response.Fail("NewTicket error:"+err.Error(), nil, ctx)
		return
	}
	//通知
	global.GoroutinePool.Submit(func() {
		if !global.Server.Notice.WhenNewTicket {
			return
		}
		user, err := service.UserSvc.FirstUser(&model.User{ID: uID})
		if err != nil {
			return
		}
		msg := service.MessageInfo{
			UserID:      uID,
			MessageType: service.MESSAGE_TYPE_USER,
			User:        user,
			Message: strings.Join([]string{
				"【新工单提醒】",
				fmt.Sprintf("用户ID：%d", user.ID),
				fmt.Sprintf("用户名：%s", user.UserName),
				fmt.Sprintf("工单标题：%s", ticket.Title),
				fmt.Sprintf("工单详情：%s\n", ticket.Details),
			}, "\n"),
		}
		service.PushMessageSvc.PushMessage(&msg)
	})
	response.OK("NewTicket success", nil, ctx)

}

// UpdateUserTicket
// @Tags [customer api] ticket
// @Summary 更新工单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Ticket true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/ticket/updateUserTicket [post]
func UpdateUserTicket(ctx *gin.Context) {
	var ticket model.Ticket
	err := ctx.ShouldBind(&ticket)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uID, _ := api.GetUserIDFromGinContext(ctx)
	ticket.UserID = uID
	err = service.TicketSvc.UpdateUserTicket(&ticket)
	if err != nil {
		response.Fail("UpdateTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateTicket success", nil, ctx)
}

// GetUserTicketList
// @Tags [customer api] ticket
// @Summary 获取工单列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/ticket/getUserTicketList [post]
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

// SendTicketMessage
// @Tags [customer api] ticket
// @Summary 发送工单消息
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.TicketMessage true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/ticket/sendTicketMessage [post]
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
	_, err = service.TicketSvc.FirstTicket(&model.Ticket{ID: msg.TicketID, UserID: uID})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	err = service.TicketSvc.NewTicketMessage(&msg)
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("SEND_TICKET_MESSAGE success", nil, ctx)
}

// FirstTicket
// @Tags [customer api] ticket
// @Summary 获取工单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Ticket true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/ticket/firstTicket [post]
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
	ticket, err := service.TicketSvc.FirstTicket(&model.Ticket{ID: params.ID, UserID: uID})

	if err != nil {
		response.Fail("FirstTicket error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("FirstTicket success", ticket, ctx)
}
