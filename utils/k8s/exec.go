package k8s

import (
	"bytes"
	"fmt"

	"github.com/astaxie/beego"

	"k8s.io/client-go/tools/remotecommand"
)

type RemoteCommand struct {
	Resource  string   `json:"resource"`  // pods
	PodName   string   `json:"podname"`   // pod名称
	Container string   `json:"container"` // pod内container名称
	Namespace string   `json:"namespace"` // 命名空间
	Command   []string `json:"command"`   // 执行命令
}

// stdout,stderr,err
func (this *RemoteCommand) Exec() (string, string, error) {
	clientset, config := GetClientSetAndConfig()

	execRequest := clientset.CoreV1().RESTClient().Post().
		Resource(this.Resource).
		Name(this.PodName).
		Namespace(this.Namespace).
		SubResource("exec").
		Param("container", this.Container).
		Param("stdin", "false").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false").
		Param("command", "/bin/bash").
		Param("command", "-c")

	for _, x := range this.Command {
		execRequest.Param("command", x)
	}
	beego.Critical(execRequest.URL().String())
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", execRequest.URL())
	if err != nil {
		return "", "", err
	}

	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
		Tty:    false,
	})

	if err != nil {
		return "", "", err
	}

	if execErr.Len() > 0 {
		return execOut.String(), fmt.Sprintf("stderr: %v", execErr.String()), nil
	}

	return execOut.String(), "", nil
}
