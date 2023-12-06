package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// 查询全部已启用商品
func GetAllEnabledGoods(ctx *gin.Context) {
	goodsArr, ok := global.LocalCache.Get(global.AllEnabledGoods)
	if ok {
		response.OK("GetRouteList success", goodsArr, ctx)
		return
	}
	goodsArr, _, err := service.CommonSqlFind[model.Goods, string, []model.Goods]("status = true ORDER BY goods_order")
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
