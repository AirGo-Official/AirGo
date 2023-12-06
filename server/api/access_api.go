package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

func NewRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewRoutes error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlCreate[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewRoutes success", nil, ctx)

}

// 修改路由控制
func UpdateRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateRoutes error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlSave[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateRoutes success", nil, ctx)

}

// 删除路由控制
func DeleteRoutes(ctx *gin.Context) {
	var acc model.Access
	err := ctx.ShouldBind(&acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteRoutes error:"+err.Error(), nil, ctx)
		return
	}
	err = service.CommonSqlDelete[model.Access](acc)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DeleteRoutes error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteRoutes success", nil, ctx)

}

// 查询路由控制列表
func GetRoutesList(ctx *gin.Context) {
	var p model.FieldParamsReq
	err := ctx.ShouldBind(&p)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetRoutesList error:"+err.Error(), nil, ctx)
		return
	}
	list, total, err := service.CommonSqlFindWithFieldParams(&p)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetRoutesList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetRoutesList success", gin.H{
		"total": total,
		"data":  list,
	}, ctx)

}
