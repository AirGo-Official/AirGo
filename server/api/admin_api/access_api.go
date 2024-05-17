package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// NewAccessRoutes
// @Tags [admin api] access
// @Summary 新建访问控制
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Access true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/access/newAccess [post]
func NewAccessRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlCreate[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewAccessRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewAccessRoutes success", nil, ctx)

}

// UpdateAccessRoutes
// @Tags [admin api] access
// @Summary 修改路由控制
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Access true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/access/updateAccess [post]
func UpdateAccessRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateAccessRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateAccessRoutes success", nil, ctx)

}

// DeleteAccessRoutes
// @Tags [admin api] access
// @Summary 删除访问控制
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Access true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/access/deleteAccess [delete]
func DeleteAccessRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlDelete[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteAccessRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteAccessRoutes success", nil, ctx)

}

// GetAccessRoutesList
// @Tags [admin api] access
// @Summary 查询路由控制列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/access/getAccessList [post]
func GetAccessRoutesList(ctx *gin.Context) {
	var p model.QueryParams
	err := ctx.ShouldBind(&p)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	list, total, err := service.CommonSqlFindWithFieldParams(&p)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetAccessRoutesList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetAccessRoutesList success", gin.H{
		"total": total,
		"data":  list,
	}, ctx)

}
