package utils

import (
	"fmt"
	"os"

	log "log/slog"

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

func (e *Email) Send(mailTo []string, subject string, body string, video string, attachs ...string) error {
	// 设置邮箱主体
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.User, "补网人")) // 添加别名
	m.SetHeader("To", mailTo...)                        // 发送给用户(可以多个)
	m.SetHeader("Subject", subject)                     // 设置邮件主题
	m.SetBody("text/html", body)                        // 设置邮件正文
	// 附件
	log.Info("发送邮件", "MailTo", mailTo, "attachs", attachs, "video", video)
	for _, v := range attachs {
		m.Attach(v, gomail.SetHeader(map[string][]string{
			"Content-ID":          {"<myImage>"},
			"Content-Disposition": {fmt.Sprintf("inline; filename='%s'", v)},
		}))
		m.Attach(v)
	}
	if video != "" {
		m.Attach(video)
	}
	d := gomail.NewDialer(e.Host, e.Port, e.User, e.Password) // 设置邮件正文
	err := d.DialAndSend(m)
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		for _, v := range attachs {
			e.Clean(v)
		}
		e.Clean(video)
	}()
	return err
}

// 清空文件
func (e *Email) Clean(path string) error {
	return os.Remove(path)
}
