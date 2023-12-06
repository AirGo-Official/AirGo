package model

// Casbin info structure
type CasbinItem struct {
	Path   string `json:"path"`   // 路径
	Method string `json:"method"` // 方法
}

// Casbin structure for input parameters
type CasbinInfo struct {
	RoleID      int64        `json:"roleID"` // 权限id
	CasbinItems []CasbinItem `json:"casbinItems"`
}

// 修改角色权限 请求
type ChangeRoleCasbinReq struct {
	RoleID      int64    `json:"roleID"` // 权限id
	CasbinItems []string `json:"casbinItems"`
}
