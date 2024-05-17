package service

import (
	"errors"
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type Shop struct {
}

var ShopSvc *Shop

// 查询商品
func (s *Shop) FirstGoods(goodsParams *model.Goods) (*model.Goods, error) {
	var goods model.Goods
	key := fmt.Sprintf("%s%d", constant.CACHE_GOODS_BY_ID, goodsParams.ID)
	// 查缓存
	value, ok := global.LocalCache.Get(key)
	if ok {
		goods = value.(model.Goods)
		return &goods, nil
	}
	err := global.DB.Where(&goodsParams).First(&goods).Error
	if err != nil {
		return nil, err
	} else {
		// 设置缓存
		global.LocalCache.SetNoExpire(key, goods)
		return &goods, nil
	}
}

// 查询商品列表
func (s *Shop) GetGoodsList(goodsParams *model.Goods) (*[]model.Goods, error) {
	var goodsList *[]model.Goods
	key := ""
	switch goodsParams.GoodsType {
	case constant.GOODS_TYPE_GENERAL:
		key = constant.CACHE_ALL_ENABLED_GOODS_GENERAL
	case constant.GOODS_TYPE_SUBSCRIBE:
		key = constant.CACHE_ALL_ENABLED_GOODS_SUBSCRIBE
	case constant.GOODS_TYPE_RECHARGE:
		key = constant.CACHE_ALL_ENABLED_GOODS_RECHARGE
	default:
	}
	value, ok := global.LocalCache.Get(key)
	if ok {
		goodsList = value.(*[]model.Goods)
		return goodsList, nil
	}
	//为普通用户屏蔽掉敏感字段 deliver_text：自动发货的内容
	err := global.DB.Where(&goodsParams).Omit("deliver_text").Order("goods_order").Find(&goodsList).Error
	if err != nil {
		return nil, err
	}
	global.LocalCache.SetNoExpire(key, goodsList)
	return goodsList, nil
}

// 更新商品（库存等）
func (s *Shop) UpdateGoods(goodsParams *model.Goods) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
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
	// 重新加载缓存
	_, _ = s.FirstGoods(&model.Goods{ID: goodsParams.ID})
	_, _ = s.GetGoodsList(&model.Goods{GoodsType: goodsParams.GoodsType, IsShow: true})

	return nil
}

// 支付
func (s *Shop) Purchase(sysOrder *model.Order) (*model.Order, error) {
	//fmt.Println("user_logic Purchase:")
	//Show(sysOrder)

	//根据支付id查询支付参数
	pay, err := PaySvc.FirstPayment(&model.Pay{ID: sysOrder.PayID})
	if err != nil {
		return nil, err
	}
	sysOrder.PayType = pay.PayType

	//判断支付方式
	switch sysOrder.PayType {
	case constant.PAY_TYPE_EPAY: // epay
		res, err := PaySvc.EpayPreByHTML(sysOrder, pay)
		if err != nil {
			return nil, err
		}
		sysOrder.TradeStatus = constant.ORDER_STATUS_WAIT_BUYER_PAY //订单状态：等待付款
		err = OrderSvc.UpdateOrder(sysOrder)                        //更新数据库
		if err != nil {
			return nil, err
		}
		sysOrder.PayInfo.EpayInfo = *res //返回易支付订单参数，采用易支付网页支付
		return sysOrder, nil
	case constant.PAY_TYPE_ALIPAY: // alipay
		//创建alipay client
		client, err := PaySvc.InitAlipayClient(pay)
		if err != nil {
			return nil, err
		}
		res, err := PaySvc.TradePreCreatePay(client, sysOrder, pay)
		if err != nil {
			return nil, err
		}
		sysOrder.TradeStatus = constant.ORDER_STATUS_WAIT_BUYER_PAY //订单状态：等待付款
		err = OrderSvc.UpdateOrder(sysOrder)                        //更新数据库
		if err != nil {
			return nil, err
		}
		sysOrder.PayInfo.AlipayInfo.QRCode = res.QRCode //返回用户qrcode
		return sysOrder, nil
	case constant.PAY_TYPE_BALANCE: // 余额支付
		err = UserSvc.UserBalancePayHandler(sysOrder)
		if err != nil {
			return nil, err
		}
		//修改支付状态
		sysOrder.TradeStatus = constant.ORDER_STATUS_TRADE_SUCCESS
		sysOrder.BuyerPayAmount = sysOrder.TotalAmount
		err = OrderSvc.PaymentSuccessfullyOrderHandler(sysOrder)
		if err != nil {
			return nil, err
		}
		return sysOrder, nil

	default:
		return nil, errors.New(constant.ERROR_INVALID_PAY_TYPE)
	}
}
