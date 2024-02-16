package model

import "time"

// Server 全局配置
type Server struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id"   gorm:"primaryKey"`

	Website  Website  `json:"website"  gorm:"embedded"`
	Email    Email    `json:"email"    gorm:"embedded"`
	Security Security `json:"security" gorm:"embedded"`
	Notice   Notice   `json:"notice"   gorm:"embedded"`
}
type Notice struct {
	BotToken           string `json:"bot_token"            gorm:"comment:tg bot token"`
	TGAdmin            string `json:"tg_admin"             gorm:"comment:tg admin"`
	TGSocks5           string `json:"tg_socks5"            gorm:"comment:tg socks5代理"`
	WhenUserRegistered bool   `json:"when_user_registered" gorm:"comment:用户注册后通知"`
	WhenUserPurchased  bool   `json:"when_user_purchased"  gorm:"comment:用户购买成功后通知"`
	WhenNodeOffline    bool   `json:"when_node_offline"    gorm:"comment:节点离线时通知"`
}

type Email struct {
	EmailFrom      string `json:"email_from"       gorm:"comment:发件人"`
	EmailFromAlias string `json:"email_from_alias" gorm:"comment:发件人别名"`
	EmailSecret    string `json:"email_secret"   gorm:"comment:密钥"`
	EmailHost      string `json:"email_host"     gorm:"comment:服务器地址"`
	EmailPort      int64  `json:"email_port"     gorm:"comment:端口"`
	EmailIsSSL     bool   `json:"email_is_ssl"   gorm:"comment:是否SSL"`
	EmailNickname  string `json:"email_nickname" gorm:"comment:昵称"`
	EmailSubject   string `json:"email_subject"  gorm:"comment:邮件主题;default:hello!"`
	EmailContent   string `json:"email_content"  gorm:"comment:邮件内容;type:text"`
}
type Security struct {
	Captcha         Captcha         `json:"captcha" gorm:"embedded"`
	JWT             JWT             `json:"jwt"     gorm:"embedded"`
	RateLimitParams RateLimitParams `json:"rate_limit_params"gorm:"embedded"`
}

type Captcha struct {
	KeyLong            int64 `json:"key_long"        gorm:"default:6;comment:验证码长度"`
	ImgWidth           int64 `json:"img_width"       gorm:"default:240;comment:验证码宽度"`
	ImgHeight          int64 `json:"img_height"      gorm:"default:80;comment:验证码高度"`
	OpenCaptcha        int64 `json:"open_captcha"    gorm:"default:2"`
	OpenCaptchaTimeOut int64 `json:"open_captcha_time_out" gorm:"default:300;comment:验证码超时时间，单位：s(秒)"`
}
type JWT struct {
	SigningKey  string `json:"signing_key"  gorm:"default:AirGo;comment:jwt签名"`
	ExpiresTime string `json:"expires_time" gorm:"default:30d;comment:过期时间"`
	BufferTime  string `json:"buffer_time"  gorm:"default:1d;comment:缓冲时间"`
	Issuer      string `json:"issuer"       gorm:"default:AirGo;comment:签发者"`
}

// RateLimitParams 限流参数
type RateLimitParams struct {
	IPRoleParam int64 `json:"ip_role_param" gorm:"default:600"`
	VisitParam  int64 `json:"visit_param"   gorm:"default:600"`
}
type Website struct {
	EnableRegister          bool   `json:"enable_register"           gorm:"default:true;comment:是否开启注册"`
	AcceptableEmailSuffixes string `json:"acceptable_email_suffixes" gorm:"comment:可接受的邮箱后缀"`
	EnableEmailCode         bool   `json:"enable_email_code"         gorm:"default:false;comment:是否开启注册email 验证码"`
	EnableLoginEmailCode    bool   `json:"enable_login_email_code"   gorm:"default:false;comment:是否开启登录email 验证码"`
	IsMultipoint            bool   `json:"is_multipoint"     gorm:"default:true;comment:是否多点登录"`
	BackendUrl              string `json:"backend_url"       gorm:"comment:后端地址"`
	FrontendUrl             string `json:"frontend_url"      gorm:"comment:官网地址"`
	SubName                 string `json:"sub_name"          gorm:"default:AirGo;comment:订阅名称"`
	TEK                     string `json:"tek"               gorm:"default:airgo;comment:前后端通信密钥"`
	EnabledClockIn          bool   `json:"enabled_clock_in"  gorm:"default:true;comment:是否开启打卡"`
}

// 公共配置参数
type PublicSystem struct {
	EnableRegister          bool    `json:"enable_register"`           // 是否开启注册
	AcceptableEmailSuffixes string  `json:"acceptable_email_suffixes"` // 可接受的邮箱后缀
	EnableEmailCode         bool    `json:"enable_email_code"`         // 是否开启注册email 验证码
	EnableLoginEmailCode    bool    `json:"enable_login_email_code"`   // 是否开启登录email 验证码
	RebateRate              float64 `json:"rebate_rate"`               // 佣金率
	BackendUrl              string  `json:"backend_url"`               // 后端地址
	EnabledClockIn          bool    `json:"enabled_clock_in"`          // 是否开启打卡
}
