package service

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
)

func InitPushMessageSvc() {
	PushMessageSvc = NewPushMessageService()
	PushMessageSvc.StartTask()
}

var PushMessageSvc *PushMessageService

type PushMessageService struct {
	engine []Send
}

const (
	MESSAGE_TYPE_ADMIN = "admin"
	MESSAGE_TYPE_USER  = "user"
)

type MessageInfo struct {
	UserID      int64
	Message     string
	MessageType string
	User        *model.User
}
type Send interface {
	SendToAdmin(*MessageInfo) error
	SendToUser(*MessageInfo) error
}

func NewPushMessageService() *PushMessageService {
	return &PushMessageService{engine: []Send{
		&PushMessageByTGBot{},
		&PushMessageByWebMail{},
		&PushMessageByEmail{},
	}}
}

func (pm *PushMessageService) StartTask() {
	ch, err := global.Queue.Subscribe(constant.PUSH_MESSAGE, 1000)
	if err != nil {
		global.Logrus.Error("Tg bot StartTask error:", err)
		return
	}
	go func() {
		for v := range ch {
			msg := v.(*MessageInfo)
			for _, e := range pm.engine {
				switch msg.MessageType {
				case MESSAGE_TYPE_ADMIN:
					_ = e.SendToAdmin(msg)
				case MESSAGE_TYPE_USER:
					_ = e.SendToUser(msg)
				default:
				}
			}
		}
	}()
}
func (pm *PushMessageService) PushMessage(msg *MessageInfo) {
	_ = global.Queue.Publish(constant.PUSH_MESSAGE, msg)
}

type PushMessageByTGBot struct{}

func (p *PushMessageByTGBot) SendToAdmin(m *MessageInfo) error {
	if !global.Server.Notice.EnableTGBot {
		return errors.New("TGBot is disabled")
	}
	user, err := UserSvc.FirstUser(&model.User{ID: m.UserID})
	if err != nil {
		return err
	}
	if user.TgID != 0 {
		var msg tgbotapi.Chattable
		msg = tgbotapi.NewMessage(user.TgID, m.Message)
		TgBotSvc.TGBotSendMessage(msg)
	}
	return nil
}
func (p *PushMessageByTGBot) SendToUser(m *MessageInfo) error {
	if !m.User.EnableTGBot {
		return errors.New("TGBot is disabled")
	}
	if m.User.TgID != 0 {
		var msg tgbotapi.Chattable
		msg = tgbotapi.NewMessage(m.User.TgID, m.Message)
		TgBotSvc.TGBotSendMessage(msg)
	}
	return nil
}

type PushMessageByWebMail struct{}

func (p *PushMessageByWebMail) SendToAdmin(m *MessageInfo) error {
	//TODO implement me
	return nil
}
func (p *PushMessageByWebMail) SendToUser(m *MessageInfo) error {
	return nil
}

type PushMessageByEmail struct{}

func (p *PushMessageByEmail) SendToAdmin(m *MessageInfo) error {
	//TODO implement me
	return nil
}
func (p *PushMessageByEmail) SendToUser(m *MessageInfo) error {
	return nil
}
