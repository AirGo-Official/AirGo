package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/response"
)

func GetCustomerServiceList(ctx *gin.Context) {
	var cs model.CustomerService
	err := ctx.ShouldBind(&cs)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	list, err := customerService.GetCustomerServiceList(&cs)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetCustomerServiceList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetCustomerServiceList success", list, ctx)

}

// 更新客户服务
func UpdateCustomerService(ctx *gin.Context) {
	var cs model.CustomerService
	err := ctx.ShouldBind(&cs)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = customerService.UpdateCustomerService(&cs)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateCustomerService error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateCustomerService success", nil, ctx)
}
