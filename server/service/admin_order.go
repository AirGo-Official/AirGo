package service

import (
	"fmt"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
	"time"
)

type AdminOrder struct{}

var AdminOrderSvc *AdminOrder

// 更新数据库订单
func (o *AdminOrder) UpdateOrder(order *model.Order) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&order).Error
	})
}

func (o *AdminOrder) OrderSummary(params *model.QueryParams) (*[]model.OrderSummary, error) {
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
