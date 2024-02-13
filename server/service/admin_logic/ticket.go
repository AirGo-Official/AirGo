package admin_logic

import (
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"gorm.io/gorm"
)

type Ticket struct {
}

var TicketService *Ticket

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
