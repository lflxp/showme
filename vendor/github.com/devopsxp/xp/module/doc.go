package module

/*
模块文件夹主要用来实现INPUT|FILTER|OUTPUT 三类插件的功能。以满足pipeline 插件池的实现和功能填充。

主要接口定义：

/ 输入插件，用于接收消息，解析内容
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

插件生命周期管理：

// 插件抽象接口定义
type Plugin interface {
	// 启动插件
	Start()
	// 停止插件
	Stop()
	// 返回插件当前的运行状态
	Status() StatusPlugin
	// 新增初始化方法，在插件工厂返回实例前调用
	Init()
}
*/
