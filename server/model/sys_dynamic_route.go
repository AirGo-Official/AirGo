package model

import "time"

// 动态路由
type DynamicRoute struct {
	//gorm.Model
	CreatedAt time.Time                                        `json:"created_at"`
	UpdatedAt time.Time                                        `json:"updated_at"`
	DeletedAt *time.Time                                       `json:"-" gorm:"index"`
	ID        int64                                            `json:"id" gorm:"primary_key"`
	ParentID  int64                                            `json:"parent_id"  gorm:"comment:父菜单ID"`   // 父菜单ID
	Path      string                                           `json:"path"      gorm:"comment:路由path"`   // 路由path
	Name      string                                           `json:"name"      gorm:"comment:路由name"`   // 路由name
	Component string                                           `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Children  []DynamicRoute                                   `json:"children"   gorm:"-"`
	Roles     []Role                                           `json:"roles"        gorm:"many2many:role_and_menu;"`
	Meta      `json:"meta"       gorm:"embedded;comment:附加属性"` // 附加属性
}

type Meta struct {
	Title       string `json:"title"             gorm:"comment:route名称"`                                           //菜单栏及 tagsView 栏、菜单搜索名称（国际化）
	IsLink      string `json:"isLink"            gorm:"comment:是否超链接菜单,开启外链条件，1、isLink: 链接地址不为空 2、isIframe:false"` //是否超链接菜单，开启外链条件，`1、isLink: 链接地址不为空 2、isIframe:false`
	IsIframe    bool   `json:"isIframe"          gorm:"default:false;comment:是否内嵌窗口"`                              //是否内嵌窗口，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
	IsHide      bool   `json:"isHide"         gorm:"default:false;comment:是否隐藏此路由"`                                //是否隐藏此路由
	IsKeepAlive bool   `json:"isKeepAlive"    gorm:"default:true;comment:是否缓存组件状态"`                                //是否缓存组件状态
	IsAffix     bool   `json:"isAffix"        gorm:"default:false;comment:是否固定在 tagsView 栏上"`                      //是否固定在 tagsView 栏上
	Icon        string `json:"icon"           gorm:"default:iconfont icon-caidan;comment:菜单、tagsView 图标"`          // 菜单、tagsView 图标，阿里：加 `iconfont xxx`，fontawesome：加 `fa xxx`

}

// 菜单node，前端编辑角色的时候显示全部菜单节点
type RouteNode struct {
	ID       int64       `json:"route_id"`
	Title    string      `json:"title"`
	ParentID int64       `json:"-"`
	Children []RouteNode `gorm:"-" json:"children"`
}
