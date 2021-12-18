package k8s

import (
	"context"

	apiv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetJobList(namespace string) (*apiv1.JobList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	return config, err
}

func GetJobListByLabels(namespace, label string) (*apiv1.JobList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	listoptions := metav1.ListOptions{
		LabelSelector: label,
	}

	config, err := cli.BatchV1().Jobs(namespace).List(context.TODO(), listoptions)
	return config, err
}

func GetJobByName(namespace, name string) (*apiv1.Job, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return config, err
}

func DeleteJob(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

// func GetJobLogByJobId(namespace, Jobid string) (string, error) {
// 	var result string
// 	f := func(s int64) *int64 {
// 		return &s
// 	}

// 	cli, err := GetClientSet()
// 	if err != nil {
// 		return result, err
// 	}

// 	options := apiv1.JobLogOptions{
// 		Follow:    false,
// 		TailLines: f(1000),
// 	}
// 	req := cli.CoreV1().Jobs(namespace).GetLogs(Jobid, &options)
// 	readCloser, err := req.Stream()
// 	if err != nil {
// 		return result, err
// 	}
// 	defer readCloser.Close()
// 	var out bytes.Buffer
// 	_, err = io.Copy(&out, readCloser)
// 	if err != nil {
// 		return result, err
// 	}
// 	result = out.String()
// 	return result, nil
// }
