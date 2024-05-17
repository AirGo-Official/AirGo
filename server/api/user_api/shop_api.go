package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// Purchase
// @Tags [customer api] shop
// @Summary 支付主逻辑
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Order true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/shop/purchase [post]
func Purchase(ctx *gin.Context) {
	// 前端传的订单信息
	var orderRequest model.Order
	err := ctx.ShouldBind(&orderRequest)
	if err != nil {
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//fmt.Println("api Purchase:")
	//admin_logic.Show(orderRequest)
	//根据订单号查询订单
	value, ok := global.LocalCache.Get(constant.CACHE_SUBMIT_ORDER_BY_ORDERID + orderRequest.OutTradeNo)
	if !ok {
		response.Fail(constant.ERROR_ORDER_TIMED_OUT, nil, ctx)
		return
	}
	sysOrder := value.(*model.Order)
	sysOrder.PayID = orderRequest.PayID //提取前端传的pay_id
	// 获取支付信息
	orderResult, err := service.ShopSvc.Purchase(sysOrder)

	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Purchase error:"+err.Error(), nil, ctx)
		return
	}
	//1、如果是余额支付，到这一步说明处理完毕。余额不足或其他错误会在上一步报错
	//2、如果是alipay和epay，则返回支付信息，由前端展示。alipay和epay的错误会在上一步报错
	response.OK("success", orderResult, ctx)
}

// GetEnabledGoodsList
// @Tags [customer api] shop
// @Summary 查询已启用商品列表
// @Produce json
// @Param Authorization header string true "Bearer 用户token"
// @Param goods_type query string true "商品类型"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/shop/getEnabledGoodsList [get]
func GetEnabledGoodsList(ctx *gin.Context) {
	//获取查询参数
	goods_type, ok := ctx.GetQuery("goods_type")
	if !ok || goods_type == "" {
		global.Logrus.Error("GetQuery error")
		response.Fail("GetEnabledGoodsList error:GetQuery error", nil, ctx)
		return
	}
	goodsArr, err := service.ShopSvc.GetGoodsList(&model.Goods{GoodsType: goods_type, IsShow: true, IsSale: true})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetEnabledGoodsList error:"+err.Error(), nil, ctx)
		return
	}
	//global.LocalCache.SetNoExpire(cacheKey, goodsArr)
	response.OK("GetEnabledGoodsList success", goodsArr, ctx)
}
