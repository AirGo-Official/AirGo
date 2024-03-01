package model

import "time"

type Coupon struct {
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-"             gorm:"index"`
	ID           int64      `json:"id"            gorm:"primaryKey"`
	Name         string     `json:"name"          gorm:"comment:名称"`
	DiscountRate float64    `json:"discount_rate" gorm:"default:0.9;comment:折扣率,实际价格=原价*折扣率"`
	Limit        int64      `json:"limit"         gorm:"default:1;comment:限制次数"`
	ExpiredAt    time.Time  `json:"expired_at"    gorm:"comment:过期时间"`
	MinAmount    float64    `json:"min_amount"    gorm:"comment:最低使用金额"`
	Goods        []Goods    `json:"goods"         gorm:"many2many:goods_and_coupon"`
}
