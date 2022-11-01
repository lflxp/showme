package model

import (
	"context"

	"github.com/lflxp/tools/sdk/clientgo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CoreV1 struct {
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
