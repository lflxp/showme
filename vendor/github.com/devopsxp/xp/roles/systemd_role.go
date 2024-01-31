package roles

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/devopsxp/xp/utils"
)

func init() {
	// 初始化copy role插件映射关系表
	addRoles(SystemdType, reflect.TypeOf(SystemdRole{}))
}

type SystemdRole struct {
	RoleLC
	service      string // 服务名
	state        string // 操作行为
	daemonReload bool   // 是否执行daemonReload
	enabled      bool   // 是否执行enabled
	masked       bool
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (s *SystemdRole) Init(args *RoleArgs) error {
	err := s.Common(args)
	if err != nil {
		return err
	}

	copyData := args.currentConfig["systemd"].(map[interface{}]interface{})
	// 获取原始shell命令
	if service, ok := copyData["name"].(string); ok {
		s.service = service
	} else {
		return errors.New(fmt.Sprintf("服务名未提供 %s", s.name))
	}
	if state, ok := copyData["state"].(string); ok {
		s.state = state
	}
	if daemon, ok := copyData["daemonReload"].(bool); ok {
		s.daemonReload = daemon
	}
	if enable, ok := copyData["enabled"].(bool); ok {
		s.enabled = enable
	}
	if masked, ok := copyData["masked"].(bool); ok {
		s.masked = masked
	}

	slog.Debug(fmt.Sprintf("Systemd %v", s))

	return nil
}

func (s *SystemdRole) setDaemonreload() error {
	var err error
	if s.terminial {
		err = utils.New(s.host, s.remote_user, s.remote_pwd, s.remote_port).RunTerminal("sudo systemctl daemon-reload", os.Stdout, os.Stderr)
	} else {
		_, err = utils.New(s.host, s.remote_user, s.remote_pwd, s.remote_port).Run("sudo systemctl daemon-reload")
	}
	return err
}

// 是否enable service
func (s *SystemdRole) setEnable() error {
	var (
		err error
		cmd string
	)

	if s.enabled {
		cmd = fmt.Sprintf("sudo systemctl enable %s", s.service)
	} else {
		cmd = fmt.Sprintf("sudo systemctl disable %s", s.service)
	}

	if s.terminial {
		err = utils.New(s.host, s.remote_user, s.remote_pwd, s.remote_port).RunTerminal(cmd, os.Stdout, os.Stderr)
	} else {
		_, err = utils.New(s.host, s.remote_user, s.remote_pwd, s.remote_port).Run(cmd)
	}
	return err
}

// 修改service状态
func (s *SystemdRole) setState() error {
	var (
		rs  string
		err error
	)
	switch s.state {
	case "start", "stop", "restart", "reload", "status":
		cmd := fmt.Sprintf("sudo systemctl %s %s", s.state, s.service)
		if s.terminial {
			err = utils.New(s.host, s.remote_user, s.remote_pwd, s.remote_port).RunTerminal(cmd, os.Stdout, os.Stderr)
		} else {
			rs, err = utils.New(s.host, s.remote_user, s.remote_pwd, s.remote_port).Run(cmd)
		}
		if err != nil {
			// log.WithFields(log.Fields{"耗时": time.Now().Sub(s.starttime)}).Errorln(fmt.Sprintf("[Item: %s] => %s <=> %s", cmd, rs, err.Error()))
			slog.Error(err.Error())
		} else {
			// log.WithFields(log.Fields{"耗时": time.Now().Sub(s.starttime)}).Infoln(fmt.Sprintf("[Item: %s] => %s ", cmd, rs))
			slog.Info(fmt.Sprintf("[Item: %s] => %s ", cmd, rs))
		}
	default:
		err = errors.New(fmt.Sprintf("没有匹配到正确的systemRole %s", s.state))
	}
	return err
}

func (s *SystemdRole) Run() error {
	// 检查命令是否正确
	if s.service == "" {
		return errors.New("未提供服务名")
	}

	// daemonReload适配
	if s.daemonReload {
		err := s.setDaemonreload()
		if err != nil {
			return err
		}
	}

	// enabled适配
	if s.enabled {
		err := s.setEnable()
		if err != nil {
			return err
		}
	}

	// state适配
	if s.state != "" {
		err := s.setState()
		if err != nil {
			return err
		}
	}

	return nil
}

// 处理返回日志
func (s *SystemdRole) After() {
	stoptime := time.Now()
	s.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(s.starttime))
	s.msg.CallBack[fmt.Sprintf("%s-%s-%s", s.host, s.stage, s.name)] = s.logs
}
