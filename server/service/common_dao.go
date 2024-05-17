package service

import (
	"fmt"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
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

// 通用查询
func CommonSqlFindWithFieldParams(fieldParams *model.QueryParams) (any, int64, error) {
	totalSql, dataSql := CommonSqlFindSqlHandler(fieldParams)
	//fmt.Println("totalSql:", totalSql)
	//fmt.Println("dataSql:", dataSql)
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

func CommonSqlFindSqlHandler(fieldParams *model.QueryParams) (string, string) {
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
	sql := strings.Join(sqlArr, " ")
	//判断sql是否为空
	sqlTemp := strings.TrimSpace(sql)
	sqlTemp = strings.ReplaceAll(sql, "'", "")
	if sqlTemp == "" {
		sql = "id > 0"
	}
	totalSql := fmt.Sprintf("SELECT COUNT(id) FROM `%s` WHERE %s", fieldParams.TableName, sql)
	dataSql := fmt.Sprintf("SELECT * FROM `%s` WHERE %s ",
		fieldParams.TableName, sql)
	//排序
	var orderBy = "id ASC"
	if fieldParams.Pagination.OrderBy != "" {
		orderBy = fieldParams.Pagination.OrderBy
	}
	//分页
	if fieldParams.Pagination.PageNum != 0 && fieldParams.Pagination.PageSize != 0 {
		dataSql = fmt.Sprintf("%s ORDER BY %s LIMIT %d OFFSET %d",
			dataSql, orderBy, fieldParams.Pagination.PageSize, fieldParams.Pagination.PageSize*(fieldParams.Pagination.PageNum-1))
	}
	return totalSql, dataSql
}

// 通用删除
func CommonSqlDelete[T1, T2 any](params T2) error {
	var m T1
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		if reflect.TypeOf(params).String() == reflect.String.String() {
			err = tx.Where(params).Delete(&m).Error
		} else {
			err = tx.Delete(&params).Error
		}
		return err
	})
}

// 通用更新 save
func CommonSqlSave[T1 any](data T1) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&data).Error
	})
}

// 通用更新 update
func CommonSqlUpdate[T1 any](params T1, data map[string]any) error {
	var m T1
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&m).Where(params).Updates(&data).Error
	})
}

// 通用更新多行数据，针对数据中多行已存在数据进行更新
func CommonSqlUpdateMultiLine[T1 any](data T1, name string, columns []string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: name}},     //冲突字段
			DoUpdates: clause.AssignmentColumns(columns), //需要更新的字段
		}).Create(&data).Error
	})
}

// 通用增加
func CommonSqlCreate[T1 any](data T1) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&data).Error
	})
}
