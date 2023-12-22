package model

import "time"

// 数据库 traffic_log 流量统计表
type TrafficLog struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"      gorm:"primary_key"`
	NodeID    int64     `json:"node_id" gorm:"comment:节点ID"`
	U         int64     `json:"u"       gorm:"comment:上行流量 bit"`
	D         int64     `json:"d"       gorm:"comment:下行流量 bit"`
}
type UserTrafficLog struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"        gorm:"primary_key"`
	UserID    int64     `json:"user_id"   gorm:"comment:用户ID"`
	UserName  string    `json:"user_name" gorm:"comment:用户名"`
	U         int64     `json:"u"            gorm:"comment:上行流量 bit"`
	D         int64     `json:"d"            gorm:"comment:下行流量 bit"`
}
