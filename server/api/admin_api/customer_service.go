package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetCustomerServiceList
// @Tags [admin api] customer service
// @Summary 获取用户服务列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.CustomerService true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/customerService/getCustomerServiceList [post]
func GetCustomerServiceList(ctx *gin.Context) {
	var cs model.CustomerService
	err := ctx.ShouldBind(&cs)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	list, err := service.AdminCustomerServiceSvc.GetCustomerServiceList(&cs)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetCustomerServiceList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetCustomerServiceList success", list, ctx)

}

// UpdateCustomerService
// @Tags [admin api] customer service
// @Summary 更新客户服务
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.CustomerService true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/customerService/updateCustomerService [post]
func UpdateCustomerService(ctx *gin.Context) {
	var cs model.CustomerService
	err := ctx.ShouldBind(&cs)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminCustomerServiceSvc.UpdateCustomerService(&cs)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateCustomerService error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateCustomerService success", nil, ctx)
}

// DeleteCustomerService
// @Tags [admin api] customer service
// @Summary 删除客户服务
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.CustomerService true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/customerService/deleteCustomerService [delete]
func DeleteCustomerService(ctx *gin.Context) {
	var cs model.CustomerService
	err := ctx.ShouldBind(&cs)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminCustomerServiceSvc.DeleteCustomerService(&cs)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteCustomerService error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteCustomerService success", nil, ctx)
}
