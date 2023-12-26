package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// 查询全部已启用商品
func GetAllEnabledGoods(ctx *gin.Context) {
	//获取查询参数
	goods_type, ok := ctx.GetQuery("goods_type")
	if !ok || goods_type == "" {
		global.Logrus.Error("GetQuery error")
		response.Fail("GetAllEnabledGoods error:GetQuery error", nil, ctx)
		return
	}
	var localGoodsArr any
	switch goods_type {
	case model.GoodsTypeGeneral:
		localGoodsArr, ok = global.LocalCache.Get(model.AllEnabledGoodsGeneral)
	case model.GoodsTypeSubscribe:
		localGoodsArr, ok = global.LocalCache.Get(model.AllEnabledGoodsSubscribe)
	case model.GoodsTypeRecharge:
		localGoodsArr, ok = global.LocalCache.Get(model.AllEnabledGoodsRecharge)
	}
	if ok {
		goodsArr := localGoodsArr.([]model.Goods)
		if len(goodsArr) > 0 {
			response.OK("GetRouteList success", goodsArr, ctx)
			return
		}
	}
	goodsArr, _, err := service.CommonSqlFind[model.Goods, string, []model.Goods](fmt.Sprintf("goods_type = '%s' AND status = true ORDER BY goods_order", goods_type))
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetAllEnabledGoods error:"+err.Error(), nil, ctx)
		return
	}
	global.GoroutinePool.Submit(func() {
		global.LocalCache.SetNoExpire(global.AllEnabledGoods, goodsArr)
	})
	response.OK("GetAllEnabledGoods success", goodsArr, ctx)
}

// 查询全部商品
func GetAllGoods(ctx *gin.Context) {
	goodsArr, err := service.GetAllGoods()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetAllGoods error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetAllGoods success", goodsArr, ctx)

}

// 新建商品
func NewGoods(ctx *gin.Context) {
	var goods model.Goods
	err := ctx.ShouldBind(&goods)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("NewGoods error:"+err.Error(), nil, ctx)
		return
	}
	//根据商品类型，修改一些默认参数
	switch goods.GoodsType {
	case model.GoodsTypeSubscribe:
		goods.DeliverType = model.DeliverTypeNone

	case model.GoodsTypeGeneral:

	case model.GoodsTypeRecharge:
		goods.DeliverType = model.DeliverTypeNone
	}
	err = service.NewGoods(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewGoods error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(global.AllEnabledGoods)
	response.OK("NewGoods success", nil, ctx)
}

// 删除商品
func DeleteGoods(ctx *gin.Context) {
	var goods model.Goods
	err := ctx.ShouldBind(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteGoods error:"+err.Error(), nil, ctx)
		return
	}
	err = service.DeleteGoods(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteGoods error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(global.AllEnabledGoods)
	response.OK("DeleteGoods success", nil, ctx)

}

// 更新商品
func UpdateGoods(ctx *gin.Context) {
	var goods model.Goods
	err := ctx.ShouldBind(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateGoods error:"+err.Error(), nil, ctx)
		return
	}
	err = service.UpdateGoods(&goods)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateGoods error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(global.AllEnabledGoods)
	response.OK("UpdateGoods success", nil, ctx)
}

// 排序
func GoodsSort(ctx *gin.Context) {
	var arr []model.Goods
	err := ctx.ShouldBind(&arr)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GoodsSort error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlUpdateMultiLine[[]model.Goods](arr, "id", []string{"goods_order"})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GoodsSort error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(global.AllEnabledGoods)
	response.OK("GoodsSort success", nil, ctx)
}
