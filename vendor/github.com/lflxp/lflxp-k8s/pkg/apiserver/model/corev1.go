package model

import (
	"context"
	"io"

	"github.com/lflxp/tools/sdk/clientgo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CoreV1 struct {
	Namespace     string `json:"namespace"`
	ContainerName string `json:"container_name"`
	PodName       string `json:"pod_name"`
	TailLines     int64  `json:"tail_lines"`
	Previous      bool   `json:"previous"`
	Fellow        bool   `json:"fellow"`
}

func (g *CoreV1) Namespaces() (*v1.NamespaceList, error) {
	client := clientgo.InitClient()
	data, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	return data, err
}

func (g *CoreV1) Nodes() (*v1.NodeList, error) {
	client := clientgo.InitClient()
	data, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	return data, err
}

func (g *CoreV1) RestClient() *kubernetes.Clientset {
	return clientgo.InitClient()
}

// getlogs
/**
    req := client.RESTClient.Get().
        Namespace(namespace).
        Name(podID).
        Resource("pods").
        SubResource("log").
        Param("follow", strconv.FormatBool(logOptions.Follow)).
        Param("container", logOptions.Container).
        Param("previous", strconv.FormatBool(logOptions.Previous)).
        Param("timestamps", strconv.FormatBool(logOptions.Timestamps))

    if logOptions.SinceSeconds != nil {
        req.Param("sinceSeconds", strconv.FormatInt(*logOptions.SinceSeconds, 10))
    }
    if logOptions.SinceTime != nil {
        req.Param("sinceTime", logOptions.SinceTime.Format(time.RFC3339))
    }
    if logOptions.LimitBytes != nil {
        req.Param("limitBytes", strconv.FormatInt(*logOptions.LimitBytes, 10))
    }
    if logOptions.TailLines != nil {
        req.Param("tailLines", strconv.FormatInt(*logOptions.TailLines, 10))
    }
    readCloser, err := req.Stream()
    if err != nil {
        return err
    }

    defer readCloser.Close()
    _, err = io.Copy(out, readCloser)
    return err
**/
// https://blog.51cto.com/u_13622854/5320810
// https://stackoverflow.com/questions/32983228/kubernetes-go-client-api-for-log-of-a-particular-pod
func (g *CoreV1) GetLogs() (io.ReadCloser, error) {
	client := clientgo.InitClient()
	logOpt := &v1.PodLogOptions{
		Container: g.ContainerName,
		Follow:    g.Fellow,
		TailLines: &g.TailLines,
		Previous:  g.Previous,
	}
	req := client.CoreV1().Pods(g.Namespace).GetLogs(g.PodName, logOpt)

	return req.Stream(context.TODO())
}

// func (g *CoreV1) Exec(cmd []string, ptyHandler PtyHandler, namespace, podName, containerName string) error {
// 	defer func() {
// 		ptyHandler.Done()
// 	}()

// 	client := clientgo.InitClient()
// 	req := client.CoreV1().RESTClient().Post().
// 		Resource("pods").Name(podName).
// 		Namespace(namespace).SubResource("exec").
// 		VersionedParams(&v1.PodExecOptions{
// 			Container: containerName,
// 			Command:   cmd,
// 			Stdin:     !(ptyHandler.Stdin() == nil),
// 			Stdout:    !(ptyHandler.Stdout() == nil),
// 			Stderr:    !(ptyHandler.Stderr() == nil),
// 			TTY:       ptyHandler.Tty(),
// 		}, scheme.ParameterCodec)

// 	exec, err := remotecommand.NewSPDYExecutor(clientgo.RestConfig(), "POST", req.URL())
// 	if err != nil {
// 		return err
// 	}

// 	err = exec.Stream(remotecommand.StreamOptions{
// 		Stdin:             ptyHandler.Stdin(),
// 		Stdout:            ptyHandler.Stdout(),
// 		Stderr:            ptyHandler.Stderr(),
// 		TerminalSizeQueue: ptyHandler,
// 		Tty:               ptyHandler.Tty(),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
