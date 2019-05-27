package k8s

import (
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeploymentList() (*v1.DeploymentList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	d, err := cli.AppsV1().Deployments("").List(metav1.ListOptions{})
	return d, err
}

func GetDeploymentListByNamespace(namespace string) (*v1.DeploymentList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	d, err := cli.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	return d, err
}

func GetDeploymentByName(namespace, name string) (*v1.Deployment, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	return config, err
}

func DeleteDeployment(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}

func ScaleDeployment(namespace, name string, num int32) (*v1.Deployment, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}
	dp, err := GetDeploymentByName(namespace, name)
	if err != nil {
		return nil, err
	}

	dp.Spec.Replicas = &num
	dps, err := cli.AppsV1().Deployments(namespace).Update(dp)
	return dps, err
}

func ChangeImageDeployment(namespace, name, image string) (*v1.Deployment, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}
	dp, err := GetDeploymentByName(namespace, name)
	if err != nil {
		return nil, err
	}

	dp.Spec.Template.Spec.Containers[0].Image = image
	dps, err := cli.AppsV1().Deployments(namespace).Update(dp)
	return dps, err
}
