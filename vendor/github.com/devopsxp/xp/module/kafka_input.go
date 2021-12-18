package module

import (
	"fmt"
	"reflect"

	. "github.com/devopsxp/xp/plugin"
)

func init() {
	// 在输入插件映射关系中加入kafka，用于通过反射创建input对象
	AddInput("kafka", reflect.TypeOf(KafkaInput{}))
}

// 原始struct
type Records struct {
	Items []string
}

type Consumer interface {
	Poll() *Records
}

// 将上面原始struct转换成目标struct KafkaInput
// 重点：转换
// 特殊功能：添加Plugin Func
type KafkaInput struct {
	status   StatusPlugin
	consumer Consumer
}

func (k *KafkaInput) Receive() *Message {
	records := k.consumer.Poll()
	if k.status != Started {
		fmt.Println("Kafka input plugin is not running, input nothing.")
		return nil
	}

	return Builder().WithRaw("{'name':'kafka'}").WithItems("name", "kafka").WithTarget(records.Items).Build()
}

func (k *KafkaInput) Start() {
	k.status = Started
	fmt.Println("KafkaInput plugin started.")
}

func (k *KafkaInput) Stop() {
	k.status = Stopped
	fmt.Println("KafkaInput plugin stopped.")
}

func (k *KafkaInput) Status() StatusPlugin {
	return k.status
}

// KakkaInput的Init函数实现
func (k *KafkaInput) Init(data interface{}) {
	k.consumer = &MockConsumer{}
}

// 上述代码中的kafka.MockConsumer为我们模式Kafka消费者的一个实现，代码如下
type MockConsumer struct{}

func (m *MockConsumer) Poll() *Records {
	records := &Records{}
	records.Items = append(records.Items, "i am mock consumer.")
	return records
}
