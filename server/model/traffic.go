package model

import "time"

// 节点流量统计
type NodeTrafficLog struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"      gorm:"primaryKey"`
	NodeID    int64     `json:"node_id" gorm:"comment:节点ID"`
	U         int64     `json:"u"       gorm:"comment:上行流量 bit"`
	D         int64     `json:"d"       gorm:"comment:下行流量 bit"`
}

// 用户订阅流量统计
type UserTrafficLog struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"          gorm:"primaryKey"`
	UserName  string    `json:"user_name"   gorm:"comment:用户名"`
	SubUserID int64     `json:"sub_user_id" gorm:"comment:订阅user id(使用服务的id)"`
	U         int64     `json:"u"           gorm:"comment:上行流量 bit"`
	D         int64     `json:"d"           gorm:"comment:下行流量 bit"`
}
