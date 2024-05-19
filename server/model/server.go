package model

import "time"

// Server 全局配置
type Server struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id"   gorm:"primaryKey"`

	Website   Website   `json:"website"  gorm:"embedded"`
	Email     Email     `json:"email"    gorm:"embedded"`
	Security  Security  `json:"security" gorm:"embedded"`
	Notice    Notice    `json:"notice"   gorm:"embedded"`
	Subscribe Subscribe `json:"subscribe" gorm:"embedded"`
	Finance   Finance   `json:"finance"   gorm:"embedded"`
}
type Notice struct {
	EnableTGBot          bool               `json:"enable_tg_bot"`
	EnableEmail          bool               `json:"enable_email"`
	EnableWebMail        bool               `json:"enable_web_mail"`
	AdminID              string             `json:"admin_id"`
	AdminIDCache         map[int64]struct{} `json:"-" gorm:"-"` //用来存储通知消息时的管理员id，key=user_id
	AdminIDCacheWithTGID map[int64]struct{} `json:"-" gorm:"-"` //用来存储通知消息时的管理员id，key=tg_id
	BotToken             string             `json:"bot_token"            gorm:"comment:tg bot token"`
	TGSocks5             string             `json:"tg_socks5"            gorm:"comment:tg socks5代理"`
	WhenUserRegistered   bool               `json:"when_user_registered" gorm:"comment:用户注册后通知"`
	WhenUserPurchased    bool               `json:"when_user_purchased"  gorm:"comment:用户购买成功后通知"`
	WhenNodeOffline      bool               `json:"when_node_offline"    gorm:"comment:节点离线时通知"`
	WhenNewTicket        bool               `json:"when_new_ticket"      gorm:"comment:新工单时通知"`
}

type Email struct {
	EmailFrom      string `json:"email_from"       gorm:"comment:发件人"`
	EmailFromAlias string `json:"email_from_alias" gorm:"comment:发件人别名"`
	EmailSecret    string `json:"email_secret"   gorm:"comment:密钥"`
	EmailHost      string `json:"email_host"     gorm:"comment:服务器地址"`
	EmailPort      int64  `json:"email_port"     gorm:"comment:端口"`
	EmailIsSSL     bool   `json:"email_is_ssl"   gorm:"comment:是否SSL"`
	EmailNickname  string `json:"email_nickname" gorm:"comment:昵称"`
	EmailSubject   string `json:"email_subject"  gorm:"comment:邮件主题"`
	EmailContent   string `json:"email_content"  gorm:"comment:邮件内容;type:text"`
}
type Security struct {
	Captcha         Captcha         `json:"captcha" gorm:"embedded"`
	JWT             JWT             `json:"jwt"     gorm:"embedded"`
	RateLimitParams RateLimitParams `json:"rate_limit_params"gorm:"embedded"`
}

type Captcha struct {
	KeyLong            int64 `json:"key_long"        gorm:"comment:验证码长度"`
	ImgWidth           int64 `json:"img_width"       gorm:"comment:验证码宽度"`
	ImgHeight          int64 `json:"img_height"      gorm:"comment:验证码高度"`
	OpenCaptcha        int64 `json:"open_captcha"`
	OpenCaptchaTimeOut int64 `json:"open_captcha_time_out" gorm:"comment:验证码超时时间，单位：s(秒)"`
}
type JWT struct {
	SigningKey  string `json:"signing_key"  gorm:"comment:jwt签名"`
	ExpiresTime string `json:"expires_time" gorm:"comment:过期时间"`
	BufferTime  string `json:"buffer_time"  gorm:"comment:缓冲时间"`
	Issuer      string `json:"issuer"       gorm:"comment:签发者"`
}

// RateLimitParams 限流参数
type RateLimitParams struct {
	IPRoleParam int64 `json:"ip_role_param"`
	VisitParam  int64 `json:"visit_param"`
}
type Website struct {
	EnableRegister          bool   `json:"enable_register"           gorm:"comment:是否开启注册"`
	AcceptableEmailSuffixes string `json:"acceptable_email_suffixes" gorm:"comment:可接受的邮箱后缀"`
	EnableBase64Captcha     bool   `json:"enable_base64_captcha"     gorm:"comment:是否开启注册图片验证码"`
	EnableEmailCode         bool   `json:"enable_email_code"         gorm:"comment:是否开启注册email 验证码"`
	EnableLoginEmailCode    bool   `json:"enable_login_email_code"   gorm:"comment:是否开启登录email 验证码"`
	IsMultipoint            bool   `json:"is_multipoint"      gorm:"comment:是否多点登录"`
	FrontendUrl             string `json:"frontend_url"       gorm:"comment:官网地址"`
	EnableSwaggerApi        bool   `json:"enable_swagger_api" gorm:"comment:swagger api"`
	EnableAssetsApi         bool   `json:"enable_assets_api"  gorm:"comment:assets api"`
}

type Subscribe struct {
	BackendUrl                 string `json:"backend_url"       gorm:"comment:后端地址"`
	SubscribeDomainBindRequest bool   `json:"subscribe_domain_bind_request" gorm:"comment:订阅域名只接受更新订阅的请求"`
	SubName                    string `json:"sub_name"          gorm:"comment:订阅名称"`
	TEK                        string `json:"tek"               gorm:"comment:前后端通信密钥"`
	SurgeRule                  string `json:"surge_rule"        gorm:"comment:Surge 规则;type:text"`
	ClashRule                  string `json:"clash_rule"        gorm:"comment:Clash 规则;type:text"`
}
type Finance struct {
	EnableInvitationCommission bool    `json:"enable_invitation_commission"` //是否开启邀请佣金
	CommissionRate             float64 `json:"commission_rate"`              //佣金率, 范围 0~1, 佣金 = 订单金额 * 佣金率 ( 100.50 * 0.50 )
	WithdrawThreshold          float64 `json:"withdraw_threshold"`           //提取到余额的阈值

	EnableLottery bool    `json:"enable_lottery"` //是否开启每日打卡抽奖
	Jackpot       Jackpot `json:"jackpot"`        //奖池
}

// 公共配置参数
type PublicSystem struct {
	EnableRegister          bool    `json:"enable_register"`           // 是否开启注册
	AcceptableEmailSuffixes string  `json:"acceptable_email_suffixes"` // 可接受的邮箱后缀
	EnableBase64Captcha     bool    `json:"enable_base64_captcha"`     // 是否开启注册图片验证码
	EnableEmailCode         bool    `json:"enable_email_code"`         // 是否开启注册email 验证码
	EnableLoginEmailCode    bool    `json:"enable_login_email_code"`   // 是否开启登录email 验证码
	BackendUrl              string  `json:"backend_url"`               // 后端地址
	CommissionRate          float64 `json:"commission_rate"`           // 佣金率, 范围 0~1, 佣金 = 订单金额 * 佣金率 ( 100.50 * 0.50 )
	WithdrawThreshold       float64 `json:"withdraw_threshold"`        // 提取到余额的阈值
	EnableLottery           bool    `json:"enable_lottery"`            //是否开启每日打卡抽奖
	Jackpot                 Jackpot `json:"jackpot"`                   //奖池
	SubName                 string  `json:"sub_name"`                  //订阅名称
}
