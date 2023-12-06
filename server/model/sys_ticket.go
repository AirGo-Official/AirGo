package model

import (
	"time"
)

const (
	TicketProcessing = "TicketProcessing"
	TicketClosed     = "TicketClosed"
)

type Ticket struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        int64     `json:"id"           gorm:"primary_key"`
	UserID    int64     `json:"user_id"      gorm:"comment:用户id"`
	Title     string    `json:"title"        gorm:"comment:工单标题"`
	Details   string    `json:"details"      gorm:"comment:工单详情"`
	Status    string    `json:"status"       gorm:"comment:工单状态"` //TicketProcessing TicketClosed
	//一对多关联
	TicketMessage []TicketMessage `json:"ticket_message" gorm:"foreignKey:TicketID;references:ID"`
}

type TicketMessage struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        int64     `json:"id"           gorm:"primary_key"`
	TicketID  int64     `json:"ticket_id"    gorm:"comment:工单id"`
	IsAdmin   bool      `json:"is_admin"     gorm:"comment:是否管理员"`
	Message   string    `json:"message"      gorm:"comment:消息"`
}
