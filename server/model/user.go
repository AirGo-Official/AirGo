package model

import (
	"sync"
	"time"
)

// 用户
type User struct {
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"-" gorm:"index"`
	ID             int64      `json:"id"           gorm:"primaryKey"`
	UserName       string     `json:"user_name"    gorm:"comment:用户名"`
	Password       string     `json:"password"     gorm:"comment:用户登录密码"`
	NickName       string     `json:"nick_name"    gorm:"default:系统用户;comment:用户昵称"`
	Avatar         string     `json:"avatar"       gorm:"comment:用户头像"`
	Enable         bool       `json:"enable"       gorm:"default:true;comment:用户是否被启用 true启用 false冻结"`
	InvitationCode string     `json:"invitation_code" gorm:"comment:我的邀请码"`
	ReferrerUserID int64      `json:"referrer_user_id"   gorm:"comment:推荐人user id"`
	Balance        float64    `json:"balance"         gorm:"default:0;comment:余额"`
	TgID           int64      `json:"tg_id"           gorm:"comment:tg id"`
	//关联参数
	RoleGroup []Role  `json:"role_group" gorm:"many2many:user_and_role;"`    //角色组，多对多
	Orders    []Order `json:"orders" gorm:"foreignKey:UserID;references:ID"` //订单，has many
	//通知参数
	EnableTGBot              bool `json:"enable_tg_bot"`
	EnableEmail              bool `json:"enable_email"`
	EnableWebMail            bool `json:"enable_web_mail"`
	WhenPurchased            bool `json:"when_purchased"`
	WhenServiceAlmostExpired bool `json:"when_service_almost_expired"`
	WhenBalanceChanged       bool `json:"when_balance_changed"`
}

// 用户与角色 多对多 表
type UserAndRole struct {
	UserID int64
	RoleID int64
}

// 用户登录/重置密码 请求
type UserLoginRequest struct {
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
type UserChangePasswordRequest struct {
	Password   string `json:"password" binding:"required,max=20,min=4"`                     // 密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password,max=20,min=4"` // 密码
	EmailCode  string `json:"email_code"`
}

// 修改头像 请求
type UserChangeAvatarRequest struct {
	Avatar string `json:"avatar"`
}

// 用户在线设备信息
type OnlineUserInfo struct {
	NodeConnector int64                    `json:"node_connector"` //连接客户端数
	NodeIPMap     map[int64]OnlineNodeInfo `json:"node_ip_map"`
}

type OnlineUserItem struct {
	NodeConnector int64     `json:"node_connector"` //连接客户端数
	NodeIPMap     *sync.Map `json:"node_ip_map"`    //key: nodeID int64   value: OnlineNodeInfo
}
type OnlineNodeInfo struct {
	NodeIP         []string  `json:"node_ip"`
	LastUpdateTime time.Time `json:"last_update_time"`
}
