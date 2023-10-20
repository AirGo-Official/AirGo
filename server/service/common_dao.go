package service

import (
	"AirGo/global"
	"AirGo/model"
	"fmt"
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
		err = global.DB.Model(&modelType).Where(params).Count(&total).Find(&res).Error

	} else {
		err = global.DB.Model(&modelType).Where(&params).Count(&total).Find(&res).Error
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
		err = global.DB.Model(&modelType).Where(params).Count(&total).Limit(int(paginationParams.PageSize)).Offset((int(paginationParams.PageNum) - 1) * int(paginationParams.PageSize)).Find(&res).Error

	} else {
		err = global.DB.Model(&modelType).Where(&params).Count(&total).Limit(int(paginationParams.PageSize)).Offset((int(paginationParams.PageNum) - 1) * int(paginationParams.PageSize)).Find(&res).Error
	}
	return res, total, err
}

// 通用查询,支持更多字段参数
func CommonSqlFindWithFieldParams(fieldParams model.FieldParamsReq) (any, int64, error) {
	var sqlArr []string
	for _, v := range fieldParams.FieldParamsList {
		if v.Field == "" || v.Condition == "" || v.ConditionValue == "" {
			continue
		}
		if v.Condition == "like" {
			sqlArr = append(sqlArr, v.Field+"  "+v.Condition+" '%"+v.ConditionValue+"%'")

		} else {
			sqlArr = append(sqlArr, v.Field+" "+v.Condition+" '"+v.ConditionValue+"'")
		}
	}
	totalSql := "select count(id) from " + fieldParams.TableName + " where " + strings.Join(sqlArr, " and ")
	dataSql := "select * from " + fieldParams.TableName + " where " + strings.Join(sqlArr, " and ") + " limit " + fmt.Sprintf("%d", fieldParams.PaginationParams.PageSize) + " offset " + fmt.Sprintf("%d", fieldParams.PaginationParams.PageNum-1)

	//fmt.Println("totalSql:", totalSql)
	//fmt.Println("dataSql:", dataSql)

	var data any
	data = model.StringAndSlice[fieldParams.TableName]

	var total int64
	err := global.DB.Raw(totalSql).Scan(&total).Error
	err = global.DB.Raw(dataSql).Scan(&data).Error
	return data, total, err
}

// 通用删除
func CommonSqlDelete[T1, T2 any](modelType T1, params T2) error {
	var err error
	if reflect.TypeOf(params).String() == reflect.String.String() {
		err = global.DB.Where(params).Delete(&modelType).Error
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
func CommonSqlUpdate[T1, T2 any](modelType T1, data T2, params string) error {
	var err error
	//保存多列
	err = global.DB.Model(&modelType).Where(params).Updates(&data).Error
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
