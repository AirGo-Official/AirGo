package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// 购买流程：获取订单详情（前端：立即购买）->订单预创建（前端：提交订单）->购买主逻辑（前端：确认购买）
// 查询全部已启用商品
func GetAllEnabledGoods(ctx *gin.Context) {
	goodsArr, _, err := service.CommonSqlFind[model.Goods, string, []model.Goods]("status = true ORDER BY goods_order")
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetAllEnabledGoods error:"+err.Error(), nil, ctx)
		return
	}
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
	response.OK("GoodsSort success", nil, ctx)
}
