package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"sync"
)

// 结构体-字符串 对应map
var (
	StringAndStruct = map[string]any{
		"user":          User{},
		"order":         Order{},
		"dynamic_route": Menu{},
		"role":          Role{},
		"casbin_rule":   gormadapter.CasbinRule{},
		"goods":         Goods{},
		"traffic_log":   NodeTrafficLog{},
		"theme":         Theme{},
		"server":        Server{},
		"article":       Article{},
		"coupon":        Coupon{},
		"node":          Node{},
		"pay":           Pay{},
		"access":        Access{},
		"ticket":        Ticket{},
	}
	StringAndSlice = map[string]any{
		"user":          []User{},
		"order":         []Order{},
		"dynamic_route": []Menu{},
		"role":          []Role{},
		"casbin_rule":   []gormadapter.CasbinRule{},
		"goods":         []Goods{},
		"traffic_log":   []NodeTrafficLog{},
		"theme":         []Theme{},
		"server":        []Server{},
		"article":       []Article{},
		"coupon":        []Coupon{},
		"node":          []Node{},
		"pay":           []Pay{},
		"access":        []Access{},
		"ticket":        []Ticket{},
	}
)

// gorm 字符串切片类型
type SliceForGorm []string

func (s *SliceForGorm) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, s)
}
func (s SliceForGorm) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type ContextGroup struct {
	MapLock   sync.RWMutex
	CtxMap    map[string]*context.Context    //
	CancelMap map[string]*context.CancelFunc //
}
