package roles

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 初始化template role插件映射关系表
	addRoles(TemplateType, reflect.TypeOf(TemplateRole{}))
}

type TemplateRole struct {
	RoleLC
	src  string // 源地址
	dest string // 目的地址
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (r *TemplateRole) Init(args *RoleArgs) error {
	err := r.Common(args)
	if err != nil {
		return err
	}

	copyData := args.currentConfig["template"].(map[interface{}]interface{})
	// 获取原始shell命令
	r.src = copyData["src"].(string)
	r.dest = copyData["dest"].(string)

	return nil
}

// 操作流程：1. 获取vars和template 2. 解析 3. 上传
func (r *TemplateRole) Run() error {
	// 读取j2文件
	templateFile, err := ioutil.ReadFile(r.src)
	if err != nil {
		return err
	}

	destFile, err := utils.ApplyTemplate(string(templateFile), r.vars)
	if err != nil {
		return err
	}

	log.Debugf("template is %s", destFile)

	err = utils.New(r.host, r.remote_user, r.remote_pwd, r.remote_port).SftpUploadTemplateString(destFile, r.dest)
	if err != nil {
		log.WithFields(log.Fields{
			"template": r.src,
			"dest":     r.dest,
			"耗时":       time.Now().Sub(r.starttime),
		}).Errorln(err.Error())
		r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = err.Error()
		if strings.Contains(err.Error(), "ssh:") {
			err = errors.New("ssh: handshake failed")
			return err
		}
	} else {
		log.WithFields(log.Fields{
			"template": r.src,
			"dest":     r.dest,
			"耗时":       time.Now().Sub(r.starttime),
		}).Infof("模板上传成功 %s", r.dest)
		r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = fmt.Sprintf("模板上传成功 %s", r.dest)
	}

	return nil
}

// 处理返回日志
func (r *TemplateRole) After() {
	stoptime := time.Now()
	r.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(r.starttime))
	r.msg.CallBack[fmt.Sprintf("%s-%s-%s", r.host, r.stage, r.name)] = r.logs
}
