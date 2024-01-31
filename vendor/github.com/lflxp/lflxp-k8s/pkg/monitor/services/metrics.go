package services

import (
	"context"

	"github.com/lflxp/lflxp-k8s/utils/metrics"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

//	for _,v := range metric.Items {
//		for _,c := range v.Containers {
//			fmt.Println(c.Usage.Memory().Value())
//		}
//	}
func GetPodMetrics(namespace string) (data *v1beta1.PodMetricsList, err error) {
	cli, err := metrics.InitClientMetrics()
	if err != nil {
		return nil, err
	}

	data, err = cli.MetricsV1beta1().PodMetricses(namespace).List(context.Background(), v1.ListOptions{})

	return
}

func GetNodeMetrics() (data *v1beta1.NodeMetricsList, err error) {
	cli, err := metrics.InitClientMetrics()
	if err != nil {
		return nil, err
	}

	data, err = cli.MetricsV1beta1().NodeMetricses().List(context.Background(), v1.ListOptions{})

	return
}
