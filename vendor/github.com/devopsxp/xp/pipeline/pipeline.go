package pipeline

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/devopsxp/xp/plugin"
	"github.com/google/gops/agent"
	log "github.com/sirupsen/logrus"
)

// 性能分析
func init() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
}

// Pipeline Config
type PipeConfig struct {
	Name    string
	TmpArgs interface{} // 临时变量插入
	Check   plugin.Config
	Input   plugin.Config
	Filter  plugin.Config
	Output  plugin.Config
}

// Pipeline Config 工厂模式
func DefaultPipeConfig(name string) *PipeConfig {
	return &PipeConfig{
		Name:   name,
		Check:  plugin.Config{PluginTypes: plugin.CheckType},
		Input:  plugin.Config{PluginTypes: plugin.InputType},
		Filter: plugin.Config{PluginTypes: plugin.FilterType},
		Output: plugin.Config{PluginTypes: plugin.OutputType},
	}
}

// 设置临时变量
func (p *PipeConfig) SetArgs(data interface{}) *PipeConfig {
	p.TmpArgs = data
	return p
}

func (p *PipeConfig) WithCheckName(name string) *PipeConfig {
	p.Check.Name = name
	return p
}

func (p *PipeConfig) WithInputName(name string) *PipeConfig {
	p.Input.Name = name
	return p
}

func (p *PipeConfig) WithFilterName(name string) *PipeConfig {
	p.Filter.Name = name
	return p
}

func (p *PipeConfig) WithOutputName(name string) *PipeConfig {
	p.Output.Name = name
	return p
}

// 对于插件化的系统，一切皆是插件，因此将pipeline也设计成一个插件，实现plugin接口
// pipeline管道的定义
type Pipeline struct {
	start   time.Time // 计时器
	status  plugin.StatusPlugin
	check   plugin.Check
	input   plugin.Input
	filter  plugin.Filter
	output  plugin.Output
	tmpargs interface{}
	spinner *spinner.Spinner
}

// 一个消息的处理流程 check -> input -> filter -> output
func (p *Pipeline) Exec() {
	// msg := p.check.Conn()
	msg := p.input.Receive()
	if msg.Status == plugin.Ok {
		msg = p.filter.Process(msg)
	}
	p.output.Send(msg)
}

// 启动的顺序 output -> filter -> input -> check
func (p *Pipeline) Start() {
	// p.spinner.Start()
	p.output.Start()
	p.filter.Start()
	p.input.Start()
	// p.check.Start()
	p.status = plugin.Started
	log.Debugln("Pipeline started.")
}

// 停止的顺序 check -> input -> filter -> output
func (p *Pipeline) Stop() {
	// defer p.spinner.Stop()
	// p.check.Stop()
	p.input.Stop()
	p.filter.Stop()
	p.output.Stop()
	p.status = plugin.Stopped
	log.Debugln("Pipeline stopped.")
	log.Infof("Pipeline执行完毕，耗时：%v", time.Now().Sub(p.start))
}

func (p *Pipeline) Status() plugin.StatusPlugin {
	return p.status
}

func (p *Pipeline) Init() {
	// p.spinner = spinner.New(spinner.CharSets[38], 100*time.Millisecond)
	p.start = time.Now()
	// p.check.Init()
	p.input.Init(p.tmpargs)
	p.filter.Init(p.tmpargs)
	p.output.Init(p.tmpargs)
}

// 最后定义pipeline的工厂方法，调用plugin.Factory抽象工厂完成pipelien对象的实例化：
// 保存用于创建Plugin的工厂实例，其中map的key为插件类型，value为抽象工厂接口
var pluginFactories = make(map[plugin.PluginType]plugin.Factory)

// 根据plugin.PluginType返回对应Plugin类型的工厂实例
func factoryOf(t plugin.PluginType) plugin.Factory {
	factory, _ := pluginFactories[t]
	return factory
}

// pipeline工厂方法，根据配置创建一个Pipeline实例
func Of(conf PipeConfig) *Pipeline {
	// 临时参数设置
	p := &Pipeline{}
	if conf.TmpArgs != nil {
		p.tmpargs = conf.TmpArgs
	}
	// p.check = factoryOf(plugin.CheckType).Create(conf.Check).(plugin.Check)
	p.input = factoryOf(plugin.InputType).Create(conf.Input).(plugin.Input)
	p.filter = factoryOf(plugin.FilterType).Create(conf.Filter).(plugin.Filter)
	p.output = factoryOf(plugin.OutputType).Create(conf.Output).(plugin.Output)
	return p
}

// 初始化插件工厂对象
func init() {
	// pluginFactories[plugin.CheckType] = &plugin.CheckFactory{}
	pluginFactories[plugin.InputType] = &plugin.InputFactory{}
	pluginFactories[plugin.FilterType] = &plugin.FilterFactory{}
	pluginFactories[plugin.OutputType] = &plugin.OutputFactory{}
}
