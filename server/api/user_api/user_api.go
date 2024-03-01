package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/response"
)

// 获取自身信息
func GetUserInfo(ctx *gin.Context) {
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	user, err := userService.FirstUser(&model.User{ID: uIDInt})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetUserInfo error:"+err.Error(), nil, ctx)
		return
	}
	user.Password = ""
	response.OK("GetUserInfo success", user, ctx)
}

// 修改密码
func ChangeUserPassword(ctx *gin.Context) {
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	var u model.UserChangePasswordRequest
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = userService.UpdateUser(&model.User{ID: uIDInt}, map[string]any{
		"password": encrypt_plugin.BcryptEncode(u.Password),
	})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ChangeUserPassword error:"+err.Error(), nil, ctx)
		return
	}
	// TODO 该用户token校验依然有效，需优化
	userService.DeleteUserCacheTokenByID(&model.User{
		ID: uIDInt,
	})
	response.OK("ChangeUserPassword success", nil, ctx)
}

// 修改头像
func ChangeUserAvatar(ctx *gin.Context) {
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	var params model.UserChangeAvatarRequest
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = userService.UpdateUser(&model.User{ID: uIDInt}, map[string]any{
		"avatar": params.Avatar,
	})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ChangeUserAvatar error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ChangeUserAvatar success", nil, ctx)
}

// 打卡
func ClockIn(ctx *gin.Context) {
	// TODO
}
