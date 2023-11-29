package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// 获取数据库的所有数据库名
func GetDB(ctx *gin.Context) {
	res, err := service.GetDB()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetDB error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetDB success", res, ctx)
}

// 获取数据库的所有表名
func GetTables(ctx *gin.Context) {
	var dbName model.DbNameAndTableReq
	err := ctx.ShouldBind(&dbName)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetTables error:"+err.Error(), nil, ctx)
		return
	}
	if dbName.Database == "" {
		response.Fail("Database name is empty", nil, ctx)
		return
	}
	res, err := service.GetTables(dbName.Database)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetTables error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetTables success", res, ctx)
}

// 获取字段名,类型值
func GetColumn(ctx *gin.Context) {
	var dbNameAndTable model.DbNameAndTableReq
	err := ctx.ShouldBind(&dbNameAndTable)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetColumn error:"+err.Error(), nil, ctx)
		return
	}
	m1, m2, m3 := service.GetColumnByReflect(dbNameAndTable.TableName)
	response.OK("GetColumn success", gin.H{
		"field_list":              m1,
		"field_chinese_name_list": m2,
		"field_type_list":         m3,
	}, ctx)
}

// 获取报表
func ReportSubmit(ctx *gin.Context) {
	var fieldParams model.FieldParamsReq
	err := ctx.ShouldBind(&fieldParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("ReportSubmit error:"+err.Error(), nil, ctx)
		return
	}
	res, total, err := service.CommonSqlFindWithFieldParams(fieldParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("ReportSubmit error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ReportSubmit success", gin.H{
		"table_data": res,
		"total":      total,
	}, ctx)
}
