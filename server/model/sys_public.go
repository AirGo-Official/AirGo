package model

// base64验证码
type Base64CaptchaInfo struct {
	ID   string `json:"id"`
	B64s string `json:"b64s"` //响应时存base64数据，请求时存前端看到的验证码。响应，请求共用该结构体
}
