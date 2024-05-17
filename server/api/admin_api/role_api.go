package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetRoleList
// @Tags [admin api] role
// @Summary 获取角色列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/role/getRoleList [get]
func GetRoleList(ctx *gin.Context) {
	res, err := service.AdminRoleSvc.GetRoleList()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetRoleList error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetRoleList success", res, ctx)

}

// UpdateRole
// @Tags [admin api] role
// @Summary 修改角色信息
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Role true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/role/updateRole [post]
func UpdateRole(ctx *gin.Context) {
	var roleParams model.Role
	err := ctx.ShouldBind(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//处理角色
	err = service.AdminRoleSvc.UpdateRole(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateRole error:"+err.Error(), nil, ctx)
		return
	}
	//处理casbin
	err = service.AdminCasbinSvc.UpdateCasbinPolicy(&model.CasbinInfo{RoleID: roleParams.ID, CasbinItems: roleParams.CasbinItems})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateRole error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateRole success", nil, ctx)

}

// NewRole
// @Tags [admin api] role
// @Summary 新建角色
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Role true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/role/newRole [post]
func NewRole(ctx *gin.Context) {
	var roleParams model.Role
	err := ctx.ShouldBind(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//处理角色
	err = service.AdminRoleSvc.NewRole(&roleParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewRole error:"+err.Error(), nil, ctx)
		return
	}
	//处理casbin
	err = service.AdminCasbinSvc.UpdateCasbinPolicy(&model.CasbinInfo{RoleID: roleParams.ID, CasbinItems: roleParams.CasbinItems})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewRole error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewRole success", nil, ctx)

}

// DelRole
// @Tags [admin api] role
// @Summary 删除角色
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.Role true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/role/delRole [delete]
func DelRole(ctx *gin.Context) {
	var role model.Role
	err := ctx.ShouldBind(&role)
	if err != nil || role.ID == 1 || role.ID == 2 { //默认admin 普通用户不可删除
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminRoleSvc.DelRole(&role)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("DelRole error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DelRole success", nil, ctx)

}

// GetAllPolicy
// @Tags [admin api] role
// @Summary 获取全部权限
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/role/getAllPolicy [get]
func GetAllPolicy(ctx *gin.Context) {
	res := service.AdminCasbinSvc.GetAllPolicy()
	response.OK("GetAllPolicy success", res, ctx)
}

// GetPolicyByID
// @Tags [admin api] role
// @Summary 获取权限列表ByRoleId
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.CasbinInfo true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/role/getPolicyByID [post]
func GetPolicyByID(ctx *gin.Context) {
	var casbinInfo model.CasbinInfo
	err := ctx.ShouldBind(&casbinInfo)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	res := service.AdminCasbinSvc.GetPolicyByRoleID(&model.CasbinInfo{RoleID: casbinInfo.RoleID})
	response.OK("GetPolicyByID success", res, ctx)
}
