package model

import gormadapter "github.com/casbin/gorm-adapter/v3"

// 分页参数
type PaginationParams struct {
	Search   string   `json:"search"`
	PageNum  int64    `json:"page_num"`
	PageSize int64    `json:"page_size"`
	Total    int64    `json:"total"`
	Date     []string `json:"date"`
}

// 结构体-字符串 对应map
var (
	StringAndStruct = map[string]any{
		"user":          User{},
		"orders":        Orders{},
		"dynamic_route": DynamicRoute{},
		"role":          Role{},
		"casbin_rule":   gormadapter.CasbinRule{},
		"goods":         Goods{},
		"traffic_log":   TrafficLog{},
		"theme":         Theme{},
		"server":        Server{},
		"gallery":       Gallery{},
		"article":       Article{},
		"coupon":        Coupon{},
		"isp":           ISP{},
		"node_shared":   NodeShared{},
		"node":          Node{},
		"pay":           Pay{},
	}
	StringAndSlice = map[string]any{
		"user":          []User{},
		"orders":        []Orders{},
		"dynamic_route": []DynamicRoute{},
		"role":          []Role{},
		"casbin_rule":   []gormadapter.CasbinRule{},
		"goods":         []Goods{},
		"traffic_log":   []TrafficLog{},
		"theme":         []Theme{},
		"server":        []Server{},
		"gallery":       []Gallery{},
		"article":       []Article{},
		"coupon":        []Coupon{},
		"isp":           []ISP{},
		"node_shared":   []NodeShared{},
		"node":          []Node{},
		"pay":           []Pay{},
	}
)
