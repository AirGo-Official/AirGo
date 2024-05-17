package model

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 结构体-字符串 对应map
var (
	StringAndStruct = map[string]any{
		"user":                 User{},
		"order":                Order{},
		"dynamic_route":        Menu{},
		"role":                 Role{},
		"casbin_rule":          gormadapter.CasbinRule{},
		"goods":                Goods{},
		"traffic_log":          NodeTrafficLog{},
		"theme":                Theme{},
		"server":               Server{},
		"article":              Article{},
		"coupon":               Coupon{},
		"node":                 Node{},
		"pay":                  Pay{},
		"access":               Access{},
		"ticket":               Ticket{},
		"balance_statement":    BalanceStatement{},
		"commission_statement": CommissionStatement{},
	}
	StringAndSlice = map[string]any{
		"user":                 []User{},
		"order":                []Order{},
		"dynamic_route":        []Menu{},
		"role":                 []Role{},
		"casbin_rule":          []gormadapter.CasbinRule{},
		"goods":                []Goods{},
		"traffic_log":          []NodeTrafficLog{},
		"theme":                []Theme{},
		"server":               []Server{},
		"article":              []Article{},
		"coupon":               []Coupon{},
		"node":                 []Node{},
		"pay":                  []Pay{},
		"access":               []Access{},
		"ticket":               []Ticket{},
		"balance_statement":    []BalanceStatement{},
		"commission_statement": []CommissionStatement{},
	}
)
