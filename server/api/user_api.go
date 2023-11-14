package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/encrypt_plugin"
	"AirGo/utils/jwt_plugin"
	"AirGo/utils/other_plugin"
	timeTool "AirGo/utils/time_plugin"
	"fmt"
	"net/http"
	"strings"
	"time"

	//"AirGo/utils/encrypt_plugin"

	"AirGo/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	//uuid "github.com/satori/go.uuid"
	uuid "github.com/satori/go.uuid"
)

// 用户注册
func Register(ctx *gin.Context) {
	if !global.Server.Subscribe.EnableRegister {
		response.Fail("Registration closed", nil, ctx)
		return
	}
	var u model.UserRegister
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Register error:"+err.Error(), nil, ctx)
		return
	}
	//判断邮箱后缀
	ok := other_plugin.In(u.EmailSuffix, strings.Fields(global.Server.Subscribe.AcceptableEmailSuffixes))
	if !ok {
		response.Fail("The suffix name of this email is not supported!", nil, ctx)
	}
	//处理base64Captcha
	if !global.Base64CaptchaStore.Verify(u.Base64Captcha.ID, u.Base64Captcha.B64s, true) {
		response.Fail("Verification code error,Please refresh the page and try again!", nil, ctx) //验证错校验失败会清除store中的value，需要前端重新获取
		return
	}
	//校验邮箱验证码
	if global.Server.Subscribe.EnableEmailCode {
		cacheEmail, ok := global.LocalCache.Get(u.UserName + "emailcode")
		if ok {
			if cacheEmail != u.EmailCode {
				response.Fail("Email verification code verification error", nil, ctx)
				return
			}
		} else {
			//cache获取验证码失败,原因：1超时 2系统错误
			response.Fail("Email verification code timeout, please obtain it again", nil, ctx)
			return
		}
	}
	//处理邮箱后缀
	err = service.Register(&model.User{
		UserName:     u.UserName + u.EmailSuffix,
		Password:     u.Password,
		ReferrerCode: u.ReferrerCode,
	})
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Register error:"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(u.UserName + "emailcode")
	response.OK("Register success", nil, ctx)
}

// 用户登录
func Login(c *gin.Context) {
	var l model.UserLogin
	err := c.ShouldBind(&l)
	//key := c.ClientIP()

	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Login error:"+err.Error(), nil, c)
		return
	}
	//校验邮箱验证码
	if global.Server.Subscribe.EnableLoginEmailCode {
		cacheEmail, ok := global.LocalCache.Get(l.UserName + "emailcode")
		global.LocalCache.Delete(l.UserName + "emailcode")
		if ok {
			if cacheEmail != l.EmailCode {
				response.Fail("Email verification code verification error", nil, c)
				return
			}
		} else {
			//cache获取验证码失败,原因：1超时 2系统错误
			response.Fail("Email verification code timeout, please obtain it again", nil, c)
			return
		}
	}
	//查询用户
	user, err := service.Login(&l)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("Login error:"+err.Error(), nil, c)
		return
	}
	//登录以后签发jwt，先查询是否有token缓存
	var token string
	cacheToken, ok := global.LocalCache.Get(l.UserName + "token")
	if ok {
		token = cacheToken.(string)
	} else {
		myCustomClaimsPrefix := jwt_plugin.MyCustomClaimsPrefix{
			UserID:   user.ID,
			UserName: user.UserName,
		}
		ep, _ := timeTool.ParseDuration(global.Server.Security.JWT.ExpiresTime)
		registeredClaims := jwt.RegisteredClaims{
			Issuer:    global.Server.Security.JWT.Issuer,      // 签发者
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)), //过期时间
			NotBefore: jwt.NewNumericDate(time.Now()),         //生效时间
		}
		tokenNew, err := jwt_plugin.GenerateTokenUsingHs256(myCustomClaimsPrefix, registeredClaims, global.Server.Security.JWT.SigningKey)
		if err != nil {
			global.Logrus.Error(err.Error())
			return
		} else {
			global.GoroutinePool.Submit(func() {
				global.LocalCache.Set(l.UserName+"token", tokenNew, ep)
			})
		}
	}
	response.OK("Login success", gin.H{
		"user":  user,
		"token": token,
	}, c)
}

// 修改混淆
func ChangeSubHost(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	var host model.SubHost
	err := ctx.ShouldBind(&host)
	if err != nil || len(host.Host) > 100 {
		global.Logrus.Error(err.Error())
		response.Fail("ChangeSubHost error:"+err.Error(), nil, ctx)
		return
	}
	err = service.ChangeSubHost(uIDInt, host.Host)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("ChangeSubHost error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ChangeSubHost success", nil, ctx)
}

// 获取自身信息
func GetUserInfo(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	u, err := service.GetUserInfo(uIDInt)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetUserInfo error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("GetUserInfo success", u, ctx)

}

// 获取用户列表
func GetUserlist(ctx *gin.Context) {
	var params model.PaginationParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetUserlist error:"+err.Error(), nil, ctx)
	}
	userList, err := service.GetUserlist(&params)
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
		response.Fail("NewUser error:"+err.Error(), nil, ctx)
		return
	}
	err = service.NewUser(u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("NewUser error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("NewUser success", nil, ctx)
}

// 编辑用户信息
func UpdateUser(ctx *gin.Context) {
	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("UpdateUser error:"+err.Error(), nil, ctx)
		return
	}
	//处理角色
	service.DeleteUserRoleGroup(&u)
	var roleArr []string
	for _, v := range u.RoleGroup {
		roleArr = append(roleArr, v.RoleName)
	}
	roles, _ := service.FindRoleIdsByRoleNameArr(roleArr)
	u.RoleGroup = roles

	err = service.SaveUser(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("UpdateUser error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("UpdateUser success", nil, ctx)
}

// 删除用户
func DeleteUser(ctx *gin.Context) {
	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("DeleteUser error:"+err.Error(), err.Error(), ctx)
		return
	}
	//删除用户关联的角色
	service.DeleteUserRoleGroup(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("DeleteUser error:"+err.Error(), nil, ctx)
		return
	}
	// 删除用户
	err = service.DeleteUser(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("DeleteUser error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("DeleteUser success", nil, ctx)

}

// 修改密码
func ChangeUserPassword(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	var u model.UserChangePassword
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ChangeUserPassword error:"+err.Error(), nil, ctx)
		return
	}
	var user = model.User{
		ID:       uIDInt,
		Password: encrypt_plugin.BcryptEncode(u.Password),
	}
	//
	err = service.UpdateUser(&user)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ChangeUserPassword error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ChangeUserPassword success", nil, ctx)
}

// 重置密码
func ResetUserPassword(ctx *gin.Context) {
	var u model.UserLogin
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ResetUserPassword error:"+err.Error(), nil, ctx)
		return
	}
	//判断邮箱后缀
	ok := other_plugin.In(u.UserName[strings.Index(u.UserName, "@"):], strings.Fields(global.Server.Subscribe.AcceptableEmailSuffixes))
	if !ok {
		response.Fail("The suffix name of this email is not supported!", nil, ctx)
	}

	//校验邮箱验证码
	emailcode, _ := global.LocalCache.Get(u.UserName + "emailcode")
	if emailcode != u.EmailCode {
		response.Fail("Email verification code error", nil, ctx)
		return
	}
	var user = model.User{
		UserName: u.UserName,
		Password: encrypt_plugin.BcryptEncode(u.Password),
	}
	err = service.ResetUserPassword(&user)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ResetUserPassword error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ResetUserPassword success", nil, ctx)

}

// 获取订阅
func GetSub(ctx *gin.Context) {
	//订阅参数
	link := ctx.Query("link")
	subType := ctx.Query("type")

	res := service.GetUserSub(link, subType)
	if res == "" {
		return
	}
	ctx.String(http.StatusOK, res)

}

// 重置订阅
func ResetSub(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	var u = model.User{
		ID:            uIDInt,
		UUID:          uuid.NewV4(),
		SubscribeInfo: model.SubscribeInfo{SubscribeUrl: encrypt_plugin.RandomString(8)}, //随机字符串订阅url
	}
	err := service.UpdateUser(&u)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ResetSub error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("ResetSub success", nil, ctx)
}

// 打卡
func ClockIn(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("user id error", nil, ctx)
		return
	}
	//判断是否已打卡
	_, ok = global.LocalCache.Get(fmt.Sprintf("%d", uIDInt) + "clockin")
	if ok { //已打卡
		response.Fail("You have already checked in ", nil, ctx)
		return
	}

	msg, err := service.ClockIn(uIDInt)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("ClockIn error:"+err.Error(), nil, ctx)
		return
	}
	now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	global.LocalCache.Set(fmt.Sprintf("%d", uIDInt)+"clockin", nil, zeroTime.Sub(now))
	response.OK("ClockIn success", msg, ctx)
}
