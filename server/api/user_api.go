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
	if !global.Server.System.EnableRegister {
		response.Fail("已关闭注册", nil, ctx)
		return
	}
	var u model.UserRegister
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error("注册参数错误:", err.Error())
		response.Fail("注册参数错误"+err.Error(), nil, ctx)
		return
	}
	//fmt.Println("注册", u)
	//处理base64Captcha
	if !global.Base64CaptchaStore.Verify(u.Base64Captcha.ID, u.Base64Captcha.B64s, true) {
		response.Fail("验证码错误", nil, ctx) //验证错误会清除store中的value，需要前端重新获取
		return
	}
	//校验邮箱验证码
	if global.Server.System.EnableEmailCode {
		cacheEmail, ok := global.LocalCache.Get(u.UserName + "emailcode")
		if ok {
			if cacheEmail != u.EmailCode {
				response.Fail("邮箱验证码校验错误", nil, ctx)
				return
			}
		} else {
			//cache获取验证码失败,原因：1超时 2系统错误
			response.Fail("邮箱验证码超时，请重新获取", nil, ctx)
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
		global.Logrus.Error("注册错误:", err.Error())
		response.Fail("注册错误"+err.Error(), nil, ctx)
		return
	}
	global.LocalCache.Delete(u.UserName + "emailcode")
	response.OK("注册成功", nil, ctx)
}

// 用户登录
func Login(c *gin.Context) {
	var l model.UserLogin
	err := c.ShouldBind(&l)
	//key := c.ClientIP()

	if err != nil {
		global.Logrus.Error("用户登录参数错误:", err.Error())
		response.Fail("用户登录参数错误"+err.Error(), nil, c)
		return
	}
	//校验邮箱验证码
	if global.Server.System.EnableLoginEmailCode {
		cacheEmail, ok := global.LocalCache.Get(l.UserName + "emailcode")
		global.LocalCache.Delete(l.UserName + "emailcode")
		if ok {
			if cacheEmail != l.EmailCode {
				response.Fail("邮箱验证码校验错误", nil, c)
				return
			}
		} else {
			//cache获取验证码失败,原因：1超时 2系统错误
			response.Fail("邮箱验证码超时，请重新获取", nil, c)
			return
		}
	}
	//查询用户
	user, err := service.Login(&l)
	if err != nil {
		response.Fail("查询用户"+err.Error(), nil, c)
		global.Logrus.Error("查询用户", err.Error())
		return
	}
	//登录以后签发jwt，先查询是否有token缓存
	var token string
	cacheToken, ok := global.LocalCache.Get(l.UserName + "token")
	if ok {
		fmt.Println("cacheToken:", cacheToken)
		token = cacheToken.(string)
	} else {
		myCustomClaimsPrefix := jwt_plugin.MyCustomClaimsPrefix{
			UserID:   user.ID,
			UserName: user.UserName,
		}
		ep, _ := timeTool.ParseDuration(global.Server.JWT.ExpiresTime)
		registeredClaims := jwt.RegisteredClaims{
			Issuer:    global.Server.JWT.Issuer,               // 签发者
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)), //过期时间
			NotBefore: jwt.NewNumericDate(time.Now()),         //生效时间
		}
		tokenNew, err := jwt_plugin.GenerateTokenUsingHs256(myCustomClaimsPrefix, registeredClaims, global.Server.JWT.SigningKey)
		if err != nil {
			global.Logrus.Error("生成token err", err.Error())
			return
		} else {
			token = tokenNew
			//fmt.Println("tokenNew:", tokenNew)
			go func(l *model.UserLogin, token string) {
				global.LocalCache.Set(l.UserName+"token", token, ep)
			}(&l, token)
		}
	}
	//fmt.Println("生成token :", token)
	response.OK("登录成功", gin.H{
		"user":  user,
		"token": token,
	}, c)
}

// 修改混淆
func ChangeSubHost(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}
	var host model.SubHost
	err := ctx.ShouldBind(&host)
	if err != nil || len(host.Host) > 100 {
		global.Logrus.Error("修改混淆参数错误:", err.Error())
		response.Fail("修改混淆参数错误"+err.Error(), nil, ctx)
		return
	}
	err = service.ChangeSubHost(uIDInt, host.Host)
	if err != nil {
		global.Logrus.Error("修改混淆错误:", err.Error())
		response.Fail("修改混淆错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("修改混淆成功", nil, ctx)
}

// 获取自身信息
func GetUserInfo(ctx *gin.Context) {
	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}

	u, err := service.GetUserInfo(uIDInt)
	if err != nil {
		global.Logrus.Error("获取自身信息:", err.Error())
		response.Fail("获取自身信息"+err.Error(), nil, ctx)
		return
	}
	response.OK("获取信息成功", u, ctx)

}

// 获取用户列表
func GetUserlist(ctx *gin.Context) {
	var params model.PaginationParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		global.Logrus.Error("获取用户列表参数错误:", err.Error())
		response.Fail("获取用户列表参数错误"+err.Error(), nil, ctx)
	}
	userList, err := service.GetUserlist(&params)
	if err != nil {
		global.Logrus.Error("获取用户列表错误:", err.Error())
		response.Fail("获取用户列表错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("获取用户列表成功", userList, ctx)
}

// 新建用户
func NewUser(ctx *gin.Context) {
	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error("新建用户参数错误:", err.Error())
		response.Fail("新建用户参数错误"+err.Error(), nil, ctx)
		return
	}
	//fmt.Println("新建用户:", u.RoleGroup)
	err = service.NewUser(u)
	if err != nil {
		global.Logrus.Error("新建用户错误:", err.Error())
		response.Fail("新建用户错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("新建用户成功", nil, ctx)
}

// 编辑用户信息
func UpdateUser(ctx *gin.Context) {
	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error("修改用户参数错误:", err.Error())
		response.Fail("修改用户参数错误"+err.Error(), nil, ctx)
		return
	}
	fmt.Println("编辑用户信息:", u.SubscribeInfo.ExpiredAt)
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
		global.Logrus.Error("修改用户错误 error:", err)
		response.Fail("修改用户错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("修改用户成功", nil, ctx)
}

// 删除用户
func DeleteUser(ctx *gin.Context) {
	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error("删除用户参数错误 error:", err)
		response.Fail("删除用户参数错误"+err.Error(), err.Error(), ctx)
		return
	}
	//删除用户关联的角色
	service.DeleteUserRoleGroup(&u)
	if err != nil {
		global.Logrus.Error("删除用户角色错误 error:", err)
		response.Fail("删除用户角色错误"+err.Error(), nil, ctx)
		return
	}
	// 删除用户
	err = service.DeleteUser(&u)
	if err != nil {
		global.Logrus.Error("删除用户错误 error:", err)
		response.Fail("删除用户错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("删除用户成功", nil, ctx)

}

// 修改密码
func ChangeUserPassword(ctx *gin.Context) {

	uIDInt, ok := other_plugin.GetUserIDFromGinContext(ctx)
	if !ok {
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}
	var u model.UserChangePassword
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error("修改密码参数错误 error:", err)
		response.Fail("修改密码参数错误"+err.Error(), nil, ctx)
		return
	}
	//
	var user = model.User{
		ID:       uIDInt,
		Password: encrypt_plugin.BcryptEncode(u.Password),
	}
	//
	err = service.UpdateUser(&user)
	if err != nil {
		global.Logrus.Error("修改密码错误 error:", err)
		response.Fail("修改密码错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("修改密码成功", nil, ctx)
}

// 重置密码
func ResetUserPassword(ctx *gin.Context) {
	var u model.UserLogin
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error("重置密码参数错误 error:", err)
		response.Fail("重置密码参数错误"+err.Error(), nil, ctx)
		return
	}
	//校验邮箱验证码
	emailcode, _ := global.LocalCache.Get(u.UserName + "emailcode")
	if emailcode != u.EmailCode {
		response.Fail("邮箱验证码错误", nil, ctx)
		return
	}
	var user = model.User{
		UserName: u.UserName,
		Password: encrypt_plugin.BcryptEncode(u.Password),
	}
	//
	err = service.ResetUserPassword(&user)
	if err != nil {
		global.Logrus.Error("重置密码错误 error:", err)
		response.Fail("重置密码错误"+err.Error(), nil, ctx)
		return
	}
	//global.LocalCache.Delete(u.UserName + "emailcode")
	response.OK("重置密码成功", nil, ctx)

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
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}
	var u = model.User{
		ID:            uIDInt,
		UUID:          uuid.NewV4(),
		SubscribeInfo: model.SubscribeInfo{SubscribeUrl: encrypt_plugin.RandomString(8)}, //随机字符串订阅url
	}
	err := service.UpdateUser(&u)
	if err != nil {
		global.Logrus.Error("重置订阅错误 error:", err)
		response.Fail("重置订阅错误"+err.Error(), nil, ctx)
		return
	}
	response.OK("重置订阅成功", nil, ctx)
}
