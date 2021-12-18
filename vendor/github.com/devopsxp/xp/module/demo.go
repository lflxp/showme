package module

import (
	"fmt"
	"reflect"
	"strings"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	"github.com/spf13/viper"
)

func init() {
	// 初始化check插件映射关系表
	AddCheck("ssh", reflect.TypeOf(SshCheck{}))
	// 初始化input插件映射关系表
	AddInput("hello", reflect.TypeOf(HelloInput{}))
	// 初始化filter插件映射关系表
	AddFilter("upper", reflect.TypeOf(UpperFilter{}))
}

// ssh主机check插件
type SshCheck struct {
	status StatusPlugin
}

// 连接目标主机22端口进行测试
func (s *SshCheck) Conn() *Message {
	if s.status != Started {
		fmt.Println("Hello input plugin is not running,input nothing.")
		return nil
	}

	fmt.Println("Check target is connect")
	ip := viper.GetStringSlice("host")

	message := Builder().WithRaw("{'name':'xp'}").WithTarget(ip)

	for _, i := range ip {
		if utils.ScanPort(i, 22) {
			fmt.Printf("%s:22 is success\n", i)
			message.WithStatus(Ok).WithItems(i, "success")
		} else {
			fmt.Printf("%s:22 is failed\n", i)
			message.WithStatus(Error).WithItems(i, "failed")
		}
	}
	viper.SetDefault("ssh", message.Build().Data.Items)
	fmt.Printf("%v\n", message.Build().Data.Items)
	// 造假数据
	return message.Build()
}

func (s *SshCheck) Start() {
	s.status = Started
	fmt.Println("Check SshCheck plugin started.")
}

func (s *SshCheck) Stop() {
	s.status = Stopped
	fmt.Println("Check SshCheck plugin stopped.")
}

func (s *SshCheck) Status() StatusPlugin {
	return s.status
}

func (s *SshCheck) Init(data interface{}) {
	fmt.Println("Get machine and connecting test init")
}

// Hello input插件，接收“Hello World”消息
type HelloInput struct {
	status StatusPlugin
}

func (h *HelloInput) Receive() *Message {
	// 如果插件未启动，则返回nil
	if h.status != Started {
		fmt.Println("Hello input plugin is not running,input nothing.")
		return nil
	}
	return Builder().WithRaw("{'name':'xp'}").WithItems("thisis", "world").WithTarget([]string{"127.0.0.1", "192.168.0.1"}).WithStatus(Ok).Build()
}

func (h *HelloInput) Start() {
	h.status = Started
	fmt.Println("Hello input plugin started.")
}

func (h *HelloInput) Stop() {
	h.status = Stopped
	fmt.Println("Hello input plugin stopped.")
}

func (h *HelloInput) Status() StatusPlugin {
	return h.status
}

func (h *HelloInput) Init(data interface{}) {}

// Upper filter插件，将消息全部字母转成大写
type UpperFilter struct {
	status StatusPlugin
}

func (u *UpperFilter) Process(msgs *Message) *Message {
	if u.status != Started {
		fmt.Println("Upper filter plugin is not running ,filter nothing.")
		return msgs
	}

	for i, val := range msgs.Data.Target {
		msgs.Data.Target[i] = strings.ToUpper(val)
	}
	return msgs
}

func (u *UpperFilter) Start() {
	u.status = Started
	fmt.Println("Upper filter plugin started.")
}

func (u *UpperFilter) Stop() {
	u.status = Stopped
	fmt.Println("Upper filter plugin stopped.")
}

func (u *UpperFilter) Status() StatusPlugin {
	return u.status
}

func (u *UpperFilter) Init(data interface{}) {}
