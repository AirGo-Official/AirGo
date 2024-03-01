package admin_api

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	//"github.com/ppoonk/AirGo/utils/encrypt_plugin"

	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/utils/response"
)

// 获取用户列表
func GetUserlist(ctx *gin.Context) {
	var params model.QueryParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
	}
	userList, err := userService.GetUserlist(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetUserlist error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetUserlist success", userList, ctx)
}

// 新建用户
func NewUser(ctx *gin.Context) {
	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = userService.NewUser(u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewUser error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewUser success", nil, ctx)
}

// 编辑用户信息
func UpdateUser(ctx *gin.Context) {
	var userParams model.User
	err := ctx.ShouldBind(&userParams)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//查找数据库用户数据
	userData, err := userService.FirstUser(&model.User{ID: userParams.ID})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateUser error:"+err.Error(), nil, ctx)
		return
	}
	admin_logic.Show(userParams)
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
		ReferrerCode:   userParams.ReferrerCode,
		Balance:        userParams.Balance,
		TgID:           userParams.TgID,
		RoleGroup:      userParams.RoleGroup,
		//Orders:         nil,
	}
	//处理密码,不为空且大于4位时更新，否则不修改
	if userParams.Password != "" && len(userParams.Password) > 4 {
		newUser.Password = encrypt_plugin.BcryptEncode(userParams.Password)
	}
	err = userService.SaveUser(&newUser)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateUser error:"+err.Error(), nil, ctx)
		return
	}
	//删除该用户token cache
	userService.DeleteUserCacheTokenByID(&userParams)
	response.OK("UpdateUser success", nil, ctx)
}

// 删除用户
func DeleteUser(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), err.Error(), ctx)
		return
	}
	// 删除用户
	err = userService.DeleteUser(&user)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("DeleteUser error:"+err.Error(), nil, ctx)
		return
	}
	//删除该用户token cache
	userService.DeleteUserCacheTokenByID(&user)
	response.OK("DeleteUser success", nil, ctx)
}
