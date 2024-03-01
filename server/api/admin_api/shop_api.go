package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

// 查询全部商品
func GetGoodsList(ctx *gin.Context) {
	goodsArr, err := shopService.GetGoodsList()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetGoodsList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetGoodsList success", goodsArr, ctx)

}

// 新建商品
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
	err = shopService.NewGoods(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewGoods error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(cacheKey)
	response.OK("NewGoods success", nil, ctx)
}

// 删除商品
func DeleteGoods(ctx *gin.Context) {
	var goods model.Goods
	err := ctx.ShouldBind(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = shopService.DeleteGoods(&goods)
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

// 更新商品
func UpdateGoods(ctx *gin.Context) {
	var goods model.Goods
	err := ctx.ShouldBind(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = shopService.UpdateGoods(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateGoods error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateGoods success", nil, ctx)
}

// 排序
func GoodsSort(ctx *gin.Context) {
	var arr []model.Goods
	err := ctx.ShouldBind(&arr)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlUpdateMultiLine[[]model.Goods](arr, "id", []string{"goods_order"})
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
