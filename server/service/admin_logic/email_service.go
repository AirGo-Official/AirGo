package admin_logic

import (
	"crypto/tls"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"gopkg.in/gomail.v2"
)

// 订阅
const (
	EmailCode         = "EMAIL_CODE"
	Base64CaptchaCode = "Base64CaptchaCode"
)

var EmailSvc *Email

func InitEmailSvc() {
	EmailSvc = NewEmailService(global.Server.Email.EmailHost, int(global.Server.Email.EmailPort), global.Server.Email.EmailFrom, global.Server.Email.EmailSecret)
	EmailSvc.StartTask()
}

type Email struct {
	Dialer *gomail.Dialer
	//TaskChannel <-chan any
}
type EmailMsg struct {
	From      string
	To        string
	NickName  string
	Subject   string
	EmailText string
}

func NewEmailService(host string, port int, username string, password string) *Email {
	var e Email
	e.Dialer = gomail.NewDialer(host, port, username, password)
	// 关闭SSL协议认证
	e.Dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &e
}
func (es *Email) SendEmail(from, fromNickname, to, mailSubject, mailText string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(from, fromNickname)) // 发件人
	//m.SetHeader("From", from) // 发件人
	m.SetHeader("Subject", mailSubject) // 邮件主题
	m.SetHeader("To", to)               // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	m.SetBody("text/html", mailText)
	return es.Dialer.DialAndSend(m)
}
func (es *Email) StartTask() {
	ch, err := global.Queue.Subscribe(constant.EMAIL_CODE)
	if err != nil {
		global.Logrus.Error("Email StartTask error:", err)
		return
	}
	go func() {
		for v := range ch {
			emailMsg := v.(EmailMsg)
			err = es.SendEmail(emailMsg.From, emailMsg.NickName, emailMsg.To, emailMsg.Subject, emailMsg.EmailText)
			if err != nil {
				global.Logrus.Error("SendEmail error:", err)
			}
		}
	}()
}
