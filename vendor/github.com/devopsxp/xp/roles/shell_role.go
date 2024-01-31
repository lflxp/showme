package roles

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/devopsxp/xp/utils"
)

func init() {
	// 初始化shell role插件映射关系表
	addRoles(ShellType, reflect.TypeOf(ShellRole{}))
}

type ShellRole struct {
	RoleLC
	shell string   // 原生命令
	items []string // 多命令集合
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (r *ShellRole) Init(args *RoleArgs) error {
	err := r.Common(args)
	if err != nil {
		return err
	}

	// 获取原始shell命令
	r.shell = args.currentConfig["shell"].(string)

	// 获取name
	r.name = args.currentConfig["name"].(string)

	// 获取with_items迭代
	if item, ok := args.currentConfig["with_items"]; ok {
		for _, it := range item.([]interface{}) {
			r.items = append(r.items, it.(string))
		}
	}

	return nil
}

// 执行
func (r *ShellRole) Run() error {
	var (
		err error
		rs  string
	)
	if r.items == nil {
		cmd := fmt.Sprintf("bash -c \"%s\"", r.shell)
		if r.terminial {
			err = utils.New(r.host, r.remote_user, r.remote_pwd, r.remote_port).RunTerminal(cmd, os.Stdout, os.Stderr)
			rs = fmt.Sprintf("%s over", r.shell)
		} else {
			rs, err = utils.New(r.host, r.remote_user, r.remote_pwd, r.remote_port).Run(cmd)
		}

		if err != nil {
			// log.WithFields(log.Fields{"耗时": time.Now().Sub(r.starttime)}).Errorln(fmt.Sprintf("[Item: %s] => %s", r.shell, err.Error()))
			r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = err.Error()
			if strings.Contains(err.Error(), "ssh:") {
				err = errors.New("ssh: handshake failed")
				// goto OVER
			}
			// return errors.New(fmt.Sprintf("%s | %s | %s | %s => %s", r.host, r.stage, r.name, cmd, err.Error()))
		} else {
			// log.WithFields(log.Fields{"耗时": time.Now().Sub(r.starttime)}).Info(fmt.Sprintf("[Item: %s] => %s", r.shell, rs))
			r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = rs
		}
	} else {
		for _, it := range r.items {
			// 补充go template基本语法
			// 注意：只针对with_items数组类型
			cmd, err := utils.ApplyTemplate(r.shell, map[string]interface{}{"item": it})
			if err != nil {
				// slog.Errorf("cmd %s error: %v", cmd, err)
				panic(err)
			}
			slog.Debug(fmt.Sprintf("cmd is %s", cmd))

			if r.terminial {
				err = utils.New(r.host, r.remote_user, r.remote_pwd, r.remote_port).RunTerminal(cmd, os.Stdout, os.Stderr)
				rs = fmt.Sprintf("%s over", r.shell)
			} else {
				rs, err = utils.New(r.host, r.remote_user, r.remote_pwd, r.remote_port).Run(cmd)
			}

			if err != nil {
				r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = err.Error()
				if strings.Contains(err.Error(), "ssh:") {
					err = errors.New("ssh: handshake failed")
					// goto OVER
				} else {
					// log.WithFields(log.Fields{"耗时": time.Now().Sub(r.starttime)}).Errorln(fmt.Sprintf("[序号: %d Item: %s] => %s", n, cmd, err.Error()))
				}
				// return errors.New(fmt.Sprintf("%s | %s | %s | %s => %s %s", r.host, r.stage, r.name, cmd, rs, err.Error()))
				return err
			} else {
				// log.WithFields(log.Fields{"耗时": time.Now().Sub(r.starttime)}).Info(fmt.Sprintf("[序号: %d Item: %s] => %s", n, cmd, rs))
				r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = rs
			}
		}
	}
	// OVER:
	return err
}

// 处理返回日志
func (r *ShellRole) After() {
	stoptime := time.Now()
	r.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(r.starttime))
	r.msg.CallBack[fmt.Sprintf("%s-%s-%s", r.host, r.stage, r.name)] = r.logs
}

func testhook(a, b string) error {
	slog.Info(fmt.Sprintf("%s %s test hook send"))
	return nil
}
