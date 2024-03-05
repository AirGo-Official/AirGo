package admin_logic

import (
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/user_logic"
	"gorm.io/gorm"
	"time"
)

type Order struct{}

var orderService *Order

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

func (o *Order) OrderSummary(params *model.QueryParams) (*[]model.OrderSummary, error) {
	//处理查询时间
	var startTime, endTime time.Time
	startTime, err := time.Parse("2006-01-02 15:04:05", params.FieldParamsList[0].ConditionValue)
	if err != nil {
		return nil, err
	}
	endTime, _ = time.Parse("2006-01-02 15:04:05", params.FieldParamsList[1].ConditionValue)
	if err != nil {
		return nil, err
	}
	const (
		sql1 = `SELECT
DATE(created_at) as date,
COUNT(id) AS order_total,
SUM(CAST(buyer_pay_amount AS decimal(10,2))) AS income_total,
COUNT(goods_type = 'subscribe' OR NULL) AS subscribe_total,
COUNT(goods_type = 'general' OR NULL) AS general_total,
COUNT(goods_type = 'recharge' OR NULL) AS recharge_total
	`
		sql2 = "FROM `order`"
	)
	sql3 := fmt.Sprintf(" WHERE created_at > '%s' AND created_at < '%s' AND trade_status = 'TRADE_SUCCESS'  GROUP BY date", startTime, endTime)
	var orderSummary []model.OrderSummary
	err = global.DB.
		Raw(sql1 + sql2 + sql3).
		Scan(&orderSummary).
		Error
	//fmt.Println("result:", orderSummary, "err:", err)
	return &orderSummary, err

}
