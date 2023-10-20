package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/other_plugin"
	"AirGo/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// 支付主逻辑
func Purchase(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}
	// 前端传的订单信息，前端下次优化：只传订单id和支付方式id
	var receiveOrder model.Orders
	err := ctx.ShouldBind(&receiveOrder)
	if err != nil || receiveOrder.OutTradeNo == "" {
		response.Fail("订单参数获取错误", nil, ctx)
		return
	}
	//根据订单号查询订单
	receiveOrder.UserID = uIDInt //确认user id
	sysOrder, _, err := service.CommonSqlFind[model.Orders, model.Orders, model.Orders](model.Orders{UserID: receiveOrder.UserID, OutTradeNo: receiveOrder.OutTradeNo})
	if err != nil {
		global.Logrus.Error("根据订单号查询订单error:", err.Error())
		if err == gorm.ErrRecordNotFound {
			response.Fail("订单不存在"+err.Error(), nil, ctx)
			return
		} else {
			response.Fail("订单查询错误"+err.Error(), nil, ctx)
			return
		}
	}
	//0元购，跳过支付
	totalAmountFloat64, _ := strconv.ParseFloat(sysOrder.TotalAmount, 10)
	if totalAmountFloat64 == 0 {
		sysOrder.TradeStatus = model.OrderCompleted                     //更新数据库订单状态,自定义结束状态completed
		sysOrder.ReceiptAmount = "0"                                    //实收金额
		sysOrder.BuyerPayAmount = "0"                                   //付款金额
		go service.UpdateOrder(&sysOrder)                               //更新数据库状态
		go service.UpdateUserSubscribe(&sysOrder)                       //更新用户订阅信息
		go service.RemainHandle(sysOrder.UserID, sysOrder.RemainAmount) //处理用户余额
		response.OK("购买成功", nil, ctx)
		return
	}
	//根据支付id查询支付参数
	var payParams = model.Pay{
		ID: receiveOrder.PayID,
	}
	pay, _, err := service.CommonSqlFind[model.Pay, model.Pay, model.Pay](payParams)
	sysOrder.PayID = pay.ID        //支付方式id
	sysOrder.PayType = pay.PayType //

	//判断支付方式
	switch sysOrder.PayType {
	case "epay":
		res, err := service.EpayPreByHTML(&sysOrder, &pay)
		if err != nil {
			response.Fail("epay error："+err.Error(), nil, ctx)
			return
		}
		sysOrder.TradeStatus = model.OrderWAIT_BUYER_PAY //初始订单状态：等待付款
		go service.UpdateOrder(&sysOrder)                //更新数据库
		var pcptf = model.PreCreatePayToFrontend{
			EpayInfo: *res,
		}
		response.OK("epay success:", pcptf, ctx) //返回用户易支付订单参数，采用易支付网页支付

	case "alipay":

		//创建alipay client
		client, err := service.InitAlipayClient(pay)
		if err != nil {
			response.Fail("alipay error："+err.Error(), nil, ctx)
			return
		}

		res, err := service.TradePreCreatePay(client, &sysOrder)
		if err != nil {
			response.Fail("alipay error："+err.Error(), nil, ctx)
			return
		}
		sysOrder.TradeStatus = model.OrderWAIT_BUYER_PAY //初始订单状态：等待付款
		go service.UpdateOrder(&sysOrder)                //更新数据库
		var pcptf = model.PreCreatePayToFrontend{
			AlipayInfo: model.AlipayPreCreatePayToFrontend{QRCode: res.QRCode},
		}
		response.OK("alipay success:", pcptf, ctx) //返回用户qrcode
		go service.PollAliPay(&sysOrder, client)   //5分钟等待付款，轮询
	case "wechatpay":

	}

}

// 支付宝异步回调，弃用，改为轮询
func AlipayNotify(ctx *gin.Context) {
	//noti, _ := global.AlipayClient.GetTradeNotification(ctx.Request)
	//if noti == nil {
	//	return
	//}
	////查询原始订单
	//var order = model.Orders{
	//	OutTradeNo: noti.OutTradeNo,
	//}
	//sysOrder, _ := service.CommonSqlFind[model.Orders, model.Orders, model.Orders](model.Orders{}, order)
	////根据回调参数更新数据库订单
	//sysOrder.TradeNo = noti.TradeNo
	//sysOrder.BuyerLogonId = noti.BuyerLogonId
	//sysOrder.TradeStatus = string(noti.TradeStatus)
	//sysOrder.TotalAmount = noti.TotalAmount
	//sysOrder.ReceiptAmount = noti.ReceiptAmount
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

// 易支付异步回调
func EpayNotify(ctx *gin.Context) {
	var epayRes model.EpayResultResponse
	err := ctx.ShouldBindQuery(&epayRes)
	if err != nil {
		global.Logrus.Error(err.Error())
		return
	}
	//查询原始订单
	var order = model.Orders{
		OutTradeNo: epayRes.OutTradeNo,
	}
	sysOrder, _, _ := service.CommonSqlFind[model.Orders, model.Orders, model.Orders](order)
	if sysOrder.TradeStatus == model.OrderTRADE_SUCCESS {
		ctx.String(200, "success")
		return
	}
	sysOrder.TradeNo = epayRes.TradeNo
	sysOrder.ReceiptAmount = epayRes.Money  //实收金额
	sysOrder.BuyerPayAmount = epayRes.Money //付款金额
	sysOrder.TradeStatus = epayRes.TradeStatus

	//更新数据库订单信息
	go service.UpdateOrder(&sysOrder)
	//更新用户订阅信息
	go service.UpdateUserSubscribe(&sysOrder)
	//返回success以表示服务器接收到了订单通知
	ctx.String(200, "success")

}

// 获取已激活支付列表
func GetEnabledPayList(ctx *gin.Context) {
	list, _, err := service.CommonSqlFind[model.Pay, string, []model.Pay]("status = true")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("获取已激活支付列表成功", list, ctx)
}

// 获取全部支付列表
func GetPayList(ctx *gin.Context) {
	list, _, err := service.CommonSqlFind[model.Pay, string, []model.Pay]("")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("获取支付列表成功", list, ctx)

}

// 新建支付
func NewPay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlCreate[model.Pay](receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("新建支付成功", nil, ctx)

}

// 删除支付
func DeletePay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	//fmt.Println("删除：", receivePay)
	err = service.CommonSqlDelete[model.Pay, model.Pay](model.Pay{}, receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("删除支付成功", nil, ctx)

}

// 修改支付
func UpdatePay(ctx *gin.Context) {
	var receivePay model.Pay
	err := ctx.ShouldBind(&receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Pay](receivePay)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(err.Error(), nil, ctx)
		return
	}
	response.OK("修改支付成功", nil, ctx)

}
