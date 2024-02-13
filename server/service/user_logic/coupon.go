package user_logic

import (
	"errors"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
	"time"
)

type Coupon struct{}

var couponService *Coupon

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

func (c *Coupon) VerifyCoupon(order *model.Order) (*model.Coupon, error) {
	//查询折扣
	coupon, err := c.GetCouponByName(order.CouponName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constant.ERROR_COUPON_NOT_EXIST)
		} else {
			return nil, errors.New(constant.ERROR_COUPON_QUERY_ERROR)
		}
	}
	//判断是否关联当前商品
	var ok = false
	for k := range coupon.Goods {
		if coupon.Goods[k].ID == order.GoodsID {
			ok = true //匹配到
			break
		}
	}
	if !ok {
		return nil, errors.New(constant.ERROR_COUPON_NOT_APPLICABLE)
	}
	//判断有效期
	if time.Now().After(coupon.ExpiredAt) {
		return nil, errors.New(constant.ERROR_COUPON_HAS_EXPIRED)
	}
	//判断使用次数
	//orderArr, _, err := service.CommonSqlFind[model.Order, string, []model.Order](fmt.Sprintf("user_id = %d AND coupon_id = %d", order.UserID, c.ID))
	var orders []model.Order
	err = global.DB.Where(&model.Order{UserID: order.UserID, CouponID: order.CouponID}).Find(&orders).Error

	if err != nil {
		return nil, errors.New(constant.ERROR_COUPON_QUERY_ERROR)
	}
	if int64(len(orders)) >= coupon.Limit {
		return nil, errors.New(constant.ERROR_COUPON_COUNT_EXHAUSTED)
	}
	//返回
	return coupon, nil
}
