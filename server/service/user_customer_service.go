package service

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type CustomerService struct{}

var CustomerServiceSvc *CustomerService

func (c *CustomerService) GetCustomerServiceListByUserID(params *model.QueryParams, uID int64) (*model.CommonDataResp, error) {
	var data model.CommonDataResp
	var csArr []model.CustomerService
	totalSql, dataSql := CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:]
	totalSql = totalSql[strings.Index(totalSql, "WHERE ")+6:]
	//拼接查询参数
	dataSql = fmt.Sprintf("user_id = %d AND %s", uID, dataSql)
	totalSql = fmt.Sprintf("user_id = %d AND %s", uID, totalSql)
	err := global.DB.Model(&model.CustomerService{}).Where(dataSql).Find(&csArr).Error
	if err != nil {
		return nil, err
	}
	err = global.DB.Model(&model.CustomerService{}).Where(totalSql).Count(&data.Total).Error
	if err != nil {
		return nil, err
	}
	data.Data = csArr
	return &data, nil
}

func (c *CustomerService) GetCustomerServiceList(csParams *model.CustomerService) (*[]model.CustomerService, error) {
	var csArr []model.CustomerService
	err := global.DB.Model(&model.CustomerService{}).Where(&csParams).Find(&csArr).Error
	return &csArr, err
}
func (c *CustomerService) DeleteCustomerService(csParams *model.CustomerService) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&csParams).Error
	})
}

func (c *CustomerService) FirstCustomerService(csParams *model.CustomerService) (*model.CustomerService, error) {
	var cs model.CustomerService
	err := global.DB.Model(&model.CustomerService{}).Where(&csParams).First(&cs).Error
	return &cs, err
}
func (c *CustomerService) UpdateCustomerService(id int64, values map[string]any) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.CustomerService{ID: id}).Updates(values).Error
	})
}
func (c *CustomerService) SaveCustomerService(csParams *model.CustomerService) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&csParams).Error
	})
}
func (c *CustomerService) CreateCustomerService(goods *model.Goods, order *model.Order) error {
	now := time.Now()
	cs := model.CustomerService{
		UserID:         order.UserID,
		UserName:       order.UserName,
		ServiceStatus:  true,
		ServiceStartAt: &now,

		//ServiceEndAt:   time.Now().AddDate(0, int(order.Duration), 0),
		IsRenew:        goods.IsRenew,
		RenewalAmount:  order.TotalAmount, //TODO 目前设置续费价格是订单结算时的价格（用户使用折扣码后的价格）要不要同步更新商品的价格呢？
		GoodsID:        order.GoodsID,
		Subject:        order.Subject,
		Des:            order.Des,
		Price:          order.Price,
		GoodsType:      order.GoodsType,
		Duration:       order.Duration,
		TotalBandwidth: goods.TotalBandwidth * 1024 * 1024 * 1024, // GB->MB->KB->B
		NodeConnector:  goods.NodeConnector,
		NodeSpeedLimit: goods.NodeSpeedLimit,

		//TrafficResetDay: int64(math.Min(float64(time.Now().Day()), 28)), //2月一般只有28天，流量重置日简单的设置为不超过28
		SubStatus: true,
		SubUUID:   uuid.NewV4(),
		UsedUp:    0,
		UsedDown:  0,
	}
	//处理流量重置日
	if goods.EnableTrafficReset {
		cs.TrafficResetDay = int64(math.Min(float64(time.Now().Day()), 28))
		// TODO 2月一般只有28天，流量重置日简单的设置为不超过28
		// TODO 如果对应的商品，开启重制功能的话,应该间隔30天，重制一次流量
	}
	//处理有效期
	switch order.Duration {
	case -1: // model.CustomerService  通过 gorm tag 指定字段 ServiceEndAt 默认值为 null，说明为不限时套餐
	default:
		end := now.AddDate(0, int(order.Duration), 0)
		cs.ServiceEndAt = &end
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&cs).Error
	})
}
