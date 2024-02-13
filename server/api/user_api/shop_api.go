package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

// 支付主逻辑
func Purchase(ctx *gin.Context) {
	//uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	//if !ok {
	//	response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR, nil, ctx)
	//	return
	//}
	// 前端传的订单信息
	var orderRequest model.Order
	err := ctx.ShouldBind(&orderRequest)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	fmt.Println("api Purchase:")
	admin_logic.Show(orderRequest)
	//根据订单号查询订单
	value, ok := global.LocalCache.Get(constant.CACHE_SUBMIT_ORDER_BY_ORDERID + orderRequest.OutTradeNo)
	if !ok {
		response.Fail(constant.ERROR_ORDER_TIMED_OUT, nil, ctx)
		return
	}
	sysOrder := value.(*model.Order)
	// 获取支付信息
	orderWaitPay, err := shopService.Purchase(sysOrder)

	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Purchase error:"+err.Error(), nil, ctx)
		return
	}
	//删除
	response.OK("Purchase success", orderWaitPay, ctx)
}

// 查询已启用商品列表
func GetEnabledGoodsList(ctx *gin.Context) {
	//获取查询参数
	goods_type, ok := ctx.GetQuery("goods_type")
	if !ok || goods_type == "" {
		global.Logrus.Error("GetQuery error")
		response.Fail("GetEnabledGoodsList error:GetQuery error", nil, ctx)
		return
	}
	goodsArr, err := shopService.GetGoodsList(&model.Goods{GoodsType: goods_type, IsShow: true, IsSale: true})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetEnabledGoodsList error:"+err.Error(), nil, ctx)
		return
	}
	//global.LocalCache.SetNoExpire(cacheKey, goodsArr)
	response.OK("GetEnabledGoodsList success", goodsArr, ctx)
}
