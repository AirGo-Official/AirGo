package admin_api

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	//"github.com/ppoonk/AirGo/utils/encrypt_plugin"

	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetUserlist
// @Tags [admin api] user
// @Summary 获取用户列表
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/user/getUserlist [post]
func GetUserlist(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
	}
	userList, err := service.AdminUserSvc.GetUserlist(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetUserlist error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetUserlist success", userList, ctx)
}

// NewUser
// @Tags [admin api] user
// @Summary 新建用户
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.User true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/user/newUser [post]
func NewUser(ctx *gin.Context) {
	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.AdminUserSvc.NewUser(u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewUser error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewUser success", nil, ctx)
}

// UpdateUser
// @Tags [admin api] user
// @Summary 编辑用户信息
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.User true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/user/updateUser [post]
func UpdateUser(ctx *gin.Context) {
	var userParams model.User
	err := ctx.ShouldBind(&userParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//查找数据库用户数据
	userData, err := service.AdminUserSvc.FirstUser(&model.User{ID: userParams.ID})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateUser error:"+err.Error(), nil, ctx)
		return
	}
	var newUser = model.User{
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		//DeletedAt:      nil,
		ID:             userData.ID,
		UserName:       userParams.UserName,
		Password:       userData.Password,
		NickName:       userParams.NickName,
		Avatar:         userParams.Avatar,
		Enable:         userParams.Enable,
		InvitationCode: userParams.InvitationCode,
		ReferrerUserID: userParams.ReferrerUserID,
		Balance:        userParams.Balance,
		TgID:           userParams.TgID,
		RoleGroup:      userParams.RoleGroup,
		//Orders:         nil,
	}
	//处理密码,不为空且大于4位时更新，否则不修改
	if userParams.Password != "" && len(userParams.Password) > 4 {
		newUser.Password = encrypt_plugin.BcryptEncode(userParams.Password)
	}
	err = service.AdminUserSvc.SaveUser(&newUser)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateUser error:"+err.Error(), nil, ctx)
		return
	}
	//删除该用户token cache
	service.AdminUserSvc.DeleteUserCacheTokenByID(&userParams)
	response.OK("UpdateUser success", nil, ctx)
}

// DeleteUser
// @Tags [admin api] user
// @Summary 删除用户
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.User true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/user/deleteUser [delete]
func DeleteUser(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), err.Error(), ctx)
		return
	}
	// 删除用户
	err = service.AdminUserSvc.DeleteUser(&user)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("DeleteUser error:"+err.Error(), nil, ctx)
		return
	}
	//删除该用户token cache
	service.AdminUserSvc.DeleteUserCacheTokenByID(&user)
	response.OK("DeleteUser success", nil, ctx)
}

// UserSummary
// @Tags [admin api] user
// @Summary 用户统计
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.QueryParams true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/admin/user/userSummary [post]
func UserSummary(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	res, err := service.AdminUserSvc.UserSummary(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	response.OK("UserSummary success", res, ctx)
}
