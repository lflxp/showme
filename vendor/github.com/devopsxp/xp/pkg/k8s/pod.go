package k8s

import (
	"bytes"
	"context"
	"io"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodList() (*apiv1.PodList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	return config, err
}

func GetPodListByNamespace(namespace string) (*apiv1.PodList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	return config, err
}

func GetPodListByLabels(namespace, label string) (*apiv1.PodList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	listoptions := metav1.ListOptions{
		LabelSelector: label,
	}

	config, err := cli.CoreV1().Pods(namespace).List(context.TODO(), listoptions)
	return config, err
}

// 创建pod
func CreatePod(pod *apiv1.Pod) (*apiv1.Pod, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	log.Debugf("create pod  %v\n", pod)
	podCallBack, err := cli.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(context.Background(), pod, metav1.CreateOptions{})
	return podCallBack, err
}

func GetPodByName(namespace, name string) (*apiv1.Pod, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return config, err
}

func DeletePod(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func GetPodLogByPodId(namespace, podid string) (string, error) {
	var result string
	f := func(s int64) *int64 {
		return &s
	}

	cli, err := GetClientSet()
	if err != nil {
		return result, err
	}

	options := apiv1.PodLogOptions{
		Follow:    false,
		TailLines: f(1000),
	}
	req := cli.CoreV1().Pods(namespace).GetLogs(podid, &options)
	readCloser, err := req.Stream(context.TODO())
	if err != nil {
		return result, err
	}
	defer readCloser.Close()
	var out bytes.Buffer
	_, err = io.Copy(&out, readCloser)
	if err != nil {
		return result, err
	}
	result = out.String()
	return result, nil
}

func GetPodLogByPodIdByNum(namespace, podid string, num int64) (string, error) {
	var result string
	f := func(s int64) *int64 {
		return &s
	}

	cli, err := GetClientSet()
	if err != nil {
		return result, err
	}

	options := apiv1.PodLogOptions{
		Follow:       false,
		SinceSeconds: f(num),
	}
	req := cli.CoreV1().Pods(namespace).GetLogs(podid, &options)
	readCloser, err := req.Stream(context.TODO())
	if err != nil {
		return result, err
	}
	defer readCloser.Close()
	var out bytes.Buffer
	_, err = io.Copy(&out, readCloser)
	if err != nil {
		return result, err
	}
	result = out.String()
	return result, nil
}

func GetPodLogByPodIdAll(namespace, podid string) (string, error) {
	var result string

	cli, err := GetClientSet()
	if err != nil {
		return result, err
	}

	options := apiv1.PodLogOptions{
		Follow: false,
	}
	req := cli.CoreV1().Pods(namespace).GetLogs(podid, &options)
	readCloser, err := req.Stream(context.TODO())
	if err != nil {
		return result, err
	}
	defer readCloser.Close()
	var out bytes.Buffer
	_, err = io.Copy(&out, readCloser)
	if err != nil {
		return result, err
	}
	result = out.String()
	return result, nil
}
