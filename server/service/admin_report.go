package service

import (
	"errors"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/other_plugin"
	"strings"
)

type Report struct{}

var AdminReportSvc *Report

// 获取所有数据库
func (r *Report) GetDB() (model.DbInfoReq, error) {
	var entities model.DbInfoReq
	if global.Config.SystemParams.DbType == "sqlite" {
		entities.DbType = "sqlite"
		entities.DatabaseList = append(entities.DatabaseList, global.Config.Sqlite.Path)
		return entities, nil
	}
	entities.DbType = "mysql"
	sql := "SELECT SCHEMA_NAME AS `database` FROM INFORMATION_SCHEMA.SCHEMATA"
	err := global.DB.Raw(sql).Scan(&entities.DatabaseList).Error
	return entities, err
}

// 获取数据库的所有表名
func (r *Report) GetTables(dbName string) ([]interface{}, error) {

	if global.Config.SystemParams.DbType == "mysql" && dbName == global.Config.Mysql.Dbname {
		var entities []model.DbMysqlTable
		var sql string
		sql = `select table_name as table_name from information_schema.tables where table_schema = ?`
		err := global.DB.Raw(sql, dbName).Scan(&entities).Error
		res := make([]interface{}, len(entities))
		for k, v := range entities {
			res[k] = v
		}
		return res, err
		//select table_name  from information_schema.tables where table_schema = 'airgo'
	} else if global.Config.SystemParams.DbType == "sqlite" && dbName == global.Config.Sqlite.Path {
		var entities []model.DbSqliteTable
		var sql string
		sql = `SELECT name FROM sqlite_master WHERE type='table' order by name`
		err := global.DB.Raw(sql).Scan(&entities).Error
		res := make([]interface{}, len(entities))
		for k, v := range entities {
			res[k] = v
		}
		return res, err
	}
	return nil, errors.New("未知数据库")
}

// 获取指定数据库指定数据表的所有字段名,类型值等(废弃）
func (r *Report) GetColumnByDB(dbName, tableName string) (data []interface{}, err error) {
	if global.Config.SystemParams.DbType == "mysql" {
		var entities []model.DbMysqlColumn
		var sql string
		sql = `
	SELECT COLUMN_NAME        column_name,
       DATA_TYPE          data_type,
       CASE DATA_TYPE
           WHEN 'longtext' THEN c.CHARACTER_MAXIMUM_LENGTH
           WHEN 'varchar' THEN c.CHARACTER_MAXIMUM_LENGTH
           WHEN 'double' THEN CONCAT_WS(',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE)
           WHEN 'decimal' THEN CONCAT_WS(',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE)
           WHEN 'int' THEN c.NUMERIC_PRECISION
           WHEN 'bigint' THEN c.NUMERIC_PRECISION
           ELSE '' END AS data_type_long,
       COLUMN_COMMENT     column_comment
	FROM INFORMATION_SCHEMA.COLUMNS c
	WHERE table_name = ?
	  AND table_schema = ?
	`
		err = global.DB.Raw(sql, tableName, dbName).Scan(&entities).Error
		res := make([]interface{}, len(entities))
		for k, v := range entities {
			res[k] = v
		}
		return res, err
	} else {
		var entities []model.DbSqliteColumn
		sql := `PRAGMA table_info(item)`
		newSql := strings.Replace(sql, "item", tableName, -1)
		err = global.DB.Raw(newSql).Scan(&entities).Error
		res := make([]interface{}, len(entities))
		for k, v := range entities {
			res[k] = v
		}
		return res, err
	}
}

// 获取指定数据表的所有字段名,类型值等
func (r *Report) GetColumnByReflect(tableName string) ([]string, map[string]interface{}, map[string]interface{}) {
	var table any
	table = model.StringAndStruct[tableName]
	return other_plugin.GetStructFieldMap(table)
}
