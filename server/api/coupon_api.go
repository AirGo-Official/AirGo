package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// 新建折扣
func NewCoupon(ctx *gin.Context) {
	var coupon model.Coupon
	err := ctx.ShouldBind(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewCoupon error:"+err.Error(), nil, ctx)
		return
	}
	err = service.NewCoupon(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewCoupon error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewCoupon success", nil, ctx)
}

// 删除折扣
func DeleteCoupon(ctx *gin.Context) {
	var coupon model.Coupon
	err := ctx.ShouldBind(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteCoupon error:"+err.Error(), nil, ctx)
		return
	}
	err = service.DeleteCoupon(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteCoupon error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteCoupon success", nil, ctx)
}

// 更新折扣
func UpdateCoupon(ctx *gin.Context) {
	var coupon model.Coupon
	err := ctx.ShouldBind(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateCoupon error:"+err.Error(), nil, ctx)
		return
	}
	err = service.UpdateCoupon(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateCoupon error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateCoupon success", nil, ctx)
}

// 获取折扣列表
func GetCoupon(ctx *gin.Context) {
	res, err := service.GetAllCoupon()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetCoupon error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetCoupon success", res, ctx)
}
