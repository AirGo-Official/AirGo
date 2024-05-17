package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetTables
// @Tags [admin api] report
// @Summary 获取数据库的所有表名
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.DbTableReq true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/report/getTables [post]
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
	res, err := service.AdminReportSvc.GetTables(dbName.DbName)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetTables error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetTables success", res, ctx)
}

// GetColumn
// @Tags [admin api] report
// @Summary 获取字段名,类型值
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.DbTableReq true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/report/getColumn [post]
func GetColumn(ctx *gin.Context) {
	var dbInfo model.DbTableReq
	err := ctx.ShouldBind(&dbInfo)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	m1, m2, m3 := service.AdminReportSvc.GetColumnByReflect(dbInfo.TableName)
	response.OK("GetColumn success", gin.H{
		"field_list":              m1,
		"field_chinese_name_list": m2,
		"field_type_list":         m3,
	}, ctx)
}

// ReportSubmit
// @Tags [admin api] report
// @Summary 获取报表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/report/reportSubmit [post]
func ReportSubmit(ctx *gin.Context) {
	var fieldParams model.QueryParams
	err := ctx.ShouldBind(&fieldParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//return
	res, total, err := service.CommonSqlFindWithFieldParams(&fieldParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("ReportSubmit error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ReportSubmit success", model.CommonDataResp{total, res}, ctx)
}
