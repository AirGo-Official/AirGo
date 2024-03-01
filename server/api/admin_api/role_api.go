package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/response"
)

// 获取角色列表
func GetRoleList(ctx *gin.Context) {
	res, err := roleService.GetRoleList()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetRoleList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetRoleList success", res, ctx)

}

// 修改角色信息
func UpdateRole(ctx *gin.Context) {
	var roleParams model.Role
	err := ctx.ShouldBind(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//处理角色
	err = roleService.UpdateRole(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateRole error:"+err.Error(), nil, ctx)
		return
	}
	//处理casbin
	err = casbinService.UpdateCasbinPolicy(&model.CasbinInfo{RoleID: roleParams.ID, CasbinItems: roleParams.CasbinItems})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateRole error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateRole success", nil, ctx)

}

// 新建角色
func NewRole(ctx *gin.Context) {
	var roleParams model.Role
	err := ctx.ShouldBind(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//处理角色
	err = roleService.NewRole(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewRole error:"+err.Error(), nil, ctx)
		return
	}
	//处理casbin
	err = casbinService.UpdateCasbinPolicy(&model.CasbinInfo{RoleID: roleParams.ID, CasbinItems: roleParams.CasbinItems})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewRole error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewRole success", nil, ctx)

}

// 删除角色
func DelRole(ctx *gin.Context) {
	var role model.Role
	err := ctx.ShouldBind(&role)
	if err != nil || role.ID == 1 || role.ID == 2 { //默认admin 普通用户不可删除
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = roleService.DelRole(&role)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DelRole error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DelRole success", nil, ctx)

}

// 获取全部权限
func GetAllPolicy(ctx *gin.Context) {
	res := casbinService.GetAllPolicy()
	response.OK("GetAllPolicy success", res, ctx)
}

// 获取权限列表ByRoleIds
func GetPolicyByID(ctx *gin.Context) {
	var casbinInfo model.CasbinInfo
	err := ctx.ShouldBind(&casbinInfo)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	res := casbinService.GetPolicyByRoleID(&model.CasbinInfo{RoleID: casbinInfo.RoleID})
	response.OK("GetPolicyByID success", res, ctx)
}
