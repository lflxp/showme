package clientgo

import (
	"flag"
	"log/slog"
	"path/filepath"
	"sync"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

var (
	onceMetrics    sync.Once
	clientsMetrics *versioned.Clientset
)

func InitClientMetrics() (*versioned.Clientset, error) {
	var err error
	onceMetrics.Do(func() {
		clientsMetrics, err = doMetricsInit()
		if err != nil {
			clientsMetrics, err = doMetricsInnerInit()
			slog.Debug("init out of cluster", "Error", err.Error())
			if err != nil {
				panic(err)
			}
		}
	})

	return clientsMetrics, err
}

func doMetricsInnerInit() (client *versioned.Clientset, err error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return client, err
	}

	client = versioned.NewForConfigOrDie(restConfig)

	return
}

func doMetricsInit() (client *versioned.Clientset, err error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfigMetrics", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfigMetrics", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	restConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	client = versioned.NewForConfigOrDie(restConfig)
	return
}
