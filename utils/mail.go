package utils

import (
	"gopkg.in/gomail.v2"
)

type Email struct {
	User     string
	Password string
	Host     string
	Port     int
}

func NewEmail(user string, password string, host string, port int) *Email {
	return &Email{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (e *Email) Send(mailTo []string, subject string, body string, attachs ...string) error {
	// 设置邮箱主体
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.User, "补网人")) // 添加别名
	m.SetHeader("To", mailTo...)                        // 发送给用户(可以多个)
	m.SetHeader("Subject", subject)                     // 设置邮件主题
	m.SetBody("text/html", body)                        // 设置邮件正文
	// 附件
	for _, v := range attachs {
		m.Attach(v)
	}
	d := gomail.NewDialer(e.Host, e.Port, e.User, e.Password) // 设置邮件正文
	err := d.DialAndSend(m)
	return err
}
