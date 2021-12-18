package roles

import (
	"errors"
	"fmt"
	"time"

	. "github.com/devopsxp/xp/plugin"
	log "github.com/sirupsen/logrus"
)

// 执行过程生命周期接口
// Role LifeCycle
type RolePlugin interface {
	// 初始化对象
	Init(*RoleArgs) error

	// 环境准备
	Pre()

	// 执行前
	Before()

	// 执行中
	// 返回是否执行信号
	Run() error

	// 执行后
	After()

	// 是否执行hook
	IsHook() bool

	// 钩子函数
	Hooks() error
}

type Hook struct {
	isHook   bool                    // 是否执行hook
	hookArgs []string                // hook 函数参数
	hookFunc func(...[]string) error // hook 执行函数
}

// Role生命周期 公共父类
type RoleLC struct {
	name  string // 名称
	types string
	// 通用字段
	stage       string
	remote_user string                 // 执行用户
	remote_pwd  string                 // 执行用户密码（非必须）
	remote_port int                    // ssh端口
	vars        map[string]interface{} // 环境变量
	host        string                 // 执行的目标机
	starttime   time.Time              // 计算执行时间之开始时间
	hook        *Hook
	// 上下文
	msg       *Message
	logs      map[string]string // 命令执行日志
	terminial bool              // ssh 是否交互式执行
}

// common 公共初始化函数
func (r *RoleLC) Common(args *RoleArgs) error {
	// 获取name
	if name, ok := args.currentConfig["name"]; !ok {
		return errors.New("config 无 name字段")
	} else {
		r.name = name.(string)
	}

	// 获取stage
	if current_stage, ok := args.currentConfig["stage"]; !ok {
		return errors.New("config 无 stage字段")
	} else {
		if args.stage != current_stage.(string) {
			return errors.New(fmt.Sprintf("stage not equal %s %d != %s %d", args.stage, len(args.stage), current_stage, len(current_stage.(string))))
		}
	}

	r.hook = args.hook
	r.logs = make(map[string]string)
	// 上下文消息传递
	r.msg = args.msg
	r.remote_user = args.user
	r.remote_pwd = args.pwd
	r.remote_port = args.port
	r.stage = args.stage
	r.vars = args.vars

	// 设置执行主机
	r.host = args.host

	// 是否在可执行主机范围内
	isTags := false

	// 获取tags目标执行主机
	if tags, ok := args.currentConfig["tags"]; ok {
		for _, tag := range tags.([]interface{}) {
			if args.host == tag.(string) {
				isTags = true
			}
		}
	} else {
		// 没有设置tags标签，表示不限制主机执行
		isTags = true
	}

	// 每个config shell terminial 是否交互式执行
	if terminial, ok := args.currentConfig["terminial"]; ok {
		r.terminial = terminial.(bool)
	} else {
		r.terminial = false
	}

	if !isTags {
		return errors.New(fmt.Sprintf("Stage: %s Name: %s Host: %s 不在可执行主机范围内，退出！", args.stage, r.name, args.host))
	}

	return nil
}

// 初始化
func (r *RoleLC) Init(args *RoleArgs) error {
	log.Debug("Parent Object Init Func")
	return r.Common(args)
}

// 准备环节
func (r *RoleLC) Pre() {
	log.Debugf("Role module %s Pre running.", r.name)
	// 设置开始时间
	r.starttime = time.Now()
}

// 执行前
func (r *RoleLC) Before() {
	log.Debugf("Role module %s Before running.", r.name)
	log.Infof("******************************************************** TASK [%s : %s] BY %s@%s \n", r.stage, r.name, r.remote_user, r.host)
}

// 执行环节
func (r *RoleLC) Run() {
	log.Debugf("Role module %s Run running.", r.name)
}

// 执行后环节
func (r *RoleLC) After() {
	// log.Debugf("Role module %s After running.", r.name)
	log.WithFields(log.Fields{
		"Status": "After LifeCycle",
		"Role":   r.types,
		"Name":   r.name,
		"Stage":  r.stage,
		"Host":   r.host,
	}).Infof("执行完成 耗时：%v", time.Now().Sub(r.starttime))
}

// 执行判断IsHook
// default is false
func (r *RoleLC) IsHook() bool {
	return r.hook.isHook
}

// 钩子函数，思考：是否和After以及output插件冲突
func (r *RoleLC) Hooks() error {
	log.Debugf("Role module %s Hooks Args running.", r.name, r.hook.hookArgs)
	err := r.hook.hookFunc(r.hook.hookArgs)
	return err
}
