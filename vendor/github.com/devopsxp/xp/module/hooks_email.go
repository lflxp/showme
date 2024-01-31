package module

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	"gopkg.in/gomail.v2"
)

type Email struct {
	Sid      string   `json:"sid"`
	From     string   `json:"from"`
	Name     string   `json:"name"`
	To       []string `json:"to"`
	Smtp     string   `json:"smtp"`
	SmtpPort int      `json:"smtpport"`
	Pwd      string   `json:"pwd"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	Range    string   `json:"range"` // 启动时间
}

func NewEmail(data map[interface{}]interface{}, msg *plugin.Message, vars map[string]interface{}) (*Email, error) {
	result := &Email{Range: "星期一,星期二,星期三,星期四,星期五,星期六,星期天|00:00-23:59", Sid: utils.GetRandomSalt()}
	// 发送人
	if user, ok := data["email_user"]; ok {
		result.From = user.(string)
	} else {
		return result, errors.New("email_user 字段不存在")
	}

	// 昵称
	if user, ok := data["alias"]; ok {
		result.Name = user.(string)
	} else {
		result.Name = "github.com/devopsxp/xp"
	}

	// 接收人
	if to, ok := data["email_to"]; ok {
		for _, t := range to.([]interface{}) {
			result.To = append(result.To, t.(string))
		}
	} else {
		return result, errors.New("email_to 字段不存在")
	}

	// 邮箱密码
	if pwd, ok := data["email_pwd"]; ok {
		result.Pwd = pwd.(string)
	} else {
		return result, errors.New("email_pwd 字段不存在")
	}

	// smtp 服务器
	if smtp, ok := data["email_smtp"]; ok {
		result.Smtp = smtp.(string)
	} else {
		return result, errors.New("email_smtp 字段不存在")
	}

	// smtp 端口
	if smtpport, ok := data["email_smtp_port"]; ok {
		port, err := strconv.Atoi(smtpport.(string))
		if err != nil {
			return result, err
		}
		result.SmtpPort = port
	} else {
		return result, errors.New("email_smtp_port 字段不存在")
	}

	// 模板解析
	if template, ok := data["template"]; ok {
		// 标题
		if title, ok := template.(map[interface{}]interface{})["title"]; ok {
			result.Subject = title.(string)
		} else {
			result.Subject = "未设置标题"
		}
		// 模板解析
		// 1. 优先text字段 如果text字段不存在再取读取path字段进行文件读取
		// 2. 日志内容以logs代替
		if text, ok := template.(map[interface{}]interface{})["text"]; ok {
			// 加载内置固定参数
			if vs, ok := template.(map[interface{}]interface{})["vars"]; ok {
				for key, value := range vs.(map[interface{}]interface{}) {
					vars[key.(string)] = value
				}
			}
			// 加载日志
			vars["logs"] = msg.CallBack
			temp, err := utils.ApplyTemplate(text.(string), vars)
			if err != nil {
				slog.Debug(fmt.Sprintf("告警模板解析失败: %s", err.Error()))
				return result, err
			}
			result.Body = temp
		} else {
			// 解析path字段
			if path, ok := template.(map[interface{}]interface{})["path"]; ok {
				// 加载内置固定参数
				if vs, ok := template.(map[interface{}]interface{})["vars"]; ok {
					for key, value := range vs.(map[interface{}]interface{}) {
						vars[key.(string)] = value
					}
				}
				// 加载日志
				vars["logs"] = msg.CallBack

				// 读取模板
				tempFile, err := ioutil.ReadFile(path.(string))
				if err != nil {
					return nil, err
				}

				temp, err := utils.ApplyTemplate(string(tempFile), vars)
				if err != nil {
					slog.Debug(fmt.Sprintf("告警模板解析失败: %s", err.Error()))
					return result, err
				}
				result.Body = temp
			} else {
				return result, errors.New("告警模板无 text or path字段，请检查配置")
			}
		}
	} else {
		return result, errors.New("template 字段不存在")
	}

	return result, nil
}

func (this *Email) SpecificSend() (string, error) {
	m := gomail.NewMessage()

	m.SetAddressHeader("From", this.From, this.Name)
	// m.SetHeader("From", this.From)
	m.SetHeader("To", this.To...)
	// m.SetAddressHeader("Cc", "xp@xp.com", "Li")
	m.SetHeader("Subject", this.Subject)
	m.SetBody("text/html", this.Body)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(this.Smtp, this.SmtpPort, this.From, this.Pwd)

	// Send the email to Bob, Cora and Dan.
	err := d.DialAndSend(m)
	// log.Println("邮件发送结果", err)

	slog.Info("邮件发送结果", "ERROR", err)

	// if err != nil {
	// 	this.CallBack(this.Sid, err.Error())
	// } else {
	// 	this.CallBack(this.Sid, "success")
	// }

	if err != nil {
		return err.Error(), err
	}
	return fmt.Sprintf("send email success"), err
}

// 'range': '星期一,星期二,星期三,星期四,星期五,星期六,星期天|01:05-23:59'
// 'range': '|00:00-23:59'
// 'range': '-'
func (this *Email) IsCurrent() bool {
	rs := false
	// 如果包含 | 字符且包含不为空的字符
	if strings.Contains(this.Range, "|") && strings.Count(this.Range, "星") > 0 && strings.Contains(this.Range, "-") && strings.Count(this.Range, ":") == 2 {
		tmp := strings.Split(this.Range, "|")
		rangeTime, err := utils.TransformCHN(strings.Split(tmp[0], ","))
		if err != nil {
			slog.Info(err.Error())
			return rs
		}
		now := time.Now()
		now_weekday := int(now.Weekday())

		for _, x := range rangeTime {
			// 判断当前星期几是否在可执行范围内
			if now_weekday == x {
				rangeTime2 := strings.Split(tmp[1], "-")
				// 如果时间格式不是两个范围
				if len(rangeTime2) != 2 {
					break
					return rs
				}
				rs, err = utils.IsBetweenAB(rangeTime2[0], rangeTime2[1])
			}
		}
	}
	return rs
}
