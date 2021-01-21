package plugin

// 插件类型定义
type PluginType uint8

const (
	CheckType PluginType = iota
	InputType
	FilterType
	OutputType
)

// 添加插件运行状态
type StatusPlugin uint8

const (
	Stopped StatusPlugin = iota
	Started
)

// 接口定义
// Check、Input、Filter、Output三类插件接口的定义
// 插件抽象接口定义
type Plugin interface {
	// 启动插件
	Start()
	// 停止插件
	Stop()
	// 返回插件当前的运行状态
	Status() StatusPlugin
	// 新增初始化方法，在插件工厂返回实例前调用
	Init(interface{})
}

// 检测插件，用于检测目标单位可执行状态
type Check interface {
	Plugin
	Conn() *Message
}

// 输入插件，用于接收消息，解析内容
type Input interface {
	Plugin
	Receive() *Message
}

// 过滤插件，用于处理消息
type Filter interface {
	Plugin
	Process(msg *Message) *Message
}

// 输出插件，用于发送消息
type Output interface {
	Plugin
	Send(msg *Message)
}

// 插件配置
type Config struct {
	Name        string
	PluginTypes PluginType
}
