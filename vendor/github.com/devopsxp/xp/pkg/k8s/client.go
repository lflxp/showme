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

// 初始化k8s client实例
func init() {
	InitClientSet()
}

type KubeConfig struct {
	Config string
}

func newKubeConfig() *KubeConfig {
	home := homedir.HomeDir()
	return &KubeConfig{
		Config: filepath.Join(home, ".kube", "config"),
	}
}

// 兼容k8s inner Cluster 和 Outter Cluster client实例
func InitClientSet() {
	kube := newKubeConfig()
	// clientSet, initErr = kube.GetClientSet()
	clientSet, initErr = kube.GetClientSetInner()
	if initErr != nil {
		// clientSet, initErr = kube.GetClientSetInner()
		clientSet, initErr = kube.GetClientSet()
	}
}

// 返回k8s client实例和初始化错误信息
func GetClientSet() (*kubernetes.Clientset, error) {
	return clientSet, initErr
}

// k8s集群内部client
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

// k8s config实例初始化方法
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
