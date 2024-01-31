package module

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/roles"
	"github.com/devopsxp/xp/utils"
	uuid "github.com/satori/go.uuid"
)

func init() {
	// 初始化shell filter插件映射关系表
	plugin.AddFilter("shell", reflect.TypeOf(ShellFilter{}))
}

// shell 命令运行filter插件
type ShellFilter struct {
	LifeCycle
	status plugin.StatusPlugin
}

func (s *ShellFilter) Process(msgs *plugin.Message) *plugin.Message {
	if s.status != plugin.Started {
		slog.Warn("Shell filter plugin is not running,filter nothing.")
		return msgs
	}

	// TODO:
	// 1. 封装config shell|copy|template等操作
	slog.Info(fmt.Sprintf("ShellFilter Filter 插件开始执行目标主机Config Playbook，并发数： %d", runtime.NumCPU()))

	// 解析yaml结果
	slog.Debug(fmt.Sprintf("解析yaml结果 Check %v\n", msgs.Data.Check))
	// 1. 解析stage步骤
	var stages []interface{}
	if sg, ok := msgs.Data.Items["stage"]; ok {
		stages = sg.([]interface{})
	}
	slog.Debug(fmt.Sprintf("Stage %v\n", stages))

	var configs []interface{}
	if cf, ok := msgs.Data.Items["config"]; ok {
		configs = cf.([]interface{})
	} else {
		slog.Error("未配置config模块，退出！")
		os.Exit(1)
	}

	// 解析include目录文件加载Yaml
	// include: /tmp/d.yaml
	// name: "name"
	for _, cc := range configs {
		rt, ok := roles.ParseRoleType(cc.(map[interface{}]interface{}))
		if ok && rt == roles.IncludeType {
			// 获取include路径
			includePath, ok := cc.(map[interface{}]interface{})["include"]
			if ok {
				slog.Info(fmt.Sprintf("匹配到 include 配置[%s] %s", cc.(map[interface{}]interface{})["name"], includePath))
				// include 配置格式为： [map[interface{}]interface{}]
				iData, err := utils.ReadYamlConfig(includePath.(string))
				if err != nil {
					slog.Error("读取include yaml文件", "错误", err.Error())
					os.Exit(1)
				}

				switch iData.(type) { //v表示b1 接口转换成Bag对象的值
				case []interface{}:
					configs = append(configs, iData.([]interface{})...)
				case map[interface{}]interface{}:
					configs = append(configs, iData.(map[interface{}]interface{}))
				default:
					slog.Warn("Include Yaml文件格式不能匹配 %v", iData)
				}
			}
		}
	}

	slog.Debug("configs", configs)

	slog.Debug(fmt.Sprintf("Config %v\n", configs))
	var (
		pipelineUuid                               string
		remote_user, remote_pwd, workdir, reponame string
		remote_port                                int
		timeout                                    int
		timeoutexit                                bool
	)

	// 设置超时时间
	if to, ok := msgs.Data.Items["timeout"]; ok {
		timeout = to.(int)
	} else {
		// 默认root用户
		timeout = 60
	}

	// 设置超时是否退出
	if toe, ok := msgs.Data.Items["timeoutexit"]; ok {
		timeoutexit = toe.(bool)
	} else {
		// 默认root用户
		timeoutexit = true
	}

	slog.Info("******************************************************** Prepare [DockerWorkspace : 镜像工作空间设置] ")
	pipelineUuid = uuid.NewV4().String()

	// docker共享空间
	if path, ok := msgs.Data.Items["workdir"]; ok {
		workdir = path.(string)
		isexist, err := utils.PathExists(workdir)
		if !isexist && err != nil {
			panic(err)
		}
		slog.Debug(fmt.Sprintf("判断docker共享目录是否存在: %s", workdir))
	} else {
		// 当没有设置workdir时，判断并创建在当前目录下
		workdir = fmt.Sprintf("%s/workspace", utils.GetCurrentDirectory())
		isexist, err := utils.PathExists(workdir)
		if !isexist && err != nil {
			panic(err)
		}
		slog.Debug(fmt.Sprintf("判断docker共享目录是否存在: %s", workdir))
	}

	slog.Info(fmt.Sprintf("准备docker 共享目录完毕: %s", workdir))

	// 如果设置了git仓库，则拉取repo并修改workdir路径
	if git, ok := msgs.Data.Items["git"]; ok {

		var branch, depth, cmd string
		// 如果设置了url则进行往下进行
		if url, ok := git.(map[string]interface{})["url"]; ok {
			slog.Info(fmt.Sprintf("******************************************************** Prepare [GitClone : %s] ", url))

			reponame = strings.Split(url.(string), "/")[len(strings.Split(url.(string), "/"))-1]
			if br, ok := git.(map[string]interface{})["branch"]; ok {
				branch = br.(string)
			}

			if dep, ok := git.(map[string]interface{})["depth"]; ok {
				depth = fmt.Sprintf("%d", dep.(int))
			}

			if branch != "" && depth != "" {
				cmd = fmt.Sprintf("git clone %s -b %s --depth %s %s/%s", url, branch, depth, workdir, pipelineUuid)
			} else if branch == "" && depth != "" {
				cmd = fmt.Sprintf("git clone %s --depth %s %s/%s", url, depth, workdir, pipelineUuid)
			} else if branch != "" && depth == "" {
				cmd = fmt.Sprintf("git clone %s -b %s %s/%s", url, branch, workdir, pipelineUuid)
			} else {
				cmd = fmt.Sprintf("git clone %s %s/%s", url, workdir, pipelineUuid)
			}

			rs, err := utils.ExecCommandString(cmd)
			if err != nil {
				slog.Error(rs)
			}

			slog.Info("success git clone", "url", rs)

			workdir = fmt.Sprintf("%s/%s", workdir, pipelineUuid)
			slog.Info(fmt.Sprintf("准备docker git clone共享目录完毕: %s", workdir))
		}
	}

	if user, ok := msgs.Data.Items["remote_user"]; ok {
		remote_user = user.(string)
	} else {
		// 默认root用户
		remote_user = "root"
	}

	if pwd, ok := msgs.Data.Items["remote_pwd"]; ok {
		remote_pwd = pwd.(string)
	} else {
		// 默认root用户
		remote_pwd = ""
	}

	if port, ok := msgs.Data.Items["remote_port"]; ok {
		remote_port = port.(int)
	} else {
		// 默认root用户
		remote_port = 22
	}

	// 全局动态变量
	var vars map[string]interface{}
	if vv, ok := msgs.Data.Items["vars"]; ok {
		vars = vv.(map[string]interface{})
	} else {
		vars = make(map[string]interface{})
	}

	rolesData := msgs.Data.Items["roles"].([]interface{})
	// 2. 根据stage进行解析
	for host, status := range msgs.Data.Check {
		if status == "failed" {
			slog.Debug(fmt.Sprintf("host %s is failed, next.\n", host))
		} else {
			for _, stage := range stages {
				if roles.IsRolesAllow(stage.(string), rolesData) {
					// 3. TODO: 解析yaml中shell的模块，然后进行匹配
					err := roles.NewShellRole(roles.NewRoleArgs(stage.(string), remote_user, remote_pwd, host, workdir, reponame, vars, configs, msgs, nil, remote_port, timeout, timeoutexit))
					if err != nil {
						slog.Debug(err.Error())
						// os.Exit(1)
						break
					}
				}
			}

			// execChan := make(chan string, runtime.NumCPU())
			// var w sync.WaitGroup
			// for _, stage := range stages {
			// w.Add(1)
			// execChan <- stage.(string)
			// go func() {
			// 	defer w.Done()
			// 	// 判断stage是否允许执行
			// 	if roles.IsRolesAllow(stage.(string), rolesData) {
			// 		// 3. TODO: 解析yaml中shell的模块，然后进行匹配
			// 		err := roles.NewShellRole(roles.NewRoleArgs(stage.(string), remote_user, host, vars, configs, msgs, nil))
			// 		if err != nil {
			// 			slog.Debug(err.Error())
			// 			os.Exit(1)
			// 		}
			// 	}
			// 	<-execChan
			// }()
			// }
			// w.Wait()
		}
	}

	return msgs
}

func (s *ShellFilter) Init(data interface{}) {
	s.name = "Shell Filter"
	s.status = plugin.Started
}
