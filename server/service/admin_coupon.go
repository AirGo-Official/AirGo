package service

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type AdminCoupon struct{}

var AdminCouponSvc *AdminCoupon

// 查询全部折扣
func (c *AdminCoupon) GetCouponList() (*model.CommonDataResp, error) {
	var couponArr []model.Coupon
	var total int64
	err := global.DB.Model(&model.Coupon{}).Count(&total).Preload("Goods").Find(&couponArr).Error
	return &model.CommonDataResp{total, couponArr}, err
}

// 新建折扣
func (c *AdminCoupon) NewCoupon(coupon *model.Coupon) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&coupon).Error
	})
}

// 更新折扣
func (c *AdminCoupon) UpdateCoupon(couponParams *model.Coupon) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&couponParams).Association("Goods").Replace(&couponParams.Goods)
		if err != nil {
			return err
		}
		err = tx.Save(&couponParams).Error
		if err != nil {
			return err
		}
		// 删除缓存
		c.DeleteCouponCache(couponParams)
		return nil
	})
}

// 删除折扣
func (c *AdminCoupon) DeleteCoupon(couponParams *model.Coupon) error {
	// 开启事务
	tx := global.DB.Begin()
	//删除关联
	err := tx.Model(&couponParams).Association("Goods").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	// 删除缓存
	c.DeleteCouponCache(couponParams)
	//删除
	err = tx.Delete(&couponParams).Error
	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return err
	}
}

func (c *AdminCoupon) DeleteCouponCache(couponParams *model.Coupon) {
	global.LocalCache.Delete(constant.CACHE_COUPON_BY_NAME + couponParams.Name)
}
