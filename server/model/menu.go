package model

import "time"

// 动态路由
type Menu struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id" gorm:"primaryKey"`
	ParentID  int64      `json:"parent_id" gorm:"comment:父菜单ID"`
	Remarks   string     `json:"remarks"   gorm:"comment:别名"`
	Path      string     `json:"path"      gorm:"comment:路由path"`
	Name      string     `json:"name"      gorm:"comment:路由name"`
	Component string     `json:"component" gorm:"comment:对应前端文件路径"`
	Children  []Menu     `json:"children"  gorm:"-"`
	Roles     []Role     `json:"roles"     gorm:"many2many:role_and_menu;"`
	Meta      `json:"meta"       gorm:"embedded;comment:附加属性"`
}

type Meta struct {
	Title       string `json:"title"          gorm:"comment:route名称"`
	IsLink      string `json:"isLink"         gorm:"comment:是否超链接菜单,开启外链条件，1、isLink: 链接地址不为空 2、isIframe:false"`
	IsIframe    bool   `json:"isIframe"       gorm:"default:false;comment:是否内嵌窗口"`
	IsHide      bool   `json:"isHide"         gorm:"default:false;comment:是否隐藏此路由"`
	IsKeepAlive bool   `json:"isKeepAlive"    gorm:"default:true;comment:是否缓存组件状态"`
	IsAffix     bool   `json:"isAffix"        gorm:"default:false;comment:是否固定在 tagsView 栏上"`
	Icon        string `json:"icon"           gorm:"default:iconfont icon-caidan;comment:菜单、tagsView 图标"`
}
