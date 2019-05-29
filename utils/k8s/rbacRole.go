package k8s

import (
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetRbacList() (*v1.RoleList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	role, err := cli.RbacV1().Roles("").List(metav1.ListOptions{})
	return role, err
}

func GetRbacByName(namespace, name string) (*v1.Role, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	role, err := cli.RbacV1().Roles(namespace).Get(name, metav1.GetOptions{})
	return role, err
}

func DeleteRbac(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.RbacV1().Roles(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}
