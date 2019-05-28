package k8s

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var clientSet *kubernetes.Clientset
var initErr error
var configShare *rest.Config

type KubeConfig struct {
	Config string
}

func NewKubeConfig() *KubeConfig {
	home := homedir.HomeDir()
	return &KubeConfig{
		Config: filepath.Join(home, ".kube", "config"),
	}
}

func InitClientSet() {
	kube := NewKubeConfig()
	// clientSet, initErr = kube.GetClientSet()
	clientSet, initErr = kube.GetClientSetInner()
	if initErr != nil {
		// clientSet, initErr = kube.GetClientSetInner()
		clientSet, initErr = kube.GetClientSet()
	}
}

func GetClientSet() (*kubernetes.Clientset, error) {
	return clientSet, initErr
}

func (this *KubeConfig) GetClientSetInner() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	configShare = config
	return clientset, nil
}

func (this *KubeConfig) GetClientSet() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", this.Config)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	configShare = config
	return clientset, nil
}

func GetClientSetAndConfig() (*kubernetes.Clientset, *rest.Config) {
	return clientSet, configShare
}
