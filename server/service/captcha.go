package service

import (
	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	Base64Captcha      *base64Captcha.Captcha
	Base64CaptchaStore base64Captcha.Store
}

var CaptchaSvc *Captcha

func newCaptchaSvc() *Captcha {
	return &Captcha{
		Base64Captcha:      nil,
		Base64CaptchaStore: nil,
	}
}
func InitBase64Captcha() {
	CaptchaSvc = newCaptchaSvc()
	// base64Captcha.DefaultMemStore 是默认的过期时间10分钟。也可以自己设定参数 base64Captcha.NewMemoryStore(GCLimitNumber, Expiration)
	CaptchaSvc.Base64CaptchaStore = base64Captcha.DefaultMemStore
	driver := base64Captcha.NewDriverDigit(38, 120, 4, 0.2, 10)
	CaptchaSvc.Base64Captcha = base64Captcha.NewCaptcha(driver, CaptchaSvc.Base64CaptchaStore)
}
