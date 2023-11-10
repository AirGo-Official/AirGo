package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/utils/encrypt_plugin"
	"AirGo/utils/mail_plugin"
	"AirGo/utils/response"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"time"
)

// 发送base64验证码
func GetBase64Captcha(ctx *gin.Context) {

	id, b64s, err := global.Base64Captcha.Generate()
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("SendBase64Captcha"+err.Error(), nil, ctx)
		return
	}
	var b64Captcha model.Base64CaptchaInfo
	b64Captcha.ID = id
	b64Captcha.B64s = b64s
	response.OK("GetBase64Captcha success", b64Captcha, ctx)

}

// 验证base64验证码
func VerifyBase64Captcha(ctx *gin.Context) {
	var b64Captcha model.Base64CaptchaInfo
	err := ctx.ShouldBind(&b64Captcha)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("VerifyBase64Captcha error:"+err.Error(), nil, ctx)
		return
	}
	if !global.Base64CaptchaStore.Verify(b64Captcha.ID, b64Captcha.B64s, true) {
		response.Fail("VerifyBase64Captcha error:"+err.Error(), nil, ctx)
		return
	}
	response.OK("VerifyBase64Captcha success", nil, ctx)
}

// 邮箱验证码
func GetMailCode(ctx *gin.Context) {
	var u model.UserRegisterEmail
	err := ctx.ShouldBind(&u)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetMailCode error:"+err.Error(), nil, ctx)
		return
	}
	_, ok := global.LocalCache.Get(u.UserName + "emailcode")
	if ok {
		response.Fail("The email verification code is frequently obtained. Please try again in 60 minutes!", nil, ctx)
		return
	}
	//生成验证码
	randomStr := encrypt_plugin.RandomString(4) //4位随机数
	var wg sync.WaitGroup
	wg.Add(3)
	//验证码存入local cache
	go func(wg *sync.WaitGroup) {
		global.LocalCache.Set(u.UserName+"emailcode", randomStr, 60*time.Second) //过期
		wg.Done()
	}(&wg)
	//发送邮件
	go func(wg *sync.WaitGroup) {
		//判断别名邮箱
		from := global.Server.Email.EmailFrom
		if global.Server.Email.EmailFromAlias != "" {
			from = global.Server.Email.EmailFromAlias
		}
		//选择验证码模
		originalText := strings.Replace(global.Server.Email.EmailContent, "emailcode", randomStr, -1)
		err = mail_plugin.SendEmail(global.EmailDialer, from, global.Server.Email.EmailNickname, u.UserName, global.Server.Email.EmailSubject, originalText)
		if err != nil {
			global.Logrus.Error(err.Error())
		}
		wg.Done()
	}(&wg)
	go func(ctx *gin.Context, wg *sync.WaitGroup) {
		response.OK("Email verification code has been sent, please check your email carefully.", nil, ctx)
		wg.Done()
	}(ctx, &wg)
	wg.Wait()
}
