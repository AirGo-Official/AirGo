package service

import (
	"AirGo/global"
	"AirGo/model"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func VerifyCoupon(order *model.Orders) (model.Coupon, error) {
	//查询折扣
	var c model.Coupon
	err := global.DB.Where(&model.Coupon{Name: order.CouponName}).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Coupon{}, errors.New("优惠码不存在")
		} else {
			return model.Coupon{}, errors.New("优惠码错误")
		}
	}
	//判断是否关联当前商品
	var gac model.GoodsAndCoupon
	err = global.DB.Where("goods_id = ? AND coupon_id = ?", order.GoodsID, c.ID).First(&gac).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Coupon{}, errors.New("优惠码不适用")
		} else {
			return model.Coupon{}, errors.New("优惠码错误")
		}
	}
	//判断有效期
	if time.Now().After(c.ExpiredAt) {
		return model.Coupon{}, errors.New("优惠码已过期")
	}
	//判断使用次数
	orderArr, _, err := CommonSqlFind[model.Orders, string, []model.Orders]("user_id = " + strconv.FormatInt(order.UserID, 10) + " AND coupon_id = " + strconv.FormatInt(c.ID, 10))
	if err != nil {
		return model.Coupon{}, errors.New("优惠码错误")
	}
	if int64(len(orderArr)) >= c.Limit {
		return model.Coupon{}, errors.New("优惠码次数用尽")
	}
	//返回
	return c, nil
}

// 查询全部折扣
func GetAllCoupon() (*[]model.Coupon, error) {
	var couponArr []model.Coupon
	err := global.DB.Model(&model.Coupon{}).Preload("Goods").Find(&couponArr).Error
	if len(couponArr) == 0 {
		return &couponArr, err
	} else { //处理商品关联的节点
		for k1, _ := range couponArr {
			for _, v2 := range couponArr[k1].Goods {
				couponArr[k1].CheckedGoods = append(couponArr[k1].CheckedGoods, v2.ID)
			}
		}
		return &couponArr, err
	}

}

// 新建折扣
func NewCoupon(coupon *model.Coupon) error {
	//查询关联商品
	global.DB.Where("id in ?", coupon.CheckedGoods).Find(&coupon.Goods)
	err := global.DB.Create(&coupon).Error
	return err
}

// 更新折扣
func UpdateCoupon(coupon *model.Coupon) error {
	//查询关联商品
	var goodsArr []model.Goods
	global.DB.Where("id in ?", coupon.CheckedGoods).Find(&goodsArr)
	coupon.Goods = goodsArr
	//更新关联商品
	global.DB.Model(&coupon).Association("Goods").Replace(&coupon.Goods)
	// 更新折扣
	err := global.DB.Save(&coupon).Error
	return err

}

// 删除折扣
func DeleteCoupon(coupon *model.Coupon) error {
	//删除关联
	err := global.DB.Model(&model.Coupon{ID: coupon.ID}).Association("Goods").Replace(nil)
	if err != nil {
		return err
	}
	err = global.DB.Delete(&coupon).Error
	return err
}
