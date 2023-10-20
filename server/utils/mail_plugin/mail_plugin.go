package mail_plugin

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func InitEmailDialer(host string, port int, username string, password string) *gomail.Dialer {
	d := gomail.NewDialer(host, port, username, password)
	// 关闭SSL协议认证
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d
}

func SendEmail(d *gomail.Dialer, from, fromNickname, to, mailSubject, mailText string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromNickname+"<"+from+">") // 发件人
	//m.SetHeader("From", from) // 发件人

	m.SetHeader("Subject", mailSubject) // 邮件主题
	m.SetHeader("To", to)               // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接

	m.SetBody("text/html", mailText)
	return d.DialAndSend(m)

	//m := gomail.NewMessage()
	// m.SetHeader("From", "alias"+"<"+userName+">") // 增加发件人别名
	//m.SetHeader("To", "******@qq.com", "******@qq.com") // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	//m.SetHeader("Cc", "******@qq.com")                  // 抄送，可以多个
	//m.SetHeader("Bcc", "******@qq.com")                 // 暗送，可以多个
	// text/html 的意思是将文件的 content-type 设置为 text/html 的形式，浏览器在获取到这种文件时会自动调用html的解析器对文件进行相应的处理。可以通过 text/html 处理文本格式进行特殊处理，如换行、缩进、加粗等等
	// m.SetBody("text/html", fmt.Sprintf(message, "testUser"))
	// text/plain的意思是将文件设置为纯文本的形式，浏览器在获取到这种文件时并不会对其进行处理
	// m.SetBody("text/plain", "纯文本")
	// m.Attach("test.sh")   // 附件文件，可以是文件，照片，视频等等
	// m.Attach("lolcatVideo.mp4") // 视频
	// m.Attach("lolcat.jpg") // 照片
	//d := gomail.NewDialer(global.Server.Email.EmailHost, global.Server.Email.EmailPort, global.Server.Email.EmailFrom, global.Server.Email.EmailSecret)
	//// 关闭SSL协议认证
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	//d.DialAndSend(m)

}
