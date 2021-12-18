package roles

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 初始化copy role插件映射关系表
	addRoles(UserType, reflect.TypeOf(UserRole{}))
}

type UserRole struct {
	RoleLC
	name     string // 用户名
	password string
	delete   bool   // 删除用户
	lock     bool   // 锁住用户无权更改密码
	unlock   bool   // 解除锁定
	force    bool   // 强制操作
	maximum  string // 两次密码修正的最大天数，后面接数字；仅能root权限操作；
	minimum  string // 两次密码修改的最小天数，后面接数字，仅能root权限操作；
	warning  string // 在距多少天提醒用户修改密码；仅能root权限操作
	inactive string // 在密码过期后多少天，用户被禁掉，仅能以root操作；
	status   bool   // 查询用户的密码状态，仅能root用户操作；
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (s *UserRole) Init(args *RoleArgs) error {
	err := s.Common(args)
	if err != nil {
		return err
	}

	copyData := args.currentConfig["user"].(map[interface{}]interface{})
	// 获取原始shell命令
	if name, ok := copyData["name"].(string); ok {
		s.name = name
	} else {
		return errors.New(fmt.Sprintf("未提供用户名 %s", s.name))
	}
	if pwd, ok := copyData["password"].(string); ok {
		s.password = pwd
	}
	if maximum, ok := copyData["maximum"].(string); ok {
		s.maximum = maximum
	}
	if minimum, ok := copyData["minimum"].(string); ok {
		s.minimum = minimum
	}
	if warning, ok := copyData["warning"].(string); ok {
		s.warning = warning
	}
	if inactive, ok := copyData["inactive"].(string); ok {
		s.inactive = inactive
	}

	if delete, ok := copyData["delete"].(bool); ok {
		s.delete = delete
	}
	if lock, ok := copyData["lock"].(bool); ok {
		s.lock = lock
	}
	if unlock, ok := copyData["unlock"].(bool); ok {
		s.unlock = unlock
	}
	if force, ok := copyData["force"].(bool); ok {
		s.force = force
	}
	if status, ok := copyData["status"].(bool); ok {
		s.status = status
	}

	log.Debugf("User %v", s)

	return nil
}

func (s *UserRole) setDo(cmd string) error {
	var (
		err error
		rs  string
	)
	if s.terminial {
		err = utils.New(s.host, s.remote_user, s.remote_pwd, s.remote_port).RunTerminal(cmd, os.Stdout, os.Stderr)
	} else {
		_, err = utils.New(s.host, s.remote_user, s.remote_pwd, s.remote_port).Run(cmd)
	}

	if err != nil {
		log.WithFields(log.Fields{"耗时": time.Now().Sub(s.starttime)}).Errorln(fmt.Sprintf("[Item: %s] => %s <=> %s", cmd, rs, err.Error()))
	} else {
		log.WithFields(log.Fields{"耗时": time.Now().Sub(s.starttime)}).Infoln(fmt.Sprintf("[Item: %s] => %s ", cmd, rs))
	}
	return err
}

func (s *UserRole) Run() error {
	// 检查命令是否正确
	if s.name == "" {
		return errors.New("未提供用户名")
	}

	if s.password != "" {
		err := s.setDo(fmt.Sprintf("echo \"%s:%s\"|sudo chpasswd", s.name, s.password))
		if err != nil {
			return err
		}
	}

	if s.delete {
		err := s.setDo(fmt.Sprintf("passwd -d %s", s.name))
		if err != nil {
			return err
		}
	}

	if s.lock {
		err := s.setDo(fmt.Sprintf("passwd -l %s", s.name))
		if err != nil {
			return err
		}
	}

	if s.unlock {
		err := s.setDo(fmt.Sprintf("passwd -u %s", s.name))
		if err != nil {
			return err
		}
	}

	if s.inactive != "" {
		err := s.setDo(fmt.Sprintf("passwd -i %s %s", s.inactive, s.name))
		if err != nil {
			return err
		}
	}

	if s.warning != "" {
		err := s.setDo(fmt.Sprintf("passwd -w %s %s", s.warning, s.name))
		if err != nil {
			return err
		}
	}

	if s.maximum != "" {
		err := s.setDo(fmt.Sprintf("passwd -x %s %s", s.maximum, s.name))
		if err != nil {
			return err
		}
	}

	if s.minimum != "" {
		err := s.setDo(fmt.Sprintf("passwd -n %s %s", s.minimum, s.name))
		if err != nil {
			return err
		}
	}

	if s.status {
		err := s.setDo(fmt.Sprintf("passwd -S %s ", s.name))
		if err != nil {
			return err
		}
	}

	return nil
}

// 处理返回日志
func (s *UserRole) After() {
	stoptime := time.Now()
	s.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(s.starttime))
	s.msg.CallBack[fmt.Sprintf("%s-%s-%s", s.host, s.stage, s.name)] = s.logs
}
