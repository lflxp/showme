package roles

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 初始化copy role插件映射关系表
	addRoles(CopyType, reflect.TypeOf(CopyRole{}))
}

type CopyRole struct {
	RoleLC
	src   string // 源地址
	dest  string // 目的地址
	items []string
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (c *CopyRole) Init(args *RoleArgs) error {
	err := c.Common(args)
	if err != nil {
		return err
	}

	copyData := args.currentConfig["copy"].(map[interface{}]interface{})
	// 获取原始shell命令
	c.src = copyData["src"].(string)
	c.dest = copyData["dest"].(string)

	if item, ok := args.currentConfig["with_items"]; ok {
		for _, it := range item.([]interface{}) {
			c.items = append(c.items, fmt.Sprintf("%v", it))
		}
	}

	return nil
}

func (c *CopyRole) Run() error {
	var err error
	if c.items == nil {
		err := utils.New(c.host, c.remote_user, c.remote_pwd, c.remote_port).SftpUploadToRemote(c.src, c.dest)
		if err != nil {
			log.WithFields(log.Fields{
				"src":  c.src,
				"dest": c.dest,
				"耗时":   time.Now().Sub(c.starttime),
			}).Errorln(err.Error())
			c.logs[fmt.Sprintf("%s %s %s", c.stage, c.host, c.name)] = err.Error()
			if strings.Contains(err.Error(), "ssh:") {
				err = errors.New("ssh: handshake failed")
				goto OVER
			}
		} else {
			log.WithFields(log.Fields{
				"src":  c.src,
				"dest": c.dest,
				"耗时":   time.Now().Sub(c.starttime),
			}).Infof("success upload file %s", c.dest)
			c.logs[fmt.Sprintf("%s %s %s", c.stage, c.host, c.name)] = fmt.Sprintf("success upload file %s", c.dest)
		}
	} else {
		for _, it := range c.items {
			// 补充go template基本语法
			// 注意：只针对with_items数组类型
			src, err := utils.ApplyTemplate(c.src, map[string]interface{}{"item": it})
			if err != nil {
				log.Errorf("src %s error: %v", src, err)
				panic(err)
			}
			dest, err := utils.ApplyTemplate(c.dest, map[string]interface{}{"item": it})
			err = utils.New(c.host, c.remote_user, c.remote_pwd, c.remote_port).SftpUploadToRemote(src, dest)
			if err != nil {
				log.WithFields(log.Fields{
					"src":  src,
					"dest": dest,
					"耗时":   time.Now().Sub(c.starttime),
				}).Errorln(err.Error())
				c.logs[fmt.Sprintf("%s %s %s", c.stage, c.host, c.name)] = err.Error()
				if strings.Contains(err.Error(), "ssh:") {
					err = errors.New("ssh: handshake failed")
					goto OVER
				}
			} else {
				log.WithFields(log.Fields{
					"src":  src,
					"dest": dest,
					"耗时":   time.Now().Sub(c.starttime),
				}).Infof("success upload file %s", dest)
				c.logs[fmt.Sprintf("%s %s %s", c.stage, c.host, c.name)] = fmt.Sprintf("success upload file %s", dest)
			}
		}
	}
OVER:
	return err
}

// 处理返回日志
func (c *CopyRole) After() {
	stoptime := time.Now()
	c.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(c.starttime))
	c.msg.CallBack[fmt.Sprintf("%s-%s-%s", c.host, c.stage, c.name)] = c.logs
}
