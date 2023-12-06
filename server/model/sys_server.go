package model

import "time"

// Server 全局配置
type Server struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id"   gorm:"primary_key"`

	Subscribe Subscribe `json:"subscribe"   gorm:"embedded"`
	Email     Email     `json:"email"    gorm:"embedded"`
	Security  Security  `json:"security" gorm:"embedded"`
	Notice    Notice    `json:"notice" gorm:"embedded"`
	//Acme      Acme      `json:"acme" gorm:"embedded"`
}
type Notice struct {
	BotToken           string `json:"bot_token"`            //tg bot token
	TGAdmin            string `json:"tg_admin"`             //tg admin
	TGSocks5           string `json:"tg_socks5"`            //tg socks5代理
	WhenUserRegistered bool   `json:"when_user_registered"` //用户注册后通知
	WhenUserPurchased  bool   `json:"when_user_purchased"`  //用户购买成功后通知
	WhenNodeOffline    bool   `json:"when_node_offline"`    //节点离线时通知
}

type Email struct {
	EmailFrom      string `json:"email_from"`                          // 发件人
	EmailFromAlias string `json:"email_from_alias"`                    // 发件人别名
	EmailSecret    string `json:"email_secret"`                        // 密钥
	EmailHost      string `json:"email_host"`                          // 服务器地址
	EmailPort      int64  `json:"email_port"`                          // 端口
	EmailIsSSL     bool   `json:"email_is_ssl"`                        // 是否SSL
	EmailNickname  string `json:"email_nickname"`                      // 昵称
	EmailSubject   string `json:"email_subject" gorm:"default:hello!"` // 邮件主题
	EmailContent   string `json:"email_content" gorm:"type:text"`      //邮件内容
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
	OpenCaptchaTimeOut int64 `json:"open_captcha_time_out" gorm:"default:300;comment:防爆破验证码超时时间，单位：s(秒)"`
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
	VisitParam  int64 `json:"visit_param"   gorm:"default:60"`
}
type Subscribe struct {
	EnableRegister          bool    `json:"enable_register"           gorm:"default:true;comment:是否开启注册"`
	AcceptableEmailSuffixes string  `json:"acceptable_email_suffixes" gorm:"comment:可接受的邮箱后缀"`
	EnableEmailCode         bool    `json:"enable_email_code"         gorm:"default:false;comment:是否开启注册email 验证码"`
	EnableLoginEmailCode    bool    `json:"enable_login_email_code"   gorm:"default:false;comment:是否开启登录email 验证码"`
	IsMultipoint            bool    `json:"is_multipoint"     gorm:"default:true;comment:是否多点登录"`
	BackendUrl              string  `json:"backend_url"       gorm:"comment:后端地址"`
	FrontendUrl             string  `json:"frontend_url"       gorm:"comment:官网地址"`
	ApiPrefix               string  `json:"api_prefix"        gorm:"default:/api;comment:api前缀"`
	SubName                 string  `json:"sub_name"          gorm:"default:AirGo;comment:订阅名称"`
	TEK                     string  `json:"tek"               gorm:"default:airgo;comment:前后端通信密钥"`
	DefaultGoods            int64   `json:"default_goods"     gorm:"comment:新用户默认套餐"`
	EnabledRebate           bool    `json:"enabled_rebate"    gorm:"default:false;comment:是否开启返利"`
	RebateRate              float64 `json:"rebate_rate"       gorm:"default:0.1;comment:返利率"`
	EnabledDeduction        bool    `json:"enabled_deduction" gorm:"default:false;comment:是否开启旧套餐抵扣"`
	DeductionThreshold      float64 `json:"deduction_threshold" gorm:"default:0.8;comment:旧套餐抵扣阈值,大于该值则抵扣"`
	EnabledClockIn          bool    `json:"enabled_clock_in"    gorm:"default:true;comment:是否开启打卡"`
	ClockInMinTraffic       int64   `json:"clock_in_min_traffic" gorm:"default:100;comment:打卡最小流量(MB)"`
	ClockInMaxTraffic       int64   `json:"clock_in_max_traffic" gorm:"default:1000;comment:打卡最大流量(MB)"`
	ClockInMinDay           int64   `json:"clock_in_min_day" gorm:"default:0;comment:打卡最小天数"`
	ClockInMaxDay           int64   `json:"clock_in_max_day" gorm:"default:1;comment:打卡最大天数"`
}
type Acme struct {
	//账户参数
	AcmeEmail   string `json:"acme_email"`
	AcmeMode    string `json:"acme_mode"    gorm:"default:dns"`         //dns http tls
	AccountType string `json:"account_type" gorm:"default:letsencrypt"` //letsencrypt zerossl buypass google
	KeyType     string `json:"key_type"`                                // P256 P384 2048 3072 4096
	Address     string `json:"address"`
	DNSProvider string `json:"dns_provider"`
	//证书信息
	IsExpired    bool      `json:"is_expired"`
	IsRenewal    bool      `json:"is_renewal"`
	ExpiredAt    time.Time `json:"expired_at"`
	StartAt      time.Time `json:"start_at"`
	PrivateKey   string    `json:"private_key"  gorm:"type:text"`
	Pem          string    `json:"pem" gorm:"type:text"`
	CertURL      string    `json:"cert_url"`
	CommonName   string    `json:"common_name"`
	Organization string    `json:"organization"`
	//其他
	EabKid     string `gorm:"default:null;" json:"eabKid"`
	EabHmacKey string `gorm:"default:null"  json:"eabHmacKey"`

	//TencentCloud
	TencentCloudSecretId  string `json:"tencent_cloud_secret_id"`
	TencentCloudSecretKey string `json:"tencent_cloud_secret_key"`
	//ALICLOUD
	AliCloudAccessKey string `json:"ali_cloud_access_key"`
	AliCloudSecretKey string `json:"ali_cloud_secret_key"`
	//
	CloudflareDnsApiToken string `json:"cloudflare_dns_api_token"`
	//
	GodaddyApiKey    string `json:"godaddy_api_key"`
	GodaddyApiSecret string `json:"godaddy_api_secret"`
	//
	HetznerApiKey string `json:"hetzner_api_key"`
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
