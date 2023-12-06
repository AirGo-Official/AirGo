package service

import (
	"fmt"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
)

// 通用查询
func CommonSqlFind[T1, T2, T3 any](params T2) (T3, int64, error) {
	var res T3
	var err error
	var modelType T1
	var total int64
	if reflect.TypeOf(params).String() == reflect.String.String() {
		err = global.DB.Model(&modelType).Where(params).Find(&res).Error

	} else {
		err = global.DB.Model(&modelType).Count(&total).Where(&params).Find(&res).Error
	}
	return res, total, err
}
func CommonSqlFirst[T1, T2, T3 any](params T2) (T3, int64, error) {
	var res T3
	var err error
	var modelType T1
	var total int64
	if reflect.TypeOf(params).String() == reflect.String.String() {
		err = global.DB.Model(&modelType).Where(params).First(&res).Error

	} else {
		err = global.DB.Model(&modelType).Where(&params).Count(&total).First(&res).Error
	}
	return res, total, err
}
func CommonSqlLast[T1, T2, T3 any](params T2) (T3, int64, error) {
	var res T3
	var err error
	var modelType T1
	var total int64
	if reflect.TypeOf(params).String() == reflect.String.String() {
		err = global.DB.Model(&modelType).Where(params).First(&res).Error

	} else {
		err = global.DB.Model(&modelType).Count(&total).Where(&params).Last(&res).Error
	}
	return res, total, err
}

// 通用查询带分页参数
func CommonSqlFindWithPagination[T1, T2, T3 any](params T2, paginationParams model.PaginationParams) (T3, int64, error) {
	var res T3
	var err error
	var modelType T1
	var total int64
	if reflect.TypeOf(params).String() == reflect.String.String() {
		err = global.DB.Model(&modelType).Count(&total).Where(params).Limit(int(paginationParams.PageSize)).Offset((int(paginationParams.PageNum) - 1) * int(paginationParams.PageSize)).Find(&res).Error

	} else {
		err = global.DB.Model(&modelType).Count(&total).Where(&params).Limit(int(paginationParams.PageSize)).Offset((int(paginationParams.PageNum) - 1) * int(paginationParams.PageSize)).Find(&res).Error
	}
	return res, total, err
}

// 通用查询
func CommonSqlFindWithFieldParams(fieldParams *model.FieldParamsReq) (any, int64, error) {
	totalSql, dataSql := CommonSqlFindSqlHandler(fieldParams)
	var data any
	data = model.StringAndSlice[fieldParams.TableName]
	var total int64
	err := global.DB.Raw(totalSql).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = global.DB.Raw(dataSql).Scan(&data).Error
	return data, total, err
}

func CommonSqlFindSqlHandler(fieldParams *model.FieldParamsReq) (string, string) {
	var sqlArr []string
	for _, v := range fieldParams.FieldParamsList {
		if v.Condition == "" && v.ConditionValue == "" {
			continue
		}
		switch v.Condition {
		case "like":
			sqlArr = append(sqlArr, v.Operator+" "+v.Field+"  "+v.Condition+" '%"+v.ConditionValue+"%'")
		default:
			sqlArr = append(sqlArr, v.Operator+" "+v.Field+" "+v.Condition+" '"+v.ConditionValue+"'")
		}
	}
	//排序
	var orderBy = "id ASC"
	if fieldParams.Pagination.OrderBy != "" {
		orderBy = fieldParams.Pagination.OrderBy
	}
	sql := strings.Join(sqlArr, " ")
	//判断sql是否为空
	sqlTemp := strings.TrimSpace(sql)
	sqlTemp = strings.ReplaceAll(sql, "'", "")
	if sqlTemp == "" {
		sql = "id > 0"
	}
	totalSql := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE %s", fieldParams.TableName, sql)
	dataSql := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY %s LIMIT %d OFFSET %d", fieldParams.TableName, sql, orderBy, fieldParams.Pagination.PageSize, fieldParams.Pagination.PageSize*(fieldParams.Pagination.PageNum-1))
	return totalSql, dataSql
}

// 通用删除
func CommonSqlDelete[T1, T2 any](params T2) error {
	var err error
	var m T1
	if reflect.TypeOf(params).String() == reflect.String.String() {
		err = global.DB.Where(params).Delete(&m).Error
	} else {
		err = global.DB.Delete(&params).Error
	}
	return err
}

// 通用更新 save
func CommonSqlSave[T1 any](data T1) error {
	//保存所有列
	return global.DB.Save(&data).Error
}

// 通用更新 update
func CommonSqlUpdate[T1, T2 any](data T2, params string) error {
	var err error
	var m T1
	//保存多列
	err = global.DB.Model(&m).Where(params).Updates(&data).Error
	return err
}

// 通用更新多行数据，针对数据中多行已存在数据进行更新
func CommonSqlUpdateMultiLine[T1 any](data T1, name string, columns []string) error {
	var err error
	err = global.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: name}},
		DoUpdates: clause.AssignmentColumns(columns),
	}).Create(&data).Error
	return err
}

// 通用增加
func CommonSqlCreate[T1 any](data T1) error {
	var err error
	err = global.DB.Create(&data).Error
	return err
}
