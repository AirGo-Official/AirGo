package user_logic

import (
	"fmt"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

type CustomerService struct{}

var customerService *CustomerService

func (c *CustomerService) GetCustomerServiceListByUserID(params *model.QueryParams, uID int64) (*model.CommonDataResp, error) {
	var data model.CommonDataResp
	var csArr []model.CustomerService
	_, dataSql := common_logic.CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:]
	//拼接查询参数
	dataSql = fmt.Sprintf("user_id = %d AND %s", uID, dataSql)
	err := global.DB.Model(&model.CustomerService{}).Count(&data.Total).Where(dataSql).Find(&csArr).Error
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
	cs := model.CustomerService{
		UserID:          order.UserID,
		UserName:        order.UserName,
		ServiceStatus:   true,
		ServiceStartAt:  time.Now(),
		ServiceEndAt:    time.Now().AddDate(0, int(order.Duration), 0),
		IsRenew:         goods.IsRenew,
		RenewalAmount:   order.TotalAmount, //TODO 目前设置续费价格是订单结算时的价格（用户使用折扣码后的价格）
		GoodsID:         order.GoodsID,
		Subject:         order.Subject,
		Des:             order.Des,
		Price:           order.Price,
		GoodsType:       order.GoodsType,
		Duration:        order.Duration,
		TotalBandwidth:  goods.TotalBandwidth * 1024 * 1024 * 1024, // GB->MB->KB->B
		NodeConnector:   goods.NodeConnector,
		NodeSpeedLimit:  goods.NodeSpeedLimit,
		TrafficResetDay: int64(math.Min(float64(time.Now().Day()), 28)), //2月一般只有28天，流量重置日简单的设置为不超过28
		SubStatus:       true,
		SubUUID:         uuid.NewV4(),
		UsedUp:          0,
		UsedDown:        0,
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&cs).Error
	})
}
