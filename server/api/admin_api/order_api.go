package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/service/user_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

// 获取全部订单，分页获取
func GetOrderList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	res, total, err := common_logic.CommonSqlFindWithFieldParams(&params)
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

// 获取订单统计
func OrderSummary(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	res, err := orderService.OrderSummary(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	response.OK("OrderSummary success", res, ctx)
}

// 更新用户订单
func UpdateOrder(ctx *gin.Context) {
	var order model.Order
	err := ctx.ShouldBind(&order)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = orderService.UpdateOrder(&order) //更新数据库状态
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateOrder error:"+err.Error(), nil, ctx)
		return
	}
	//如果订单状态是 支付成功 或 交易结束，将订单进行处理
	var userOrderService user_logic.Order
	if order.TradeStatus == constant.ORDER_STATUS_TRADE_SUCCESS || order.TradeStatus == constant.ORDER_STATUS_TRADE_FINISHED {
		userOrderService.DeleteOneOrderFromCache(&order) //删除缓存
		err = userOrderService.PaymentSuccessfullyOrderHandler(&order)
		if err != nil {
			if err != nil {
				global.Logrus.Error(err)
				response.Fail("UpdateOrder error:"+err.Error(), nil, ctx)
				return
			}
		}
	} else {
		userOrderService.UpdateOneOrderToCache(&order) //更新缓存
	}
	response.OK("UpdateOrder success", nil, ctx)
}
