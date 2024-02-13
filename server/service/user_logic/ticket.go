package user_logic

import (
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type Ticket struct {
}

func (t *Ticket) NewTicket(ticket *model.Ticket) error {
	ticket.Status = constant.TICKET_PROCESSING
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ticket).Error
	})
}
func (t *Ticket) UpdateUserTicket(ticketParams *model.Ticket) error {
	userTicket, err := t.FirstTicket(&model.Ticket{ID: ticketParams.ID, UserID: ticketParams.UserID})
	if err != nil {
		return err
	}
	userTicket.Title = ticketParams.Title
	userTicket.Details = ticketParams.Details
	userTicket.Status = constant.TICKET_PROCESSING
	if ticketParams.Status == constant.TICKET_CLOSED {
		userTicket.Status = constant.TICKET_CLOSED
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&userTicket).Error
	})
}
func (t *Ticket) FirstTicket(ticketParams *model.Ticket) (*model.Ticket, error) {
	var userTicket model.Ticket
	err := global.DB.Model(model.Ticket{}).Where(&ticketParams).Preload("TicketMessage").First(&userTicket).Error
	return &userTicket, err
}
func (t *Ticket) NewTicketMessage(msg *model.TicketMessage) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&msg).Error
	})
}
