package public_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/service/user_logic"
)

var orderService user_logic.Order

// 易支付异步回调
func EpayNotify(ctx *gin.Context) {
	var epayRes model.EpayResultResponse
	err := ctx.ShouldBindQuery(&epayRes)
	if err != nil {
		global.Logrus.Error(err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	//查询原始订单
	var order = model.Order{
		OutTradeNo: epayRes.OutTradeNo,
	}
	sysOrder, _, _ := common_logic.CommonSqlFind[model.Order, model.Order, model.Order](order)
	if sysOrder.TradeStatus == model.OrderTRADE_SUCCESS {
		ctx.String(200, "success")
		return
	}
	sysOrder.TradeNo = epayRes.TradeNo
	sysOrder.BuyerPayAmount = epayRes.Money //付款金额
	sysOrder.TradeStatus = epayRes.TradeStatus
	orderService.PaymentSuccessfullyOrderHandler(&sysOrder)
	//返回success以表示服务器接收到了订单通知
	ctx.String(200, "success")

}

// 支付宝异步回调，弃用，改为轮询
func AlipayNotify(ctx *gin.Context) {
	//noti, _ := global.AlipayClient.GetTradeNotification(ctx.Request)
	//if noti == nil {
	//	return
	//}
	////查询原始订单
	//var order = model.Order{
	//	OutTradeNo: noti.OutTradeNo,
	//}
	//sysOrder, _ := service.CommonSqlFind[model.Order, model.Order, model.Order](model.Order{}, order)
	////根据回调参数更新数据库订单
	//sysOrder.TradeNo = noti.TradeNo
	//sysOrder.BuyerLogonId = noti.BuyerLogonId
	//sysOrder.TradeStatus = string(noti.TradeStatus)
	//sysOrder.Price = noti.Price
	//sysOrder.BuyerPayAmount = noti.BuyerPayAmount
	//
	//err := service.UpdateOrder(&sysOrder)
	//if err != nil && noti.TradeStatus == model.OrderTRADE_SUCCESS {
	//	global.Logrus.Error("更新数据库订单错误", err.Error())
	//	return
	//}
	//// 确认收到通知消息
	//alipay.ACKNotification(ctx.Writer)
	////更新用户订阅信息
	//service.UpdateUserSubscribe(&sysOrder)
}
