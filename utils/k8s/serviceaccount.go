package k8s

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetServiceAccountList() (*apiv1.ServiceAccountList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().ServiceAccounts("").List(metav1.ListOptions{})
	return config, err
}

func GetServiceAccountByName(namespace, name string) (*apiv1.ServiceAccount, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
	return config, err
}

func DeleteServiceAccount(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().ServiceAccounts(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}
