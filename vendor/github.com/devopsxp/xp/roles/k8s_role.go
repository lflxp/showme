package roles

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/devopsxp/xp/pkg/k8s"
	"github.com/devopsxp/xp/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	// 初始化k8s role插件映射关系表
	addRoles(K8sType, reflect.TypeOf(K8sRole{}))
}

// 对标pod内containerd信息
type K8sBasic struct {
	command []string          // 执行命令
	env     map[string]string // 环境变量
	args    []string          // 命令参数
	name    string
	image   string
}

type GitRepo struct {
	url  string
	user string
	pwd  string
}

type K8sRole struct {
	RoleLC
	k8s          []K8sBasic // 容器pod组成
	workspace    string     // 共享目录空间
	workspaceRaw string     // 未设置git之前
	repo         GitRepo    // git代码
	name         string
	namespace    string
}

func (k *K8sRole) Init(args *RoleArgs) error {
	err := k.Common(args)
	if err != nil {
		return err
	}

	// 设置hook做自动删除pod处理
	k.hook.isHook = true

	k.repo.url = args.reponame
	// 获取git信息
	if g, ok := args.currentConfig["git"]; ok {
		if url, ok := g.(map[interface{}]interface{})["url"]; ok {
			k.repo.url = url.(string)
		} else {
			k.repo.url = "https://github.com/lflxp/helloworld.git"
		}

		if user, ok := g.(map[interface{}]interface{})["user"]; ok {
			k.repo.user = user.(string)
		}

		if pwd, ok := g.(map[interface{}]interface{})["pwd"]; ok {
			k.repo.pwd = pwd.(string)
		}
	} else {
		k.repo.url = "https://github.com/lflxp/helloworld.git"
	}

	// 获取name
	if na, ok := args.currentConfig["name"]; ok {
		k.name = na.(string)
	} else {
		k.name = fmt.Sprintf("unknow-%s", time.Now().Format("200601021504"))
	}

	// 获取workspace
	if ws, ok := args.currentConfig["workspace"]; ok {
		tmpUrl := strings.Split(strings.ReplaceAll(k.repo.url, ".git", ""), "/")
		k.workspace = fmt.Sprintf("%s/%s", ws.(string), tmpUrl[len(tmpUrl)-1])
		k.workspaceRaw = ws.(string)
	} else {
		k.workspace = "/workspace"
	}

	// 获取namespace
	if ns, ok := args.currentConfig["namespace"]; ok {
		k.namespace = ns.(string)
	} else {
		k.namespace = "default"
	}

	// TODO: 解析k8s字段
	tmp := args.currentConfig["k8s"].([]interface{})
	for _, x := range tmp {
		k8sbasicData := K8sBasic{
			command: []string{},
			env:     map[string]string{},
			args:    []string{},
		}

		if n, ok := x.(map[interface{}]interface{})["name"]; ok {
			if n.(string) == "" {
				k8sbasicData.name = utils.GetRandomString(32)
			} else {
				k8sbasicData.name = n.(string)
			}
		}

		if im, ok := x.(map[interface{}]interface{})["image"]; ok {
			if im.(string) == "" {
				return errors.New("image is none")
			} else {
				k8sbasicData.image = im.(string)
			}
		}

		if sc, ok := x.(map[interface{}]interface{})["command"]; ok {
			for _, it := range sc.([]interface{}) {
				k8sbasicData.command = append(k8sbasicData.command, it.(string))
			}
		}

		if args, ok := x.(map[interface{}]interface{})["args"]; ok {
			for _, arg := range args.([]interface{}) {
				k8sbasicData.args = append(k8sbasicData.args, arg.(string))
			}
		}

		if e, ok := x.(map[interface{}]interface{})["env"]; ok {
			for k, v := range e.(map[interface{}]interface{}) {
				k8sbasicData.env[k.(string)] = v.(string)
			}
		}
		k.k8s = append(k.k8s, k8sbasicData)
	}

	return nil
}

func (k *K8sRole) Run() error {
	// 组装pod
	pod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:         k.name,
			GenerateName: k.name,
			Namespace:    k.namespace,
		},
		Spec: apiv1.PodSpec{
			InitContainers: []apiv1.Container{
				apiv1.Container{
					Name:       "git-clone",
					Image:      "docker:git",
					WorkingDir: k.workspaceRaw,
					Command: []string{
						"sh",
						"-c",
						"git clone " + k.repo.url,
					},
					VolumeMounts: []apiv1.VolumeMount{
						apiv1.VolumeMount{
							Name:      "workdir",
							MountPath: "/workspace",
						},
					},
				},
			},
			Volumes: []apiv1.Volume{
				apiv1.Volume{
					Name: "workdir",
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	containers := []apiv1.Container{}
	// TODO: env from k8srole => EnvVar
	// DONE: git clone at initContainers
	// DONE: set share folder to every container
	for _, cc := range k.k8s {
		tmp := apiv1.Container{
			Name:       cc.name,
			Image:      cc.image,
			Command:    cc.command,
			Args:       cc.args,
			WorkingDir: k.workspace,
			VolumeMounts: []apiv1.VolumeMount{
				apiv1.VolumeMount{
					Name:      "workdir",
					MountPath: "/workspace",
				},
			},
		}
		containers = append(containers, tmp)
	}

	pod.Spec.Containers = containers

	slog.Debug(fmt.Sprintf("pod name: %s %v", k.name, pod))

	podinfo, err := k8s.CreatePod(pod)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	info, err := json.MarshalIndent(podinfo, "", "\t")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Pod ", "YAML", info)
	return err
}

// 处理返回日志
func (k *K8sRole) After() {
	stoptime := time.Now()
	k.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(k.starttime))
	k.msg.CallBack[fmt.Sprintf("%s-%s-%s", k.host, k.stage, k.name)] = k.logs
}

// hook钩子 进行pod删除
func (k *K8sRole) Hooks() error {
	slog.Debug(fmt.Sprintf("K8s Role module on %s => %s", k.namespace, k.name))
	// err := k8s.DeletePod(k.namespace, k.name)
	// log.WithFields(log.Fields{
	// 	"耗时": time.Now().Sub(k.starttime),
	// }).Infof("******************************************************** K8S Role Hook: 删除Pod: %s Namespace: %s  [%s By %s@%s ], Result: %v \n", k.name, k.namespace, k.stage, k.remote_user, k.host, err)
	k.msg.Tmp["namespace"] = k.namespace
	k.msg.Tmp["name"] = k.name
	return nil
}
