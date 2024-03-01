package public_api

import (
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/other_plugin"
	"github.com/ppoonk/AirGo/utils/response"
	"net/http"
	"strings"
	"time"
)

// GetBase64Captcha
// @Tags public_api
// @Summary 发送base64验证码
// @Description  发送base64验证码
// @Produce json
// @Router 	/api/public/getBase64Captcha [get]
// @Success 200 {object} model.Base64CaptchaInfo
func GetBase64Captcha(ctx *gin.Context) {
	id, b64s, _, err := global.Base64Captcha.Generate()
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

func GetEmailCode(ctx *gin.Context) {
	var e model.EmailRequest
	err := ctx.ShouldBind(&e)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail(constant.ERROR_REQUEST_PARAMETER_PARSING_ERROR+err.Error(), nil, ctx)
		return
	}
	// 判断邮箱后缀是否合法
	ok := other_plugin.In(e.TargetEmail[strings.Index(e.TargetEmail, "@"):], strings.Fields(global.Server.Website.AcceptableEmailSuffixes))
	if !ok {
		response.Fail("The suffix name of this email is not supported!", nil, ctx)
		return
	}
	// 根据邮件类型，进行处理
	switch e.EmailType {
	case constant.EMAIL_TYPE_USER_REGISTER:
		SendEmailCode(ctx, &e, constant.CACHE_USER_REGISTER_EMAIL_CODE_BY_USERNAME)
	case constant.EMAIL_TYPE_USER_RESETPWD:
		SendEmailCode(ctx, &e, constant.CACHE_USER_RESET_PWD_EMAIL_CODE_BY_USERNAME)
	case constant.EMAIL_TYPE_TEST:
		SendEmailCode(ctx, &e, "")
	default:
		response.Fail("Illegal email type", nil, ctx)
		return
	}
}

// 发送邮件验证码
func SendEmailCode(ctx *gin.Context, e *model.EmailRequest, keyPre string) {
	var randomStr string
	// 查缓存，如有，则发送，不用生成新的验证码
	cache, ok := global.LocalCache.Get(keyPre + e.TargetEmail)
	if ok {
		randomStr = cache.(string)
	} else {
		//生成验证码
		randomStr = encrypt_plugin.RandomString(4) //4位随机数
	}
	// 验证码默认3分钟缓存时间
	// 前端在1分钟后，显示可以重新获取
	global.LocalCache.Set(keyPre+e.TargetEmail, randomStr, 3*time.Minute)
	//判断别名邮箱
	from := global.Server.Email.EmailFrom
	if global.Server.Email.EmailFromAlias != "" {
		from = global.Server.Email.EmailFromAlias
	}
	//内容
	originalText := strings.Replace(global.Server.Email.EmailContent, "emailcode", randomStr, -1)
	emailMsg := admin_logic.EmailMsg{
		From:      from,
		To:        e.TargetEmail,
		NickName:  global.Server.Email.EmailNickname,
		Subject:   global.Server.Email.EmailSubject,
		EmailText: originalText,
	}
	// 入队:邮件验证码发送队列
	global.Queue.Publish(constant.EMAIL_CODE, emailMsg)
	response.OK("Email code has been sent.", nil, ctx)
	return

}

// 获取订阅
func GetSub(ctx *gin.Context) {
	//Shadowrocket/2070 CFNetwork/1325.0.1 Darwin/21.1.0
	//ClashMetaForAndroid/2.8.9.Meta
	//ClashX/1.118.0 (com.west2online.ClashX; build:1.118.0; macOS 10.15.7) Alamofire/5.8.0
	//Quantumult/627 CFNetwork/1325.0.1 Darwin/21.1.0
	//NekoBox/Android/1.2.9 (Prefer ClashMeta Format)
	//v2rayNG/1.8.9
	//V2rayU/4.0.0 CFNetwork/1128.0.1 Darwin/19.6.0 (x86_64)
	//v2rayN/6.30

	clientType := ""
	ua := ctx.Request.Header.Get("User-Agent")
	if strings.HasPrefix(ua, "NekoBox") {
		clientType = "NekoBox"
		goto next
	}
	if strings.HasPrefix(ua, "v2rayNG") {
		clientType = "v2rayNG"
		goto next
	}
	if strings.HasPrefix(ua, "v2rayN") {
		clientType = "v2rayN"
		goto next
	}
	if strings.HasPrefix(ua, "Clash") {
		clientType = "Clash"
		goto next
	}
	if strings.HasPrefix(ua, "Shadowrocket") {
		clientType = "Shadowrocket"
		goto next
	}
	if strings.HasPrefix(ua, "Surge") {
		clientType = "Surge"
		goto next
	}
	if strings.HasPrefix(ua, "Quantumult") {
		clientType = "Quantumult"
		goto next
	}
	if strings.HasPrefix(ua, "V2rayU") {
		clientType = "V2rayU"
		goto next
	}
	if clientType == "" {
		clientType = ctx.Query("type")
	}
	if clientType == "" {
		clientType = "NekoBox"
	}
next:
	id := ctx.Param("id")
	//fmt.Println("id:", id)
	res, header := customerService.GetSubscribe(id, clientType)
	if res == "" {
		return
	}
	ctx.Header("subscription-userinfo", header)
	ctx.String(http.StatusOK, res)
}
