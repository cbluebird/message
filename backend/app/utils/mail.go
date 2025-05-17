package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"message/config/config"
)

type EmailConfig struct {
	Title     string
	Sender    string
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}

var emailConfig EmailConfig
var mailSender *gomail.Dialer

func init() {
	emailConfig = EmailConfig{
		Sender:    config.Config.GetString("email.sender"),
		SPassword: config.Config.GetString("email.password"),
		SMTPAddr:  config.Config.GetString("email.smtp_addr"), // SMTP地址 QQ邮箱是smtp.qq.com
		SMTPPort:  config.Config.GetInt("email.smtp_port"),    // SMTP端口 QQ邮箱是465
		Title:     "您的电子邮箱验证码",
	}
	mailSender = gomail.NewDialer(emailConfig.SMTPAddr, emailConfig.SMTPPort, emailConfig.Sender, emailConfig.SPassword)
}

func SendEmail(mail, code string) error {
	html := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>你本次的验证码为%s,为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>    
    </div>`, code)
	m := gomail.NewMessage()
	// 第三个参数是我们发送者的名称，但是如果对方有发送者的好友，优先显示对方好友备注名
	m.SetHeader(`From`, emailConfig.Sender)
	m.SetHeader(`To`, []string{mail}...)
	m.SetHeader(`Subject`, emailConfig.Title)
	m.SetBody(`text/html`, html)
	err := mailSender.DialAndSend(m)
	if err != nil {
		log.Println(err)
	}
	return nil
}
