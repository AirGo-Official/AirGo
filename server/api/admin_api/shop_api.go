package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	service "github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetGoodsList
// @Tags [admin api] shop
// @Summary 查询全部商品
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/shop/getGoodsList [get]
func GetGoodsList(ctx *gin.Context) {
	goodsArr, err := service.AdminShopSvc.GetGoodsList()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetGoodsList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetGoodsList success", goodsArr, ctx)

}

// NewGoods
// @Tags [admin api] shop
// @Summary 新建商品
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Goods true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/shop/newGoods [post]
func NewGoods(ctx *gin.Context) {
	var goods model.Goods
	err := ctx.ShouldBind(&goods)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//根据商品类型，修改一些默认参数
	var cacheKey string
	switch goods.GoodsType {
	case constant.GOODS_TYPE_SUBSCRIBE:
		cacheKey = constant.CACHE_ALL_ENABLED_GOODS_SUBSCRIBE
		goods.DeliverType = constant.DELIVER_TYPE_NONE

	case constant.GOODS_TYPE_GENERAL:
		cacheKey = constant.CACHE_ALL_ENABLED_GOODS_GENERAL

	case constant.GOODS_TYPE_RECHARGE:
		cacheKey = constant.CACHE_ALL_ENABLED_GOODS_RECHARGE
		goods.DeliverType = constant.DELIVER_TYPE_NONE
	}
	err = service.AdminShopSvc.NewGoods(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewGoods error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(cacheKey)
	response.OK("NewGoods success", nil, ctx)
}

// DeleteGoods
// @Tags [admin api] shop
// @Summary 删除商品
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Goods true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/shop/deleteGoods [delete]
func DeleteGoods(ctx *gin.Context) {
	var goods model.Goods
	err := ctx.ShouldBind(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminShopSvc.DeleteGoods(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteGoods error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(constant.CACHE_ALL_ENABLED_GOODS_GENERAL)
	global.LocalCache.Delete(constant.CACHE_ALL_ENABLED_GOODS_SUBSCRIBE)
	global.LocalCache.Delete(constant.CACHE_ALL_ENABLED_GOODS_RECHARGE)
	response.OK("DeleteGoods success", nil, ctx)

}

// UpdateGoods
// @Tags [admin api] shop
// @Summary 更新商品
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Goods true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/shop/updateGoods [post]
func UpdateGoods(ctx *gin.Context) {
	var goods model.Goods
	err := ctx.ShouldBind(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminShopSvc.UpdateGoods(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateGoods error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateGoods success", nil, ctx)
}

// GoodsSort
// @Tags [admin api] shop
// @Summary 商品排序
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body []model.Goods true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/shop/goodsSort [post]
func GoodsSort(ctx *gin.Context) {
	var arr []model.Goods
	err := ctx.ShouldBind(&arr)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlUpdateMultiLine[[]model.Goods](arr, "id", []string{"goods_order"})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GoodsSort error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(constant.CACHE_ALL_ENABLED_GOODS_GENERAL)
	global.LocalCache.Delete(constant.CACHE_ALL_ENABLED_GOODS_SUBSCRIBE)
	global.LocalCache.Delete(constant.CACHE_ALL_ENABLED_GOODS_RECHARGE)
	response.OK("GoodsSort success", nil, ctx)
}
