package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/response"
	"github.com/gin-gonic/gin"
)

func GetRoleList(ctx *gin.Context) {
	//查询参数
	var roleParams model.PaginationParams
	err := ctx.ShouldBind(&roleParams)
	if err != nil {
		global.Logrus.Error("角色列表参数错误", err.Error())
		response.Fail("角色列表参数错误"+err.Error(), nil, ctx)
		return
	}
	res, err := service.GetRoleList(&roleParams)
	if err != nil {
		global.Logrus.Error("角色列表查询错误", err.Error())
		response.Fail("角色列表查询错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("角色列表查询成功", res, ctx)

}

// 修改角色信息
func ModifyRoleInfo(ctx *gin.Context) {
	var roleInfo model.Role
	err := ctx.ShouldBind(&roleInfo) //在 BindJson 之前有 ioutil （c.Request.Body）读取，报 EOF错误。
	if err != nil {
		global.Logrus.Error("修改角色信息参数错误", err.Error())
		response.Fail("修改角色信息参数错误"+err.Error(), nil, ctx)
		return
	}
	err = service.ModifyRoleInfo(&roleInfo)
	if err != nil {
		global.Logrus.Error("修改角色错误", err.Error())
		response.Fail("修改角色错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("修改角色成功", nil, ctx)

}

// 新建角色
func AddRole(ctx *gin.Context) {
	var role model.Role
	err := ctx.ShouldBind(&role)
	if err != nil {
		global.Logrus.Error("新建角色参数错误:", err.Error())
		response.Fail("新建角色参数错误", err.Error(), ctx)
		return
	}
	err = service.AddRole(&role)
	if err != nil {
		global.Logrus.Error("新建角色错误:", err.Error())
		response.Fail("新建角色错误", err, ctx)
		return
	}
	response.OK("新建角色成功", nil, ctx)

}

// 删除角色
func DelRole(ctx *gin.Context) {
	var role model.Role
	err := ctx.ShouldBind(&role)
	if err != nil || role.ID == 1 || role.ID == 2 { //默认admin 普通用户不可删除
		global.Logrus.Error("删除角色参数错误:", err.Error())
		response.Fail("删除角色参数错误", err.Error(), ctx)
		return
	}
	err = service.DelRole(role.ID)
	if err != nil {
		global.Logrus.Error("删除角色错误:", err.Error())
		response.Fail("删除角色错误", err, ctx)
		return
	}
	response.OK("删除角色成功", nil, ctx)

}
