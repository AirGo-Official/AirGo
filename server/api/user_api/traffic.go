package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetSubTrafficList
// @Tags [customer api] traffic
// @Summary 获取订阅流量记录
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.CustomerService true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/traffic/getSubTrafficList [post]
func GetSubTrafficList(ctx *gin.Context) {
	var params model.CustomerService
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	list, err := service.TrafficSvc.GetSubTrafficList(&params)
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("Success", list, ctx)
}
