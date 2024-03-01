package public_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/user_logic"
	"github.com/smartwalle/alipay/v3"
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
	//只有TRADE_SUCCESS是成功
	if epayRes.TradeStatus != constant.ORDER_STATUS_TRADE_SUCCESS {
		ctx.AbortWithStatus(400)
		return
	}
	//查询原始订单
	sysOrder, err := orderService.FirstUserOrder(&model.Order{OutTradeNo: epayRes.OutTradeNo})
	if err != nil {
		global.Logrus.Error(err.Error())
		ctx.AbortWithStatus(400)
		return
	}
	//如果已经交易成功，则不再进行更新
	if sysOrder.TradeStatus != constant.ORDER_STATUS_WAIT_BUYER_PAY {
		ctx.String(200, "success")
		return
	}
	sysOrder.TradeNo = epayRes.TradeNo
	sysOrder.BuyerPayAmount = epayRes.Money //付款金额
	sysOrder.TradeStatus = epayRes.TradeStatus
	_ = orderService.PaymentSuccessfullyOrderHandler(sysOrder)
	//返回success以表示服务器接收到了订单通知

	ctx.String(200, "success")

}

// AlipayNotify
// 支付宝异步回调
func AlipayNotify(ctx *gin.Context) {
	err := ctx.Request.ParseForm()
	if err != nil {
		return
	}
	out_trade_no := ctx.Request.Form.Get("out_trade_no")
	//fmt.Println("out_trade_no:", out_trade_no)

	//通过订单号查询alipay参数
	sysOrder, err := orderService.FirstUserOrder(&model.Order{OutTradeNo: out_trade_no})
	if err != nil {
		return
	}
	if sysOrder.TradeStatus != constant.ORDER_STATUS_WAIT_BUYER_PAY {
		alipay.ACKNotification(ctx.Writer)
		return
	}
	pay, err := payService.FirstPayment(&model.Pay{ID: sysOrder.PayID})
	if err != nil {
		return
	}
	//生成alipay client
	client, err := payService.InitAlipayClient(pay)
	if err != nil {
		return
	}
	//获取回调信息
	notification, err := client.GetTradeNotification(ctx.Request)
	if notification == nil || err != nil {
		return
	}
	if notification.TradeStatus != constant.ORDER_STATUS_TRADE_SUCCESS && notification.TradeStatus != constant.ORDER_STATUS_TRADE_FINISHED {
		ctx.AbortWithStatus(400)
		return
	}
	//根据回调参数更新数据库订单
	sysOrder.TradeNo = notification.TradeNo
	sysOrder.BuyerLogonId = notification.BuyerLogonId
	sysOrder.TradeStatus = string(notification.TradeStatus)
	sysOrder.BuyerPayAmount = notification.BuyerPayAmount
	_ = orderService.PaymentSuccessfullyOrderHandler(sysOrder)

	alipay.ACKNotification(ctx.Writer)
}
