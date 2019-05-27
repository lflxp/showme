package k8s

import (
	// "gitlab.yc/ares/k8sApi/models"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPvList() (*apiv1.PersistentVolumeList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	return config, err
}

func GetPvByName(name string) (*apiv1.PersistentVolume, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().PersistentVolumes().Get(name, metav1.GetOptions{})
	return config, err
}

func DeletePv(name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().PersistentVolumes().Delete(name, &metav1.DeleteOptions{})
	return err
}
