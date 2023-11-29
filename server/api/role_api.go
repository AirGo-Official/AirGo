package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

func GetRoleList(ctx *gin.Context) {
	var roleParams model.PaginationParams
	err := ctx.ShouldBind(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetRoleList error:"+err.Error(), nil, ctx)
		return
	}
	res, err := service.GetRoleList(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetRoleList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetRoleList success", res, ctx)

}

// 修改角色信息
func ModifyRoleInfo(ctx *gin.Context) {
	var roleInfo model.Role
	err := ctx.ShouldBind(&roleInfo)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("ModifyRoleInfo error:"+err.Error(), nil, ctx)
		return
	}
	err = service.ModifyRoleInfo(&roleInfo)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("ModifyRoleInfo error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ModifyRoleInfo success", nil, ctx)

}

// 新建角色
func AddRole(ctx *gin.Context) {
	var role model.Role
	err := ctx.ShouldBind(&role)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("AddRole error:"+err.Error(), nil, ctx)
		return
	}
	err = service.AddRole(&role)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("AddRole error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("AddRole success", nil, ctx)

}

// 删除角色
func DelRole(ctx *gin.Context) {
	var role model.Role
	err := ctx.ShouldBind(&role)
	if err != nil || role.ID == 1 || role.ID == 2 { //默认admin 普通用户不可删除
		global.Logrus.Error(err.Error())
		response.Fail("DelRole error:"+err.Error(), nil, ctx)
		return
	}
	err = service.DelRole(role.ID)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DelRole error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DelRole success", nil, ctx)

}
