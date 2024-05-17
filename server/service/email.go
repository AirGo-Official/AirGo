package service

import (
	"crypto/tls"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"gopkg.in/gomail.v2"
)

var EmailSvc *Email

func InitEmailSvc() {
	EmailSvc = NewEmailService()
	EmailSvc.StartTask()
}

type Email struct {
	Dialer *gomail.Dialer
}
type EmailMsg struct {
	From      string
	To        string
	NickName  string
	Subject   string
	EmailText string
}

func NewEmailService() *Email {
	var e Email
	e.Dialer = NewDialer()
	// 关闭SSL协议认证
	e.Dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &e
}
func NewDialer() *gomail.Dialer {
	return gomail.NewDialer(global.Server.Email.EmailHost, int(global.Server.Email.EmailPort), global.Server.Email.EmailFrom, global.Server.Email.EmailSecret)
}
func (es *Email) Reload() {
	es.Dialer = NewDialer()
}
func (es *Email) SendEmail(msg *EmailMsg) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(msg.From, msg.NickName)) // 发件人
	//m.SetHeader("From", from) // 发件人
	m.SetHeader("Subject", msg.Subject) // 邮件主题
	m.SetHeader("To", msg.To)           // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	m.SetBody("text/html", msg.EmailText)
	return es.Dialer.DialAndSend(m)
}
func (es *Email) StartTask() {
	ch, err := global.Queue.Subscribe(constant.SEND_EMAIL)
	if err != nil {
		global.Logrus.Error("Email StartTask error:", err)
		return
	}
	go func() {
		for v := range ch {
			emailMsg := v.(*EmailMsg)
			err = es.SendEmail(emailMsg)
			if err != nil {
				global.Logrus.Error("SendEmail error:", err)
			}
		}
	}()
}
func (es *Email) PushEmailToQueue(email *EmailMsg) {
	global.Queue.Publish(constant.SEND_EMAIL, email)
}
