package admin_logic

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/user_logic"
	"gorm.io/gorm"
	"time"
)

type Order struct{}

var orderService *Order

func (o *Order) GetMonthOrderStatistics(params *model.QueryParams) (*model.OrderStatistics, error) {
	var startTime, endTime time.Time
	//时间格式转换
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", params.FieldParamsList[0].ConditionValue, time.Local)
	if err != nil {
		return nil, err
	}
	endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", params.FieldParamsList[1].ConditionValue, time.Local)
	if err != nil {
		return nil, err
	}
	return o.GetOrderStatistics(startTime, endTime)
}

func (o *Order) GetOrderStatistics(startTime, endTime time.Time) (*model.OrderStatistics, error) {
	var orderStatistic model.OrderStatistics
	err := global.DB.
		Model(&model.Order{}).
		Where("created_at > ? and created_at < ?", startTime, endTime).
		Select("sum(buyer_pay_amount) as total_amount").
		Find(&orderStatistic).
		Count(&orderStatistic.Total).
		Error
	if err != nil {
		return &model.OrderStatistics{}, err
	}
	return &orderStatistic, err
}

// 更新数据库订单
func (o *Order) UpdateOrder(order *model.Order) error {
	var userOrderService user_logic.Order
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&order).Error
	})
	if err != nil {
		return err
	}
	//如果订单状态是 支付成功 或 交易结束，将订单进行处理
	if order.TradeStatus == constant.ORDER_STATUS_TRADE_SUCCESS || order.TradeStatus == constant.ORDER_STATUS_TRADE_FINISHED {
		userOrderService.DeleteOneOrderFromCache(order) //删除缓存
		return userOrderService.PaymentSuccessfullyOrderHandler(order)
	} else {
		userOrderService.UpdateOneOrderToCache(order) //更新缓存
	}
	return nil
}
