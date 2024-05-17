package public_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/other_plugin"
	"github.com/ppoonk/AirGo/utils/response"
	"strings"
	"time"
)

// Register
// @Tags [public api] user
// @Summary 用户注册
// @Produce json
// @Param data body model.UserRegister true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/public/user/register [post]
func Register(ctx *gin.Context) {
	if !global.Server.Website.EnableRegister {
		response.Fail("Registration closed", nil, ctx)
		return
	}
	var u model.UserRegister
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//判断邮箱后缀，目前注册时用户名和邮箱后缀是分开的，如：user_name = 123, email_suffix = @qq.com
	ok := other_plugin.In(u.EmailSuffix, strings.Fields(global.Server.Website.AcceptableEmailSuffixes))
	if !ok {
		response.Fail("The suffix name of this email is not supported!", nil, ctx)
	}
	//处理base64Captcha
	if global.Server.Website.EnableBase64Captcha {
		if !service.CaptchaSvc.Base64CaptchaStore.Verify(u.Base64Captcha.ID, u.Base64Captcha.B64s, true) {
			response.Fail("Base64Captcha Verification code error,please try again!", nil, ctx) //验证错校验失败会清除store中的value，需要前端重新获取
			return
		}
	}
	//处理邮箱验证码
	if global.Server.Website.EnableEmailCode {
		ok, err = service.UserSvc.VerifyEmailWhenRegister(u)
		if err != nil {
			response.Fail(err.Error(), nil, ctx)
			return
		}
		if !ok {
			response.Fail("Email Verification code error,please try again!", nil, ctx)
			return
		}
	}
	//构建用户信息
	var avatar string
	if u.EmailSuffix == "@qq" {
		avatar = fmt.Sprintf("https://q1.qlogo.cn/g?b=qq&nk=%s&s=100", u.UserName)
	} else {
		avatar = fmt.Sprintf("https://api.multiavatar.com/%s.svg", u.UserName)
	}
	userEmail := u.UserName + u.EmailSuffix //处理邮箱后缀,注册时，用户名和邮箱后缀是分开的
	newUser := &model.User{
		UserName:       userEmail,
		NickName:       userEmail,
		Avatar:         avatar,                                  //头像
		Password:       encrypt_plugin.BcryptEncode(u.Password), //密码
		RoleGroup:      []model.Role{{ID: 2}},                   //默认角色：普通用户角色
		InvitationCode: encrypt_plugin.RandomString(8),          //随机邀请码
	}
	//查找推荐人
	if u.ReferrerCode != "" {
		referrerUser, _ := service.UserSvc.FirstUser(&model.User{InvitationCode: u.ReferrerCode})
		if referrerUser.ID != 0 {
			newUser.ReferrerUserID = referrerUser.ID
		}
	}
	err = service.UserSvc.Register(newUser)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Register error:"+err.Error(), nil, ctx)
		return
	}
	// 推送通知
	if global.Server.Notice.WhenUserRegistered {
		global.GoroutinePool.Submit(func() {
			for k, _ := range global.Server.Notice.AdminIDCache {
				var msg = service.MessageInfo{
					MessageType: service.MESSAGE_TYPE_ADMIN,
					UserID:      k,
					Message: strings.Join([]string{
						"【新注册用户】",
						fmt.Sprintf("时间：%s", time.Now().Format("2006-01-02 15:04:05")),
						fmt.Sprintf("用户名：%s", userEmail),
					}, "\n"),
				}
				service.PushMessageSvc.PushMessage(&msg)
			}
		})
	}
	response.OK("Register success", nil, ctx)
}

// Login
// @Tags [public api] user
// @Summary 用户登录
// @Produce json
// @Param data body model.UserLoginRequest true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/public/user/login [post]
func Login(c *gin.Context) {
	var l model.UserLoginRequest
	err := c.ShouldBind(&l)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, c)
		return
	}
	//查询用户并校验有效性
	user, err := service.UserSvc.Login(&l)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Login error:"+err.Error(), nil, c)
		return
	}
	//签发jwt
	token, err := service.UserSvc.GetUserToken(user)
	response.OK("Login success", gin.H{
		"user":  user,
		"token": token,
	}, c)
}

// ResetUserPassword
// @Tags [public api] user
// @Summary 重置密码
// @Produce json
// @Param data body model.UserLoginRequest true "参数"
// @Success 200 {object} response.ResponseStruct "请求成功；正常：业务代码 code=0；错误：业务代码code=1"
// @Router /api/public/user/resetUserPassword [post]
func ResetUserPassword(ctx *gin.Context) {
	var u model.UserLoginRequest
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//校验邮箱验证码
	ok, err := service.UserSvc.VerifyEmailWhenResetPassword(u)
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	if !ok {
		response.Fail("Email Verification code error,please try again!", nil, ctx)
		return
	}

	err = service.UserSvc.UpdateUser(&model.User{UserName: u.UserName}, map[string]any{"password": encrypt_plugin.BcryptEncode(u.Password)})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ResetUserPassword error:"+err.Error(), nil, ctx)
		return
	}
	// TODO 使该用户token失效
	response.OK("ResetUserPassword success", nil, ctx)

}
