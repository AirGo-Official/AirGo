package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// NewCoupon
// @Tags [admin api] coupon
// @Summary 新建折扣
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Coupon true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/coupon/newCoupon [post]
func NewCoupon(ctx *gin.Context) {
	var coupon model.Coupon
	err := ctx.ShouldBind(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminCouponSvc.NewCoupon(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewCoupon error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewCoupon success", nil, ctx)
}

// DeleteCoupon
// @Tags [admin api] coupon
// @Summary 删除折扣
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Coupon true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/coupon/deleteCoupon [delete]
func DeleteCoupon(ctx *gin.Context) {
	var coupon model.Coupon
	err := ctx.ShouldBind(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminCouponSvc.DeleteCoupon(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteCoupon error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteCoupon success", nil, ctx)
}

// UpdateCoupon
// @Tags [admin api] coupon
// @Summary 更新折扣
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Coupon true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/coupon/updateCoupon [post]
func UpdateCoupon(ctx *gin.Context) {
	var coupon model.Coupon
	err := ctx.ShouldBind(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminCouponSvc.UpdateCoupon(&coupon)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateCoupon error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateCoupon success", nil, ctx)
}

// GetCouponList
// @Tags [admin api] coupon
// @Summary 获取折扣列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/coupon/getCouponList [post]
func GetCouponList(ctx *gin.Context) {
	res, err := service.AdminCouponSvc.GetCouponList()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetCouponList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetCouponList success", res, ctx)
}
