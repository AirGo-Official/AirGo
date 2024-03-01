package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/response"
	"time"
)

// 获取用户订单by user id
func GetOrderList(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("GetOrderList error:user id error", nil, ctx)
		return
	}
	res, err := orderService.GetUserOrders(&params, uIDInt)
	if err != nil {
		response.Fail("GetOrderList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetOrderList success", res, ctx)
}

// 获取订单详情（计算价格等）
func GetOrderInfo(ctx *gin.Context) {
	orderReq, err := OrderRequestHandler(ctx)
	if err != nil {
		response.Fail("GetOrderInfo error:"+err.Error(), nil, ctx)
		return
	}
	preOrder, msg, err := orderService.PreHandleOrder(orderReq)

	if err != nil {
		response.Fail("GetOrderInfo error:"+err.Error(), nil, ctx)
		return
	}
	if msg != "" {
		response.Response(constant.RESPONSE_WARNING, msg, nil, ctx) //目前msg是用来提示折扣码处理信息
		return
	}
	response.OK("GetOrderInfo success", preOrder, ctx)
}

// 订单预创建，生成系统订单（提交订单）
func PreCreateOrder(ctx *gin.Context) {
	// 1、订单请求预处理
	orderReq, err := OrderRequestHandler(ctx)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("PreCreateOrder error:"+err.Error(), nil, ctx)
		return
	}
	// 2、订单拦截逻辑，处理一些校验
	err = orderService.PreCheckOrder(orderReq)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("PreCreateOrder error:"+err.Error(), nil, ctx)
		return
	}
	// 3、订单预处理，计算价格
	preOrder, _, err := orderService.PreHandleOrder(orderReq)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("PreCreateOrder error:"+err.Error(), nil, ctx)
		return
	}
	// 4、幂等性，用户订单在队列处理完成时，删除cache
	switch orderReq.OrderType {
	case constant.ORDER_TYPE_NEW:
		global.LocalCache.
			Set(fmt.Sprintf("%s%d:%d", constant.CACHE_USERID_AND_GOODSID, orderReq.UserID, orderReq.GoodsID),
				orderReq.OutTradeNo,
				constant.CACHE_SUBMIT_ORDER_TIMEOUT*time.Minute)
	case constant.ORDER_TYPE_RENEW:
		global.LocalCache.
			Set(fmt.Sprintf("%s%d:%d",
				constant.CACHE_USERID_AND_CUSTOMERSERVICEID, orderReq.UserID, orderReq.CustomerServiceID),
				orderReq.OutTradeNo,
				constant.CACHE_SUBMIT_ORDER_TIMEOUT*time.Minute)
	}
	// 5、异步下单，返回订单号，前端轮训订单号获取订单详情
	err = global.Queue.Publish(constant.SUBMIT_ORDER, preOrder)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("PreCreateOrder error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("PreCreateOrder success", preOrder, ctx)
}
func OrderRequestHandler(ctx *gin.Context) (*model.Order, error) {
	//
	var orderReq model.Order
	err := ctx.ShouldBind(&orderReq)
	if err != nil {
		return nil, err
	}
	uIDInt, _ := api.GetUserIDFromGinContext(ctx)
	uName, _ := api.GetUserNameFromGinContext(ctx)
	//user, _ := userService.FindUserByID(uIDInt)
	orderReq.UserID = uIDInt
	orderReq.UserName = uName
	return &orderReq, nil
}
func GetOrderInfoWaitPay(ctx *gin.Context) {
	orderRequest, err := OrderRequestHandler(ctx)

	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetOrderInfoWaitPay error:"+err.Error(), nil, ctx)
		return
	}
	var order *model.Order
	value, ok := global.LocalCache.Get(constant.CACHE_SUBMIT_ORDER_BY_ORDERID + orderRequest.OutTradeNo)
	if !ok {
		response.Fail(constant.ERROR_ORDER_TIMED_OUT, nil, ctx)
		return
	}
	order = value.(*model.Order)
	response.OK("GetOrderInfoWaitPay success", order, ctx)
}
