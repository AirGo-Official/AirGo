package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

var reportService admin_logic.Report

// 获取数据库的所有表名
func GetTables(ctx *gin.Context) {
	var dbName model.DbTableReq
	err := ctx.ShouldBind(&dbName)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	if dbName.DbName == "" {
		response.Fail("Database name is empty", nil, ctx)
		return
	}
	res, err := reportService.GetTables(dbName.DbName)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetTables error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetTables success", res, ctx)
}

// 获取字段名,类型值
func GetColumn(ctx *gin.Context) {
	var dbInfo model.DbTableReq
	err := ctx.ShouldBind(&dbInfo)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	m1, m2, m3 := reportService.GetColumnByReflect(dbInfo.TableName)
	response.OK("GetColumn success", gin.H{
		"field_list":              m1,
		"field_chinese_name_list": m2,
		"field_type_list":         m3,
	}, ctx)
}

// 获取报表
func ReportSubmit(ctx *gin.Context) {
	var fieldParams model.QueryParams
	err := ctx.ShouldBind(&fieldParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//return
	res, total, err := common_logic.CommonSqlFindWithFieldParams(&fieldParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("ReportSubmit error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ReportSubmit success", gin.H{
		"data":  res,
		"total": total,
	}, ctx)
}
