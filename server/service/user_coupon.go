package service

import (
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Coupon struct{}

var CouponSvc *Coupon

// 查询优惠券，预加载商品
func (c *Coupon) GetCouponByName(couponName string) (*model.Coupon, error) {
	var coupon model.Coupon
	cache, ok := global.LocalCache.Get(constant.CACHE_COUPON_BY_NAME + couponName)
	if ok {
		coupon = cache.(model.Coupon)
		return &coupon, nil
	}
	err := global.DB.Where(&model.Coupon{Name: couponName}).Preload("Goods").First(&coupon).Error
	if err != nil {
		return nil, err
	}
	global.LocalCache.SetNoExpire(constant.CACHE_COUPON_BY_NAME+couponName, coupon)
	return &coupon, err
}

func (c *Coupon) VerifyCoupon(preOrder *model.Order) string {
	total, _ := strconv.ParseFloat(preOrder.TotalAmount, 64)
	//查询折扣
	coupon, err := c.GetCouponByName(preOrder.CouponName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constant.ERROR_COUPON_NOT_EXIST
		} else {
			return constant.ERROR_COUPON_QUERY_ERROR
		}
	}
	//判断是否关联当前商品
	var ok = false
	for k := range coupon.Goods {
		if coupon.Goods[k].ID == preOrder.GoodsID {
			ok = true //匹配到
			break
		}
	}
	if !ok {
		return constant.ERROR_COUPON_NOT_APPLICABLE
	}
	//判断有效期
	if time.Now().After(coupon.ExpiredAt) {
		return constant.ERROR_COUPON_HAS_EXPIRED
	}
	//判断最低使用金额
	if total < coupon.MinAmount {
		return constant.ERROR_COUPON_NOT_APPLICABLE
	}
	//判断使用次数
	var orders []model.Order
	err = global.DB.
		Where(&model.Order{UserID: preOrder.UserID, CouponID: coupon.ID}).
		Where("trade_status <> ?", constant.ORDER_STATUS_TRADE_CLOSED). //订单关闭，不计算折扣码使用次数
		Find(&orders).Error

	if err != nil {
		return constant.ERROR_COUPON_QUERY_ERROR
	}
	if int64(len(orders)) >= coupon.Limit {
		return constant.ERROR_COUPON_COUNT_EXHAUSTED
	}
	//处理折扣价格
	if coupon.DiscountRate != 0 {
		preOrder.CouponAmount = fmt.Sprintf("%.2f", total*coupon.DiscountRate)
		preOrder.CouponID = coupon.ID
		preOrder.TotalAmount = fmt.Sprintf("%.2f", total-total*coupon.DiscountRate)
	}
	//返回
	return ""
}
