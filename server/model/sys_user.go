package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// gorm 创建/更新时间追踪（纳秒、毫秒、秒、Time）
// CreatedAt time.Time // 在创建时，如果该字段值为零值，则使用当前时间填充
// UpdatedAt int64       // 在创建时该字段值为零值或者在更新时，使用当前时间戳秒数填充
// Updated   int64 `gorm:"autoUpdateTime:nano"` // 使用时间戳纳秒数填充更新时间
// Updated   int64 `gorm:"autoUpdateTime:milli"` // 使用时间戳毫秒数填充更新时间
// Created   int64 `gorm:"autoCreateTime"`      // 使用时间戳秒数填充创建时间
// 用户
type User struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`

	ID             int64     `json:"id"           gorm:"primary_key"`
	UUID           uuid.UUID `json:"uuid"         gorm:"comment:用户UUID"`
	UserName       string    `json:"user_name"    gorm:"comment:用户登录名"`
	Password       string    `json:"password"     gorm:"comment:用户登录密码"`
	Passwd         string    `json:"passwd"       gorm:"comment:用户密码"`
	NickName       string    `json:"nick_name"    gorm:"default:系统用户;comment:用户昵称"`
	Avatar         string    `json:"avatar"       gorm:"default:/src/assets/icon/avatar.jpeg;comment:用户头像"`
	Phone          string    `json:"phone"        gorm:"comment:用户手机号"`
	Enable         bool      `json:"enable"       gorm:"default:true;comment:用户是否被启用 true正常 false冻结"`
	InvitationCode string    `json:"invitation_code" gorm:"comment:我的邀请码"`
	ReferrerCode   string    `json:"referrer_code"   gorm:"comment:推荐人码"`
	Remain         float64   `json:"remain"          gorm:"comment:余额"`
	//角色组
	RoleGroup []Role `json:"role_group" gorm:"many2many:user_and_role;"` //多对多
	//订单
	Orders []Orders `json:"orders" gorm:"foreignKey:UserID;references:ID"` //has many
	//附加订阅信息
	SubscribeInfo SubscribeInfo `json:"subscribe_info" gorm:"embedded;comment:附加订阅信息"`
}

// 附加订阅信息
type SubscribeInfo struct {
	Host           string     `json:"host"              gorm:"comment:用户混淆"`
	ClientIP       string     `json:"client_ip"         gorm:"comment:用户的当前在线IP"`
	SubStatus      bool       `json:"sub_status"        gorm:"comment:订阅是否有效"`
	SubscribeUrl   string     `json:"subscribe_url"     gorm:"comment:订阅链接"`
	GoodsID        int64      `json:"goods_id"          gorm:"comment:商品ID"`
	GoodsSubject   string     `json:"goods_subject"     gorm:"comment:商品标题"`
	ExpiredAt      *time.Time `json:"expired_at"        gorm:"comment:过期时间"`
	T              int64      `json:"t"                 gorm:"default:0;comment:总流量（Byte）"`
	U              int64      `json:"u"                 gorm:"default:0;comment:上行流量"`
	D              int64      `json:"d"                 gorm:"default:0;comment:下行流量"`
	ResetDay       int64      `json:"reset_day"         gorm:"comment:流量重置日"`
	NodeSpeedLimit int64      `json:"node_speedlimit"   gorm:"default:0;comment:限速Mbps（Mbps）"`
	NodeConnector  int64      `json:"node_connector"    gorm:"default:3;comment:连接客户端数"`
}

// 用户与角色 多对多 表
type UserAndRole struct {
	UserID int64
	RoleID int64
}

// 用户登录/重置密码 请求
type UserLogin struct {
	UserName  string `json:"user_name" binding:"required,email,max=40,min=4"` // 用户名
	Password  string `json:"password" binding:"required,max=20,min=4"`        // 密码
	EmailCode string `json:"email_code"`                                      //邮箱验证码
}

// 用户注册 请求
type UserRegister struct {
	UserName      string            `json:"user_name"    binding:"required,max=40,min=4"`                  // 用户名
	EmailSuffix   string            `json:"email_suffix" binding:"required,max=40"`                        // 邮箱后缀
	Password      string            `json:"password"     binding:"required,max=20,min=4"`                  // 密码
	RePassword    string            `json:"re_password"  binding:"required,eqfield=Password,max=20,min=4"` // 密码
	EmailCode     string            `json:"email_code"`                                                    //邮箱验证码
	ReferrerCode  string            `json:"referrer_code"`
	Base64Captcha Base64CaptchaInfo `json:"base64_captcha"`
}

// 用户注册 校验邮箱
type UserRegisterEmail struct {
	UserName string `json:"user_name" binding:"required,email,max=40,min=4"` // 用户名
}

// 修改密码 请求
type UserChangePassword struct {
	Password   string `json:"password" binding:"required,max=20,min=4"`                     // 密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password,max=20,min=4"` // 密码
	EmailCode  string `json:"email_code"`
}

// users with total
type UsersWithTotal struct {
	Total    int64  `json:"total"`
	UserList []User `json:"user_list"`
}

// 获取当前节点可连接的用户 响应
type SSUsers struct {
	ID             int64  `json:"id"`
	UUID           string `json:"uuid"`
	U              int64  `json:"u"`
	D              int64  `json:"d"`
	NodeSpeedLimit int64  `json:"node_speedlimit"`
	NodeConnector  int64  `json:"node_connector"`
	// IsMultiUser    int64    `json:"is_multi_user"` //用不到
}
