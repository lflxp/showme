package k8s

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetEvenstList() (*apiv1.EventList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Events("").List(metav1.ListOptions{})
	return config, err
}

func GetEventsByNamespace(namespace string) (*apiv1.EventList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Events(namespace).List(metav1.ListOptions{})
	return config, err
}

func GetEventsByName(namespace, name string) (*apiv1.Event, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Events(namespace).Get(name, metav1.GetOptions{})
	return config, err
}

func DeleteEvents(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().Events(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}
