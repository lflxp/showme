package k8s

import (
	"github.com/lflxp/showme/utils/k8s/models"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNodeList() (*apiv1.NodeList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Nodes().List(metav1.ListOptions{})
	return config, err
}

func GetNodeByName(name string) (*apiv1.Node, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Nodes().Get(name, metav1.GetOptions{})
	return config, err
}

func DeleteNode(name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().Nodes().Delete(name, &metav1.DeleteOptions{})
	return err
}

func AddNodeLabel(data models.NodeLabels) (*apiv1.Node, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	node, err := cli.CoreV1().Nodes().Get(data.Node, metav1.GetOptions{})
	if err != nil {
		return node, err
	}
	if data.Overwrite {
		for k, v := range data.Value {
			node.ObjectMeta.Labels[k] = v
		}
	} else {
		for k, v := range data.Value {
			if _, ok := node.ObjectMeta.Labels[k]; !ok {
				node.ObjectMeta.Labels[k] = v
			}
		}
	}

	cnode, err := cli.CoreV1().Nodes().Update(node)
	return cnode, err
}

func DeleteNodeLabel(nodename, key string) (*apiv1.Node, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	node, err := cli.CoreV1().Nodes().Get(nodename, metav1.GetOptions{})
	if err != nil {
		return node, err
	}

	delete(node.ObjectMeta.Labels, key)

	cnode, err := cli.CoreV1().Nodes().Update(node)
	return cnode, err
}
