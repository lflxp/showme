package roles

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
)

// 执行过程生命周期接口
// Role LifeCycle
type RolePlugin interface {
	// 初始化对象
	Init(*RoleArgs) error

	// 环境准备
	Pre()

	// 执行前
	Before() error

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
	name  string `json:"name"` // 名称
	types string `json:"types"`
	// 通用字段
	stage       string                 `json:"stage"`
	remote_user string                 `json:"remote_user"` // 执行用户
	remote_pwd  string                 `json:"remote_pwd"`  // 执行用户密码（非必须）
	remote_port int                    `json:"remote_port"` // ssh端口
	vars        map[string]interface{} `json:"vars"`        // 环境变量
	host        string                 `json:"host"`        // 执行的目标机
	starttime   time.Time              `json:"starttime"`   // 计算执行时间之开始时间
	hook        *Hook                  `json:"hook"`
	// 上下文
	msg       *Message          `json:"msg"`
	logs      map[string]string `json:"logs"`      // 命令执行日志
	terminial bool              `json:"terminial"` // ssh 是否交互式执行
	args      *RoleArgs         `json:"args"`
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

	r.args = args
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
	slog.Debug("Parent Object Init Func")
	return r.Common(args)
}

// 准备环节
func (r *RoleLC) Pre() {
	slog.Debug(fmt.Sprintf("Role module %s Pre running.", r.name))
	// 设置开始时间
	r.starttime = time.Now()
}

// 执行前的条件判断
func (r *RoleLC) Before() error {
	slog.Debug(fmt.Sprintf("Role module %s Before running.", r.name))
	slog.Info(fmt.Sprintf("******************************************************** TASK [%s : %s] BY %s@%s \n", r.stage, r.name, r.remote_user, r.host))
	// when条件判断 @key shell命令 @value 命令结果
	if when, ok := r.args.currentConfig["when"]; ok {
		// 如果when条件判断存在，则执行@key命令
		whenData := when.(map[interface{}]interface{})
		key := whenData["key"].(string)
		value := whenData["value"].(string)

		rs, err := utils.New(r.host, r.remote_user, r.remote_pwd, r.remote_port).Run(key)
		if err != nil {
			return errors.New(fmt.Sprintf("%s %s when 条件 [%s] 错误: %v", r.host, r.remote_user, key, err))
		}

		if rs != value {
			return errors.New(fmt.Sprintf("%s %s when 条件 [%s|%s] %s 不满足: %v", r.host, r.remote_user, key, value, err))
		}
	}
	return nil
}

// 执行环节
func (r *RoleLC) Run() {
	slog.Debug(fmt.Sprintf("Role module %s Run running.", r.name))
}

// 执行后环节
func (r *RoleLC) After() {
	// slog.Debug(fmt.Sprintf("Role module %s After running.", r.name)
	slog.Info(fmt.Sprintf("执行完成 耗时：%v", time.Now().Sub(r.starttime)))
}

// 执行判断IsHook
// default is false
func (r *RoleLC) IsHook() bool {
	return r.hook.isHook
}

// 钩子函数，思考：是否和After以及output插件冲突
func (r *RoleLC) Hooks() error {
	slog.Debug(fmt.Sprintf("Role module %s Hooks Args running.", r.name, r.hook.hookArgs))
	err := r.hook.hookFunc(r.hook.hookArgs)
	return err
}
