package service

import (
	"errors"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/time_plugin"
	"strings"
)

type Traffic struct {
}

var TrafficSvc *Traffic

func (t *Traffic) GetSubTrafficList(params *model.CustomerService) (*[]model.UserTrafficLog, error) {
	// 默认查询当月
	start, end := time_plugin.GetFirstToTodayForMonth()
	var trafficList []model.UserTrafficLog
	err := global.DB.Model(&model.UserTrafficLog{}).Where("sub_user_id = ? AND created_at > ? AND created_at < ?", params.ID, start, end).Find(&trafficList).Error
	return &trafficList, err

}
func (t *Traffic) GetUserTraffic(params *model.QueryParams) (*model.UserTrafficLog, error) {
	var userTraffic model.UserTrafficLog
	var err error
	_, dataSql := CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:] //去掉`WHERE `
	if dataSql == "" {
		return nil, errors.New("invalid parameter")
	}
	//fmt.Println("dataSql:", dataSql)
	if global.Config.SystemParams.DbType == "mysql" {
		err = global.DB.
			Model(&model.UserTrafficLog{}).
			Where(dataSql).
			Select("sub_user_id, any_value(user_name) AS user_name, SUM(u) AS u, SUM(d) AS d").
			Group("sub_user_id").
			Find(&userTraffic).Error
	} else {
		err = global.DB.
			Model(&model.UserTrafficLog{}).
			Where(dataSql).
			Select("sub_user_id, user_name, SUM(u) AS u, SUM(d) AS d").
			Group("sub_user_id").
			Find(&userTraffic).Error
	}
	return &userTraffic, err
}

func (t *Traffic) GetAllUserTrafficRank(params *model.QueryParams) (*model.CommonDataResp, error) {
	//约定：params.FieldParamsList 数组前两项传时间，第三个开始传查询参数
	var userTraffic []model.UserTrafficLog
	var total int64
	var err error
	_, dataSql := CommonSqlFindSqlHandler(params)
	dataSql = dataSql[strings.Index(dataSql, "WHERE ")+6:] //去掉`WHERE `
	if dataSql == "" {
		return nil, errors.New("invalid parameter")
	}
	//fmt.Println("dataSql:", dataSql)

	if global.Config.SystemParams.DbType == "mysql" { //mysql only_full_group_by 问题
		err = global.DB.
			Model(&model.UserTrafficLog{}).
			Where(dataSql).
			Select("sub_user_id, any_value(user_name) AS user_name, SUM(u) u, SUM(d) AS d").
			Group("sub_user_id").
			Count(&total).
			Order(params.Pagination.OrderBy).
			Limit(int(params.Pagination.PageSize)).
			Offset((int(params.Pagination.PageNum) - 1) * int(params.Pagination.PageSize)).
			Find(&userTraffic).Error
	} else {
		err = global.DB.
			Model(&model.UserTrafficLog{}).
			Where(dataSql).
			Select("sub_user_id, user_name, SUM(u) u, SUM(d) AS d").
			Group("sub_user_id").Count(&total).
			Order(params.Pagination.OrderBy).
			Limit(int(params.Pagination.PageSize)).
			Offset((int(params.Pagination.PageNum) - 1) * int(params.Pagination.PageSize)).
			Find(&userTraffic).Error
	}

	if err != nil {
		return nil, err
	}
	return &model.CommonDataResp{
		Total: total,
		Data:  userTraffic,
	}, nil
}
