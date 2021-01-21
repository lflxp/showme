package roles

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 初始化copy role插件映射关系表
	addRoles(ScriptType, reflect.TypeOf(ScriptRole{}))
}

type ScriptRole struct {
	RoleLC
	script string // 脚本地址
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (c *ScriptRole) Init(args *RoleArgs) error {
	err := c.Common(args)
	if err != nil {
		return err
	}

	data := args.currentConfig["script"].(string)
	// 获取原始shell命令
	c.script = data

	return nil
}

// 1. copy file
// 2. exec script
// 3. delete script
func (c *ScriptRole) Run() error {
	dest := fmt.Sprintf("/tmp/%s", utils.GetRandomSalt())
	err := utils.New(c.host, c.remote_user, c.remote_pwd, c.remote_port).SftpUploadToRemote(c.script, dest)
	if err != nil {
		log.WithFields(log.Fields{
			"src": c.script,
			"耗时":  time.Now().Sub(c.starttime),
		}).Errorln(err.Error())
		c.logs[fmt.Sprintf("%s %s %s", c.stage, c.host, c.name)] = err.Error()
		if strings.Contains(err.Error(), "ssh:") {
			return errors.New("ssh: handshake failed")
		}
		return err
	} else {
		cmd := fmt.Sprintf("sh %s", dest)
		if c.terminial {
			err = utils.New(c.host, c.remote_user, c.remote_pwd, c.remote_port).RunTerminal(cmd, os.Stdout, os.Stderr)
			rs := fmt.Sprintf("%s over", c.script)
			if err != nil {
				log.WithFields(log.Fields{"耗时": time.Now().Sub(c.starttime)}).Errorln(fmt.Sprintf("[Item: %s] => %s", c.script, err.Error()))
				return err
			} else {
				log.WithFields(log.Fields{"耗时": time.Now().Sub(c.starttime)}).Info(fmt.Sprintf("[Item: %s] => %s", c.script, rs))
			}
		} else {
			rs, err := utils.New(c.host, c.remote_user, c.remote_pwd, c.remote_port).Run(cmd)
			if err != nil {
				log.WithFields(log.Fields{"耗时": time.Now().Sub(c.starttime)}).Errorln(fmt.Sprintf("[Item: %s] => %s", c.script, err.Error()))
				return err
			} else {
				log.WithFields(log.Fields{"耗时": time.Now().Sub(c.starttime)}).Info(fmt.Sprintf("[Item: %s] => %s", c.script, rs))
			}
		}

	}
	return nil
}

// 处理返回日志
func (c *ScriptRole) After() {
	stoptime := time.Now()
	c.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(c.starttime))
	c.msg.CallBack[fmt.Sprintf("%s-%s-%s", c.host, c.stage, c.name)] = c.logs
}
