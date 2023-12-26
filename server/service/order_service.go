package service

import (
	"fmt"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"strings"
	"time"
)

// 更新数据库订单
func UpdateOrder(order *model.Orders) error {
	err := global.DB.Save(&order).Error
	return err
}

func GetUserOrders(params *model.FieldParamsReq, uIDInt int64) (*model.CommonDataResp, error) {
	var data model.CommonDataResp
	var orderList []model.Orders
	_, dataSql := CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:]
	//拼接查询参数
	dataSql = fmt.Sprintf("user_id = %d AND %s", uIDInt, dataSql)
	err := global.DB.Model(&model.Orders{}).Count(&data.Total).Where(dataSql).Find(&orderList).Error
	if err != nil {
		global.Logrus.Error("GetUserOrders error:", err.Error())
		return nil, err
	}
	data.Data = orderList
	return &data, nil

}

func GetMonthOrderStatistics(params *model.FieldParamsReq) (*model.OrderStatistics, error) {
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
	return GetOrderStatistics(startTime, endTime)
}

func GetOrderStatistics(startTime, endTime time.Time) (*model.OrderStatistics, error) {
	var orderStatistic model.OrderStatistics
	err := global.DB.Model(&model.Orders{}).Where("created_at > ? and created_at < ?", startTime, endTime).Select("sum(receipt_amount) as total_amount").Find(&orderStatistic).Count(&orderStatistic.Total).Error
	if err != nil {
		global.Logrus.Error("获取月订单统计 error:", err.Error())
		return &model.OrderStatistics{}, err
	}
	return &orderStatistic, err
}

func DeleteUserAllOrder(user *model.User) error {
	return global.DB.Where("user_id = ?", user.ID).Delete(&model.Orders{}).Error
}

// 处理需要发货的订单
func DeliverOrder(order *model.Orders) error {
	//查询商品信息
	goods, err := FindGoodsByGoodsID(order.GoodsID)
	if err != nil {
		return err
	}
	switch goods.DeliverType {
	case model.DeliverTypeNone:
	case model.DeliverTypeManual:
	case model.DeliverTypeAuto:
		order.DeliverText = goods.DeliverText
	}
	return UpdateOrder(order) //更新数据库订单状态
}

// 处理支付成功后的订单
func PaymentSuccessfullyOrderHandler(order *model.Orders) {
	switch order.GoodsType {
	case model.GoodsTypeGeneral: //普通商品
		DeliverOrder(order)
	case model.GoodsTypeSubscribe: //订阅
		global.GoroutinePool.Submit(func() {
			_ = UpdateUserSubscribe(order) //更新用户订阅信息
		})
		global.GoroutinePool.Submit(func() {
			_ = RemainHandle(order.UserID, order.RemainAmount) //处理用户余额
		})
		global.GoroutinePool.Submit(func() {
			_ = UpdateOrder(order) //更新数据库订单状态
		})
	case model.GoodsTypeRecharge: //充值
		global.GoroutinePool.Submit(func() {
			_ = RechargeHandle(order) //处理用户余额
		})
		global.GoroutinePool.Submit(func() {
			_ = UpdateOrder(order) //更新数据库订单状态
		})
	}
	//通知
	global.GoroutinePool.Submit(func() {
		UnifiedPushMessage(fmt.Sprintf("用户：%s\n购买：%s\n销售价格：%s\n订单金额：%s\n支付方式：%s", order.UserName, order.Subject, order.Price, order.TotalAmount, order.PayType))
	})
}
