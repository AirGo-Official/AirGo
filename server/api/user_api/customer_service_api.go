package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/response"
	uuid "github.com/satori/go.uuid"
)

func GetCustomerServiceList(ctx *gin.Context) {
	uID, _ := api.GetUserIDFromGinContext(ctx)
	csArr, err := customerService.GetCustomerServiceList(&model.CustomerService{UserID: uID, ServiceStatus: true})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("Success", csArr, ctx)
}
func ResetSubscribeUUID(ctx *gin.Context) {
	var cs model.CustomerService
	err := ctx.ShouldBind(&cs)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	if cs.ID == 0 {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR, nil, ctx)
		return
	}
	if cs.SubUUID.String() == "" {
		cs.SubUUID = uuid.NewV4()
	}
	err = customerService.UpdateCustomerService(cs.ID, map[string]any{
		"sub_uuid": cs.SubUUID,
	})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("Success", nil, ctx)
}
func PushCustomerService(ctx *gin.Context) {
	var cs model.PushCustomerServiceRequest
	err := ctx.ShouldBind(&cs)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uID, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR, nil, ctx)
		return
	}
	toUser, err := userService.FirstUser(&model.User{UserName: cs.ToUserName})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	res, err := customerService.FirstCustomerService(&model.CustomerService{ID: cs.CustomerServiceID, UserID: uID, IsRenew: true, ServiceStatus: true})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	err = customerService.UpdateCustomerService(res.ID, map[string]any{
		"user_id":   toUser.ID,
		"user_name": toUser.UserName,
		"sub_uuid":  uuid.NewV4(),
	})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("Success", nil, ctx)
}
