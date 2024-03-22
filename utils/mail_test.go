package utils

import "testing"

const FROM = "******"
const To = "******"
const Password = "******"
const Host = "smtp.163.com"
const Port = 465

func Test_Send(t *testing.T) {
	t.Run("发送邮件测试", func(t *testing.T) {
		err := NewEmail(FROM, Password, Host, Port).Send([]string{To}, "测试邮件", "这是一封测试邮件", "./mail.go")
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Log("邮件发送成功")
	})
}
