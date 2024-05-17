package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
	uuid "github.com/satori/go.uuid"
)

// GetCustomerServiceList
// @Tags [customer api] customer service
// @Summary 获取服务列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/customerService/getCustomerServiceList [get]
func GetCustomerServiceList(ctx *gin.Context) {
	uID, _ := api.GetUserIDFromGinContext(ctx)
	csArr, err := service.CustomerServiceSvc.
		GetCustomerServiceList(&model.CustomerService{UserID: uID, ServiceStatus: true})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("Success", csArr, ctx)
}

// DeleteCustomerService
// @Tags [customer api] customer service
// @Summary 删除服务
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.CustomerService true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/customerService/deleteCustomerService [delete]
func DeleteCustomerService(ctx *gin.Context) {
	uID, _ := api.GetUserIDFromGinContext(ctx)
	var cs model.CustomerService
	err := ctx.ShouldBind(&cs)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CustomerServiceSvc.
		DeleteCustomerService(&model.CustomerService{ID: cs.ID, UserID: uID})
	if err != nil {
		response.Fail("DeleteCustomerService error:"+err.Error(), nil, ctx)
		global.Logrus.Error("DeleteCustomerService error:", err.Error())
		return
	}
	response.OK("Success", nil, ctx)
}

// ResetSubscribeUUID
// @Tags [customer api] customer service
// @Summary 重置订阅uuid
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.CustomerService true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/customerService/resetSubscribeUUID [post]
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
	err = service.CustomerServiceSvc.UpdateCustomerService(cs.ID, map[string]any{
		"sub_uuid": cs.SubUUID,
	})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("Success", nil, ctx)
}

// PushCustomerService
// @Tags [customer api] customer service
// @Summary push
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.PushCustomerServiceRequest true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/customerService/pushCustomerService [post]
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
	toUser, err := service.UserSvc.FirstUser(&model.User{UserName: cs.ToUserName})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	res, err := service.CustomerServiceSvc.FirstCustomerService(&model.CustomerService{ID: cs.CustomerServiceID, UserID: uID, IsRenew: true, ServiceStatus: true})
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	err = service.CustomerServiceSvc.UpdateCustomerService(res.ID, map[string]any{
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
