package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// 更新casbin权限
func UpdateCasbinPolicy(ctx *gin.Context) {
	var data model.ChangeRoleCasbinReq
	err := ctx.ShouldBind(&data)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateCasbinPolicy error:"+err.Error(), err, ctx)
		return
	}
	//前端传过来的没有处理，只有method，从数据库查询完整的rules，再更新casbin
	err = service.UpdateCasbinPolicy(&data)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateCasbinPolicy error:"+err.Error(), err, ctx)
		return
	}
	response.OK("UpdateCasbinPolicy success", nil, ctx)

}

// 获取全部权限
func GetAllPolicy(ctx *gin.Context) {
	res := service.GetAllPolicy()
	response.OK("GetAllPolicy success", res, ctx)
}

// 获取权限列表ByRoleIds
func GetPolicyByRoleID(ctx *gin.Context) {
	var casbinInfo model.CasbinInfo
	err := ctx.ShouldBind(&casbinInfo)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetPolicyByRoleID error:"+err.Error(), nil, ctx)
		return
	}
	res := service.GetPolicyByRoleID(&casbinInfo)
	response.OK("GetPolicyByRoleID success", res, ctx)
}
