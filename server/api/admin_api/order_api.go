package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetOrderList
// @Tags [admin api] order
// @Summary 获取全部订单，分页获取
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/order/getOrderList [post]
func GetOrderList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	res, total, err := service.CommonSqlFindWithFieldParams(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetOrderList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetOrderList success", model.CommonDataResp{
		Total: total,
		Data:  res,
	}, ctx)

}

// OrderSummary
// @Tags [admin api] order
// @Summary 获取订单统计
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/order/orderSummary [post]
func OrderSummary(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	res, err := service.AdminOrderSvc.OrderSummary(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	response.OK("OrderSummary success", res, ctx)
}

// UpdateOrder
// @Tags [admin api] order
// @Summary 更新用户订单
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Order true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/order/updateOrder [post]
func UpdateOrder(ctx *gin.Context) {
	var order model.Order
	err := ctx.ShouldBind(&order)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminOrderSvc.UpdateOrder(&order) //更新数据库状态
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateOrder error:"+err.Error(), nil, ctx)
		return
	}
	//如果订单状态是 支付成功 或 交易结束，将订单进行处理
	if order.TradeStatus == constant.ORDER_STATUS_TRADE_SUCCESS || order.TradeStatus == constant.ORDER_STATUS_TRADE_FINISHED {
		service.OrderSvc.DeleteOneOrderFromCache(&order) //删除缓存
		err = service.OrderSvc.PaymentSuccessfullyOrderHandler(&order)
		if err != nil {
			if err != nil {
				global.Logrus.Error(err)
				response.Fail("UpdateOrder error:"+err.Error(), nil, ctx)
				return
			}
		}
	} else {
		service.OrderSvc.UpdateOneOrderToCache(&order) //更新缓存
	}
	response.OK("UpdateOrder success", nil, ctx)
}
