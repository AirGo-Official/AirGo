package service

import (
	"AirGo/global"
	"AirGo/model"
	"time"
)

// 更新数据库订单
func UpdateOrder(order *model.Orders) error {
	err := global.DB.Save(&order).Error
	return err
}

// 获取全部订单,分页获取
func GetAllOrder(orderParams *model.PaginationParams) (*model.OrdersWithTotal, error) {
	var startTime, endTime time.Time
	//时间格式转换
	if len(orderParams.Date) == 2 {
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", orderParams.Date[0], time.Local)
		endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", orderParams.Date[1], time.Local)
	} else {
		//默认前1个月数据
		endTime = time.Now().Local()
		startTime = endTime.AddDate(0, 0, -30)
	}
	var ordersWithTotal model.OrdersWithTotal
	var err error
	if orderParams.Search != "" {
		err = global.DB.Model(&model.Orders{}).Count(&ordersWithTotal.Total).Where("out_trade_no = ? and created_at > ? and created_at < ?", orderParams.Search, startTime, endTime).Limit(int(orderParams.PageSize)).Offset((int(orderParams.PageNum) - 1) * int(orderParams.PageSize)).Order("id desc").Find(&ordersWithTotal.OrderList).Error
	} else {
		err = global.DB.Model(&model.Orders{}).Count(&ordersWithTotal.Total).Where("created_at > ? and created_at < ?", startTime, endTime).Limit(int(orderParams.PageSize)).Offset((int(orderParams.PageNum) - 1) * int(orderParams.PageSize)).Order("id desc").Find(&ordersWithTotal.OrderList).Error
	}
	//fmt.Println("ordersWithTotal:", ordersWithTotal)
	return &ordersWithTotal, err
}

// 获取订单统计
func GetMonthOrderStatistics(orderParams *model.PaginationParams) (*model.OrderStatistics, error) {
	var startTime, endTime time.Time
	//时间格式转换
	if len(orderParams.Date) == 2 {
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", orderParams.Date[0], time.Local)
		endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", orderParams.Date[1], time.Local)
	} else {
		//默认前1个月数据
		endTime = time.Now().Local()
		startTime = endTime.AddDate(0, 0, -30)
	}
	var orderStatistic model.OrderStatistics
	err := global.DB.Model(&model.Orders{}).Where("created_at > ? and created_at < ?", startTime, endTime).Select("sum(receipt_amount) as total_amount").Find(&orderStatistic).Count(&orderStatistic.Total).Error
	if err != nil {
		global.Logrus.Error("获取月订单统计 error:", err.Error())
		return &model.OrderStatistics{}, err
	}
	return &orderStatistic, err
}
