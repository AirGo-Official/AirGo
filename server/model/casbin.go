package model

type CasbinInfo struct {
	RoleID      int64        `json:"roleID"` // 权限id
	CasbinItems []CasbinItem `json:"casbinItems"`
}

type CasbinItem struct {
	Path   string `json:"path"`   // 路径
	Method string `json:"method"` // 方法
}
