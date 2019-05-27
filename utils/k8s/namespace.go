package k8s

import (
	"errors"

	"github.com/lflxp/showme/utils/k8s/models"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespacesByName(name string) (*apiv1.Namespace, error) {
	clientset, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	namespace, err := clientset.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	return namespace, err
}

func GetNamespaces() (*apiv1.NamespaceList, error) {
	clientset, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	return namespaces, err
}

func CreateNamespaces(name string) error {
	clientset, err := GetClientSet()
	if err != nil {
		return err
	}

	ns := &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err = clientset.CoreV1().Namespaces().Create(ns)
	return err
}

func DeleteNamespaces(name string) error {
	clientset, err := GetClientSet()
	if err != nil {
		return err
	}

	err = clientset.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
	return err
}

func MutilDeleteNamespace(name models.Namespaces) error {
	if len(name.Names) == 0 {
		return errors.New("none of namespace product")
	}
	for _, x := range name.Names {
		err := DeleteNamespaces(x)
		if err != nil {
			break
			return err
		}
	}
	return nil
}
