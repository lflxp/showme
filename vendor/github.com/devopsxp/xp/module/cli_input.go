package module

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
)

// 1. 获取cli参数
// 2. 想办法传进来
// 3. 拼装成Message.Data map[string]interface{}
// 4. 执行pipeline

func init() {
	plugin.AddInput("cli", reflect.TypeOf(CliInput{}))
}

type CliInput struct {
	LifeCycle
	status       plugin.StatusPlugin
	connectcheck map[string]string
	lock         sync.RWMutex
	data         map[string]interface{}
	faileds      int // 失败次数
}

func (c *CliInput) Receive() *plugin.Message {
	if c.status != plugin.Started {
		slog.Warn("LocalYaml input plugin is not running,input nothing.")
		return nil
	}

	return plugin.Builder().WithInit(c.faileds).WithCheck(c.connectcheck).WithItemInterface(c.data).Build()
}

func (c *CliInput) SetConnectStatus(ip, status string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connectcheck[ip] = status
}

func (c *CliInput) Start() {
	c.faileds = 0
	c.status = plugin.Started
	slog.Debug("LocalYamlInput plugin started.")

	// Check all ipsl.yaml.
	// TODO: error 没有viper取配置了
	ips, err := getips(c.data["host"].([]string))
	if err != nil {
		panic(err)
	}

	var port int
	if pt, ok := c.data["remote_port"]; ok {
		port = pt.(int)
	} else {
		port = 22
	}

	// 目标主机22端口检测并发限制
	checkchan := make(chan string, 5*runtime.NumCPU())

	var wg sync.WaitGroup

	slog.Info("******************************************************** TASK [CliCheck : 主机状态检测] ********************************************************")
	slog.Info("LocalYaml Input 插件开始执行ssh目标主机状态扫描", "并发数", 5*runtime.NumCPU())
	for n, i := range ips {
		wg.Add(1)

		checkchan <- i
		go func(ip string, num int) {
			defer wg.Done()
			now := time.Now()
			if utils.ScanPort(ip, port) {
				slog.Info(fmt.Sprintf("%d: Ssh check %s:%d success 耗时: %v", num, ip, port, time.Now().Sub(now)))
				c.SetConnectStatus(ip, "success")
			} else {
				slog.Info(fmt.Sprintf("%d: Ssh check %s:%d failed 耗时：%v", num, ip, port, time.Now().Sub(now)))
				c.faileds += 1
				c.SetConnectStatus(ip, "failed")
			}
			<-checkchan
		}(i, n)

		if n%10 == 0 {
			slog.Info(fmt.Sprintf("已完成 %d 主机连接测试, 当前GoRoutine数量: %d", n, runtime.NumGoroutine()))
		}
	}

	wg.Wait()
}

// LocalYamlInput的Init函数实现
func (c *CliInput) Init(data interface{}) {
	c.connectcheck = make(map[string]string)
	c.name = "Cli Input"
	// 配置cli
	c.data = data.(map[string]interface{})
}
