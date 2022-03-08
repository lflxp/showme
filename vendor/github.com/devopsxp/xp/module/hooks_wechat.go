package module

// https://work.weixin.qq.com/api/doc?notreplace=true#90000/90136/91770

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

type Text struct {
	Content string `json:"content"`
}

type Wechat struct {
	Sid         string   `json:"sid"`
	Msgtype     string   `json:"msgtype"`
	Articles    string   `json:"articles"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Picurl      string   `json:"picurl"`
	Text        Text     `json:"text"`
	Address     []string `json:"address"`
	Markdown    Text     `json:"markdown"`
	Origin      string
	Range       string `json:"range"`
}

func NewWechat(data map[interface{}]interface{}, msg *plugin.Message, vars map[string]interface{}) (*Wechat, error) {
	result := &Wechat{
		Sid:     utils.GetRandomSalt(),
		Msgtype: "markdown",
		Range:   "星期一,星期二,星期三,星期四,星期五,星期六,星期天|00:00-23:59",
		Address: []string{},
	}

	// 配置企业微信地址
	if address, ok := data["address"]; ok {
		for _, add := range address.([]interface{}) {
			result.Address = append(result.Address, add.(string))
		}
	} else {
		return nil, errors.New("address 字段不存在")
	}

	// 配置模板发送内容
	if template, ok := data["template"]; ok {
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
				log.Debugf("告警模板解析失败: %s", err.Error())
				return result, err
			}
			result.Markdown = Text{Content: temp}
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
					log.Debugf("告警模板解析失败: %s", err.Error())
					return result, err
				}
				result.Markdown = Text{Content: temp}
			} else {
				return result, errors.New("告警模板无 text or path字段，请检查配置")
			}
		}
	} else {
		return result, errors.New("template 字段不存在")
	}

	// 生成Origin数据
	result.Origin = result.String()

	return result, nil
}

// 转换成string
func (this *Wechat) String() string {
	info, _ := json.Marshal(this)
	return string(info)
}

func (this *Wechat) SpecificSend() (string, error) {
	// data, err := json.Marshal(this)
	// if err != nil {
	// 	return err
	// }

	var (
		rs  string
		err error
	)
	log.WithFields(log.Fields{
		"uuid":   this.Sid,
		"type":   "wechat",
		"title":  this.Title,
		"method": "SpecificSend()",
	}).Info(this.Origin)
	for _, address := range this.Address {
		rs, err = post(address, this.Origin)
		if err != nil {
			log.WithFields(log.Fields{
				"uuid":    this.Sid,
				"type":    "wechat",
				"title":   this.Title,
				"address": address,
				"method":  "SpecificSend()",
			}).Error(err.Error())
			// this.CallBack(this.Sid, err.Error())
		} else {
			log.WithFields(log.Fields{
				"uuid":    this.Sid,
				"type":    "wechat",
				"title":   this.Title,
				"address": address,
				"method":  "SpecificSend()",
			}).Info("发送成功")
			// this.CallBack(this.Sid, "success")
		}
	}

	return rs, err
}

func post(url, message string) (string, error) {
	var rs string
	client := http.Client{}
	request, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return rs, err
	}

	// request.SetBasicAuth("yunxiao", "message")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return rs, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"url":    url,
			"method": "post",
		}).Error(err.Error())
		return rs, err
	}
	log.WithFields(log.Fields{
		"url":    url,
		"method": "post",
	}).Info("Http Post Response ", string(body))
	rs = string(body)
	return rs, nil
}

func (this *Wechat) IsCurrent() bool {
	rs := false
	// 如果包含 | 字符且包含不为空的字符
	if strings.Contains(this.Range, "|") && strings.Count(this.Range, "星") > 0 && strings.Contains(this.Range, "-") && strings.Count(this.Range, ":") == 2 {
		tmp := strings.Split(this.Range, "|")
		rangeTime, err := utils.TransformCHN(strings.Split(tmp[0], ","))
		if err != nil {
			log.Println(err.Error())
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
				}
				rs, err = utils.IsBetweenAB(rangeTime2[0], rangeTime2[1])
			}
		}
	}
	return rs
}
