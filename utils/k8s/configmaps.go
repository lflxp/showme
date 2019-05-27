package k8s

import (
	"github.com/lflxp/showme/utils/k8s/models"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigmapsList() (*apiv1.ConfigMapList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().ConfigMaps("").List(metav1.ListOptions{})
	return config, err
}

func GetConfigmapsListByLabels(namespace, labels string) (*apiv1.ConfigMapList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	listoptions := metav1.ListOptions{
		LabelSelector: labels,
	}

	config, err := cli.CoreV1().ConfigMaps(namespace).List(listoptions)
	return config, err
}

func GetConfigmapsListByOnlyLabels(labels string) (*apiv1.ConfigMapList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	listoptions := metav1.ListOptions{
		LabelSelector: labels,
	}

	config, err := cli.CoreV1().ConfigMaps("").List(listoptions)
	return config, err
}

func GetConfigmapsByName(namespace, name string) (*apiv1.ConfigMap, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	return config, err
}

func DeleteConfigmaps(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().ConfigMaps(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}

func CreateConfigmap(data models.ConfigmapsJson) (*apiv1.ConfigMap, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	if data.Namespace == "" {
		data.Namespace = "default"
	}

	// 配置Configap
	cm := &apiv1.ConfigMap{
		Data: data.Data,
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
		},
	}

	if len(data.Labels) > 0 {
		cm.ObjectMeta.Labels = data.Labels
	}

	if len(data.Annotations) > 0 {
		cm.ObjectMeta.Annotations = data.Annotations
	}

	info, err := cli.CoreV1().ConfigMaps(data.Namespace).Create(cm)
	return info, err
}

func UpdateConfigmap(data models.ConfigmapsJson) (*apiv1.ConfigMap, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	if data.Namespace == "" {
		data.Namespace = "default"
	}

	// 配置Configap
	cm := &apiv1.ConfigMap{
		Data: data.Data,
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
		},
	}

	if len(data.Labels) > 0 {
		cm.ObjectMeta.Labels = data.Labels
	}

	if len(data.Annotations) > 0 {
		cm.ObjectMeta.Annotations = data.Annotations
	}

	info, err := cli.CoreV1().ConfigMaps(data.Namespace).Update(cm)
	return info, err
}
