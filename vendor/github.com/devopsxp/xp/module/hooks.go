package module

import (
	"time"

	. "github.com/devopsxp/xp/plugin"
	log "github.com/sirupsen/logrus"
)

// 对外接口
type Alert interface {
	Send(*Message)
}

// 对内实现
type HookMethod interface {
	SpecificSend() (string, error)
	IsCurrent() bool // 判断当前是否可以发送告警信
}

// hook适配器
func NewHookAdapter(in HookMethod) *hook {
	return &hook{
		HookMethod: in,
	}
}

// 转换对外接口调用对内接口
// output 钩子结构体
// 负责处理发送
type hook struct {
	HookMethod
	Type   string
	Target string
	start  time.Time // 计时器
}

func (h *hook) Send(msg *Message) {
	h.start = time.Now()
	switch h.Type {
	case "count":
		for k, v := range msg.Count {
			log.Warnf("%s : ok=%d   changed=%d failed=%d  skipped=%d rescued=%d  ignored=%d", k, v["ok"], v["changed"], v["failed"], v["skipped"], v["rescued"], v["ignored"])
		}
		log.Warnf("count日志耗时：%v", time.Now().Sub(h.start))
	case "console":
		log.Debugf("console hook send %v\n", msg.Data.Check)
		for k, v := range msg.CallBack {
			log.Warnln(k, v)
		}
		log.Warnf("console日志耗时：%v", time.Now().Sub(h.start))
	default:
		log.Debugln("email hook send")
		status := h.IsCurrent()
		if !status {
			log.Warnln("不在发送时间，停止发送")
		} else {
			rs, err := h.SpecificSend()
			if err != nil {
				log.Errorln(err)
			} else {
				log.Warnln(rs)
			}
		}
		log.Warnf("%s 发送耗时：%v", h.Type, time.Now().Sub(h.start))
	}
}

func (h *hook) SetType(t string) *hook {
	h.Type = t
	return h
}

func (h *hook) SetTarget(target string) *hook {
	h.Target = target
	return h
}
