package service

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type AdminCustomerService struct{}

var AdminCustomerServiceSvc *AdminCustomerService

func (c *AdminCustomerService) GetCustomerServiceList(csParams *model.CustomerService) (*[]model.CustomerService, error) {
	var csArr []model.CustomerService
	err := global.DB.Model(&model.CustomerService{}).Where(&csParams).Find(&csArr).Error
	return &csArr, err
}

// 更新客户服务
// gorm Save：
// 存在主键字段: 做插入操作
// 不存在主键字段：做更新操作
func (c *AdminCustomerService) UpdateCustomerService(csParams *model.CustomerService) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&csParams).Error
	})
}

// 删除
func (c *AdminCustomerService) DeleteCustomerService(csParams *model.CustomerService) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&model.CustomerService{}, csParams.ID).Error
	})
}

// 有效性检测任务
func (c *AdminCustomerService) SubExpirationCheck() error {
	return global.DB.Transaction(func(tx *gorm.DB) error {

		//服务有效性
		err := tx.Exec("UPDATE customer_service SET service_status = 0, sub_status = 0 WHERE service_end_at < ?", time.Now()).Error
		if err != nil {
			return err
		}
		//订阅有效性
		return tx.Exec("UPDATE customer_service SET sub_status = 0 WHERE ( used_up + used_down ) > total_bandwidth").Error
	})
}
func (c *AdminCustomerService) GetCustomerServiceListAlmostExpired() (*[]model.CustomerService, error) {
	var list []model.CustomerService
	//到期前3天
	d := time.Now()
	date := time.Date(d.Year(), d.Month(), d.Day()+10, d.Hour(), d.Minute(), d.Second(), 0, d.Location())
	err := global.DB.Model(&model.CustomerService{}).Where("service_end_at < ?", date).Find(&list).Error
	return &list, err
}

// 用户流量重置任务
func (c *AdminCustomerService) TrafficReset() error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		day := time.Now().Day()
		return tx.Exec("UPDATE customer_service SET used_up = 0, used_down = 0, sub_status = 1 WHERE traffic_reset_day = ? AND service_status = 1", day).Error
	})
}

// 更新客户已用流量
func (c *AdminCustomerService) UpdateCustomerServiceTrafficUsed(customerServiceArr *[]model.CustomerService, userIds []int64) error {
	var query []model.CustomerService
	err := global.DB.Where("id in ?", userIds).Select("id", "used_up", "used_down").Find(&query).Error
	if err != nil {
		return err
	}
	for item, _ := range query {
		query[item].UsedUp = query[item].UsedUp + (*customerServiceArr)[item].UsedUp
		query[item].UsedDown = query[item].UsedDown + (*customerServiceArr)[item].UsedDown
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"used_up", "used_down"}),
		}).Create(&query).Error
	})
}

// 更新客户流量记录
func (c *AdminCustomerService) UpdateCustomerServiceTrafficLog(userTrafficLogMap map[int64]model.UserTrafficLog, customerServerIDs []int64) error {
	var query []model.UserTrafficLog
	now := time.Now()
	//当日0点
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	err := global.DB.Where("created_at > ? AND sub_user_id IN ?", todayZero, customerServerIDs).Find(&query).Error
	if err != nil {
		return err
	}
	for k, _ := range query {
		if tl, ok := userTrafficLogMap[query[k].SubUserID]; ok { //已存在，叠加流量
			query[k].U += tl.U
			query[k].D += tl.D
			delete(userTrafficLogMap, query[k].SubUserID) //删除
		}
	}
	//不存在的数据，追加到最后面，一起插入数据库
	if len(userTrafficLogMap) > 0 {
		for k, _ := range userTrafficLogMap {
			query = append(query, userTrafficLogMap[k])
		}
	}
	if len(query) == 0 {
		return nil
	}
	return global.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"u", "d"}),
	}).Create(&query).Error
}

// 清理用户流量记录
func (c *AdminCustomerService) ClearCustomerServiceTrafficLog() error {
	y, m, _ := time.Now().Date()
	startTime := time.Date(y, m-2, 1, 0, 0, 0, 0, time.Local) //清除2个月之前的数据
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Where("created_at < ?", startTime).Delete(&model.UserTrafficLog{}).Error
	})
}
