package k8s

import (
	// "gitlab.yc/ares/k8sApi/models"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetStorageClassList() (*v1.StorageClassList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.StorageV1().StorageClasses().List(metav1.ListOptions{})
	return config, err
}

func GetStorageClassByName(name string) (*v1.StorageClass, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.StorageV1().StorageClasses().Get(name, metav1.GetOptions{})
	return config, err
}

func DeleteStorageClass(name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.StorageV1().StorageClasses().Delete(name, &metav1.DeleteOptions{})
	return err
}
