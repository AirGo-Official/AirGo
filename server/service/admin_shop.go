package service

import (
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type AdminShop struct{}

var AdminShopSvc *AdminShop

// 查询全部商品
func (s *AdminShop) GetGoodsList() (*[]model.Goods, error) {
	var goodsArr []model.Goods
	err := global.DB.Model(&model.Goods{}).Preload("Nodes").Order("goods_order").Find(&goodsArr).Error
	return &goodsArr, err
}

// 查询商品 by nodeID
func (s *AdminShop) FindGoodsByNodeID(nodeID int64) ([]model.Goods, error) {
	var node model.Node
	err := global.DB.Where("id = ?", nodeID).Preload("Goods").First(&node).Error
	if err != nil {
		return nil, err
	}
	return node.Goods, nil
}

// 新建商品
func (s *AdminShop) NewGoods(goods *model.Goods) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&goods).Error
	})
}

// 删除商品
func (s *AdminShop) DeleteGoods(goods *model.Goods) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&goods).Association("Nodes").Replace(nil)
		if err != nil {
			return err
		}
		return tx.Delete(&goods).Error
	})
}

// 更新商品
func (s *AdminShop) UpdateGoods(goodsParams *model.Goods) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//更新关联节点
		err := tx.Model(&goodsParams).Association("Nodes").Replace(&goodsParams.Nodes)
		if err != nil {
			return err
		}
		// 更新商品
		return tx.Save(&goodsParams).Error
	})
	if err != nil {
		return err
	}
	// 删除缓存
	key1 := fmt.Sprintf("%s%d", constant.CACHE_GOODS_BY_ID, goodsParams.ID)
	global.LocalCache.Delete(key1)
	key2 := ""
	switch goodsParams.GoodsType {
	case constant.GOODS_TYPE_GENERAL:
		key2 = constant.CACHE_ALL_ENABLED_GOODS_GENERAL
	case constant.GOODS_TYPE_SUBSCRIBE:
		key2 = constant.CACHE_ALL_ENABLED_GOODS_SUBSCRIBE
	case constant.GOODS_TYPE_RECHARGE:
		key2 = constant.CACHE_ALL_ENABLED_GOODS_RECHARGE
	default:
	}
	global.LocalCache.Delete(key2)
	return nil
}
