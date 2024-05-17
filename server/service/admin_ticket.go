package service

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type AdminTicket struct{}

var AdminTicketSvc *AdminTicket

func (t *AdminTicket) FirstTicket(ticketParams *model.Ticket) (*model.Ticket, error) {
	var userTicket model.Ticket
	err := global.DB.Model(model.Ticket{}).Where(&ticketParams).Preload("TicketMessage").First(&userTicket).Error
	return &userTicket, err
}
func (t *AdminTicket) NewTicketMessage(msg *model.TicketMessage) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&msg).Error
	})
}
