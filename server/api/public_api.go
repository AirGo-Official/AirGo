package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/mail_plugin"
	"github.com/ppoonk/AirGo/utils/other_plugin"
	"github.com/ppoonk/AirGo/utils/response"
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
	//判断邮箱后缀
	ok := other_plugin.In(u.UserName[strings.Index(u.UserName, "@"):], strings.Fields(global.Server.Subscribe.AcceptableEmailSuffixes))
	if !ok {
		response.Fail("The suffix name of this email is not supported!", nil, ctx)
	}

	_, ok = global.LocalCache.Get(u.UserName + "emailcode")
	if ok {
		response.Fail("The email verification code is frequently obtained. Please try again in 60 minutes!", nil, ctx)
		return
	}
	//生成验证码
	randomStr := encrypt_plugin.RandomString(4) //4位随机数
	var wg sync.WaitGroup
	wg.Add(2)
	//验证码存入local cache
	global.GoroutinePool.Submit(func() {
		global.LocalCache.Set(u.UserName+"emailcode", randomStr, 60*time.Second) //过期
		wg.Done()
	})
	//发送邮件
	global.GoroutinePool.Submit(func() {
		//判断别名邮箱
		from := global.Server.Email.EmailFrom
		if global.Server.Email.EmailFromAlias != "" {
			from = global.Server.Email.EmailFromAlias
		}
		//选择内容模板
		originalText := strings.Replace(global.Server.Email.EmailContent, "emailcode", randomStr, -1)
		err = mail_plugin.SendEmail(global.EmailDialer, from, global.Server.Email.EmailNickname, u.UserName, global.Server.Email.EmailSubject, originalText)
		if err != nil {
			global.Logrus.Error(err.Error())
			//"The email verification code has failed to be sent.error:gomail: could not send email 1: gomail: invalid address "=?UTF-8?q?=E5=90=8A=E7=82=B8=E5=A4=A9=E6=9C=BA=E5=9C=BA=E7=AE=A1=E7=90=86?= =?UTF-8?q?=E5=91=98<poonk@foxmail.com>?=": mail: expected single address, got "?=""
			response.Fail("The email verification code has failed to be sent. Error:"+err.Error(), nil, ctx)

		} else {
			response.OK("Email verification code has been sent, please check your email carefully.", nil, ctx)
		}
		wg.Done()
	})
	wg.Wait()
}
