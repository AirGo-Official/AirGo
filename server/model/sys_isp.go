package model

import "time"

type ISP struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt *time.Time `json:"-" gorm:"index"`
	ID int64 `json:"id"   gorm:"primary_key"`

	UserID        int64         `json:"user_id"`
	Mobile        string        `json:"mobile"`
	ISPType       string        `json:"isp_type"      gorm:"conment:运营商类型，unicom,telecom"`
	Status        bool          `json:"status"        gorm:"default:false"`
	UnicomConfig  UnicomConfig  `json:"unicom_config" gorm:"embedded"`
	TelecomConfig TelecomConfig `json:"telecom_config" gorm:"embedded"`
}

type UnicomConfig struct {
	Version      string `json:"version" gorm:"default:iphone_c@10.5"`
	APPID        string `json:"app_id"`
	Cookie       string `json:"cookie"`
	UnicomMobile string `json:"unicommobile"`
	Password     string `json:"password"` //短信验证码
}
type TelecomConfig struct {
	PhoneNum                   string `json:"phoneNum"`                   //加密手机号
	TelecomPassword            string `json:"telecomPassword"`            //登录密码
	Timestamp                  string `json:"timestamp"`                  //时间戳
	LoginAuthCipherAsymmertric string `json:"loginAuthCipherAsymmertric"` //加密字段
	DeviceUid                  string `json:"deviceUid"`                  //设备id
	TelecomToken               string `json:"telecomToken"`               //登录后token
}
