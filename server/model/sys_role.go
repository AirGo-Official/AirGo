package model

import "time"

type Role struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`

	ID          int64          `json:"id"   gorm:"primary_key;comment:角色ID"`         //
	RoleName    string         `json:"role_name" gorm:"comment:角色名"`                 // 角色名
	Status      bool           `json:"status"   gorm:"default:true;comment:角色状态"`    //
	Description string         `json:"description" gorm:"comment:描述"`                //
	UserGroup   []User         `json:"user_group"   gorm:"many2many:user_and_role;"` //多对多
	Menus       []DynamicRoute `json:"menus"       gorm:"many2many:role_and_menu;"`  //
	Nodes       []int64        `json:"nodes"       gorm:"-"`                         //接收前端选中的节点数组
}

// 角色 菜单 多对多
type RoleAndMenu struct {
	RoleID         int64
	DynamicRouteID int64
}

// roles with total
type RolesWithTotal struct {
	Total    int64  `json:"total"`
	RoleList []Role `json:"role_list"`
}
