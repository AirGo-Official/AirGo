package public_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/other_plugin"
	"github.com/ppoonk/AirGo/utils/response"
	"strings"
)

// 用户注册
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
	//判断邮箱后缀
	ok := other_plugin.In(u.EmailSuffix, strings.Fields(global.Server.Website.AcceptableEmailSuffixes))
	if !ok {
		response.Fail("The suffix name of this email is not supported!", nil, ctx)
	}
	//处理base64Captcha
	if !global.Base64CaptchaStore.Verify(u.Base64Captcha.ID, u.Base64Captcha.B64s, true) {
		response.Fail("Verification code error,please try again!", nil, ctx) //验证错校验失败会清除store中的value，需要前端重新获取
		return
	}
	//处理邮箱验证码
	userEmail := u.UserName + u.EmailSuffix //处理邮箱后缀
	if global.Server.Website.EnableEmailCode {
		cacheEmail, ok := global.LocalCache.Get(constant.CACHE_USER_REGISTER_EMAIL_CODE_BY_USERNAME + userEmail)
		if ok {
			if !strings.EqualFold(cacheEmail.(string), u.EmailCode) {
				response.Fail("Email verification error", nil, ctx)
				return
			}
		} else {
			//cache获取验证码失败,原因：1超时 2系统错误
			response.Fail("Timeout, please try again", nil, ctx)
			return
		}
	}
	global.LocalCache.Delete(constant.CACHE_USER_REGISTER_EMAIL_CODE_BY_USERNAME + userEmail)

	//构建用户信息
	var avatar string
	if u.EmailSuffix == "@qq" {
		avatar = fmt.Sprintf("https://q1.qlogo.cn/g?b=qq&nk=%s&s=100", u.UserName)
	} else {
		avatar = fmt.Sprintf("https://api.multiavatar.com/%s.svg", u.UserName)
	}
	err = userService.Register(&model.User{
		UserName:       u.UserName,
		NickName:       u.UserName,
		Avatar:         avatar,                                  //头像
		Password:       encrypt_plugin.BcryptEncode(u.Password), //密码
		RoleGroup:      []model.Role{{ID: 2}},                   //默认角色：普通用户角色
		InvitationCode: encrypt_plugin.RandomString(8),          //邀请码
		ReferrerCode:   u.ReferrerCode,                          //推荐人
	})

	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Register error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(userEmail + "emailcode")
	response.OK("Register success", nil, ctx)
}

// 用户登录
func Login(c *gin.Context) {
	var l model.UserLoginRequest
	err := c.ShouldBind(&l)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, c)
		return
	}
	//查询用户并校验有效性
	user, err := userService.Login(&l)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Login error:"+err.Error(), nil, c)
		return
	}
	//签发jwt
	token, err := userService.GetUserToken(user)
	response.OK("Login success", gin.H{
		"user":  user,
		"token": token,
	}, c)
}

// 重置密码
func ResetUserPassword(ctx *gin.Context) {
	var u model.UserLoginRequest
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	//校验邮箱验证码
	cacheEmail, ok := global.LocalCache.Get(constant.CACHE_USER_RESET_PWD_EMAIL_CODE_BY_USERNAME + u.UserName)
	if ok {
		if !strings.EqualFold(cacheEmail.(string), u.EmailCode) {
			response.Fail("Email verification error", nil, ctx)
			return
		}
	} else {
		//cache获取验证码失败,原因：1超时 2系统错误
		response.Fail("Timeout, please try again", nil, ctx)
		return
	}
	global.LocalCache.Delete(constant.CACHE_USER_RESET_PWD_EMAIL_CODE_BY_USERNAME + u.UserName)
	err = userService.UpdateUser(&model.User{UserName: u.UserName}, map[string]any{"password": encrypt_plugin.BcryptEncode(u.Password)})
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ResetUserPassword error:"+err.Error(), nil, ctx)
		return
	}
	// todo 删除该用户token cache
	response.OK("ResetUserPassword success", nil, ctx)

}
