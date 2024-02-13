package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/common_logic"
	"github.com/ppoonk/AirGo/utils/response"
)

func NewAccessRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlCreate[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewAccessRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewAccessRoutes success", nil, ctx)

}

// 修改路由控制
func UpdateAccessRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlSave[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateAccessRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateAccessRoutes success", nil, ctx)

}

// 删除路由控制
func DeleteAccessRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = common_logic.CommonSqlDelete[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteAccessRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteAccessRoutes success", nil, ctx)

}

// 查询路由控制列表
func GetAccessRoutesList(ctx *gin.Context) {
	var p model.QueryParams
	err := ctx.ShouldBind(&p)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	list, total, err := common_logic.CommonSqlFindWithFieldParams(&p)
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
