package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/response"
)

// GetUserInfo
// @Tags [customer api] user
// @Summary 获取自身信息
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/user/getUserInfo [get]
func GetUserInfo(ctx *gin.Context) {
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	user, err := service.UserSvc.FirstUser(&model.User{ID: uIDInt})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetUserInfo error:"+err.Error(), nil, ctx)
		return
	}
	user.Password = ""
	response.OK("GetUserInfo success", user, ctx)
}

// ChangeUserPassword
// @Tags [customer api] user
// @Summary 修改密码
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.UserChangePasswordRequest true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/user/changeUserPassword [post]
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
	err = service.UserSvc.UpdateUser(&model.User{ID: uIDInt}, map[string]any{
		"password": encrypt_plugin.BcryptEncode(u.Password),
	})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ChangeUserPassword error:"+err.Error(), nil, ctx)
		return
	}
	// TODO 该用户token校验依然有效，需优化
	service.UserSvc.DeleteUserCacheTokenByID(&model.User{
		ID: uIDInt,
	})
	response.OK("ChangeUserPassword success", nil, ctx)
}

// ChangeUserAvatar
// @Tags [customer api] user
// @Summary 修改头像
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.UserChangeAvatarRequest true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/user/changeUserAvatar [post]
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
	err = service.UserSvc.UpdateUser(&model.User{ID: uIDInt}, map[string]any{
		"avatar": params.Avatar,
	})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ChangeUserAvatar error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ChangeUserAvatar success", nil, ctx)
}

// ClockIn
// @Tags [customer api] user
// @Summary 打卡
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/user/clockIn [get]
func ClockIn(ctx *gin.Context) {
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	index, _, err := service.UserSvc.ClockIn(uIDInt)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ClockIn error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ClockIn success", model.CommonDataResp{
		Total: 1,
		Data:  index, //返回奖品索引
	}, ctx)

}

// SetUserNotice
// @Tags [customer api] user
// @Summary 设置用户通知
// @Produce json
// @Param Authorization header string false "Bearer 用户token"
// @Param data body model.User true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/customer/user/setUserNotice [post]
func SetUserNotice(ctx *gin.Context) {
	uIDInt, ok := api.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	err = service.UserSvc.UpdateUser(&model.User{ID: uIDInt}, map[string]any{
		"tg_id":                       u.TgID,
		"enable_tg_bot":               u.EnableTGBot,
		"enable_email":                u.EnableEmail,
		"enable_web_mail":             u.EnableWebMail,
		"when_purchased":              u.WhenPurchased,
		"when_service_almost_expired": u.WhenServiceAlmostExpired,
		"when_balance_changed":        u.WhenBalanceChanged,
	})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("SetUserNotice error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("SetUserNotice success", nil, ctx)
}
