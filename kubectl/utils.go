package kubectl

import (
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 获取集群服务发现与负载均衡 配置和存储
// service ingress  configmap secret pvc
func GetServiceConfigStatus() error {
	// service
	service, err := origin.ClientSet.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(service.Items) > 0 {
		s_rs := ClusterStatus{
			Title: "Service",
			Count: len(service.Items),
			Data:  map[string]string{},
		}
		for _, x := range service.Items {
			s_rs.Data[fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())] = fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())
		}
		origin.ServiceConfig = append(origin.ServiceConfig, s_rs)
	}

	// ingress
	ingress, err := origin.ClientSet.Extensions().Ingresses("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(ingress.Items) > 0 {
		ingress_rs := ClusterStatus{
			Title: "Ingress",
			Count: len(ingress.Items),
			Data:  map[string]string{},
		}
		for _, x := range service.Items {
			ingress_rs.Data[fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())] = fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())
		}
		origin.ServiceConfig = append(origin.ServiceConfig, ingress_rs)
	}

	// pvc
	pvc, err := origin.ClientSet.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(pvc.Items) > 0 {
		pvc_rs := ClusterStatus{
			Title: "Pvc",
			Count: len(pvc.Items),
			Data:  map[string]string{},
		}
		for _, x := range pvc.Items {
			pvc_rs.Data[fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())] = fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())
		}
		origin.ServiceConfig = append(origin.ServiceConfig, pvc_rs)
	}

	// configmap
	configmap, err := origin.ClientSet.CoreV1().ConfigMaps("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(configmap.Items) > 0 {
		configmap_rs := ClusterStatus{
			Title: "Configmap",
			Count: len(configmap.Items),
			Data:  map[string]string{},
		}
		for _, x := range service.Items {
			configmap_rs.Data[fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())] = fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())
		}
		origin.ServiceConfig = append(origin.ServiceConfig, configmap_rs)
	}

	// secrets
	secrets, err := origin.ClientSet.CoreV1().Secrets("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(secrets.Items) > 0 {
		secrets_rs := ClusterStatus{
			Title: "Secrets",
			Count: len(secrets.Items),
			Data:  map[string]string{},
		}
		for _, x := range secrets.Items {
			secrets_rs.Data[fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())] = fmt.Sprintf("%s::%s", x.GetNamespace(), x.GetName())
		}
		origin.ServiceConfig = append(origin.ServiceConfig, secrets_rs)
	}
	return nil
}

// 判断origin.Pods是否更改
func isPodsChanges(info PodStatus) bool {
	isChange := true
	for _, x := range origin.Pods {
		if x.Name == info.Name && x.Namespace == info.Namespace && x.Ready == info.Ready && x.Restarts == info.Restarts {
			isChange = false
			break
		}
	}
	return isChange
}

// 判断origin.PodController是否更改
func isPodControllersChanges(info PodController) bool {
	isChange := true
	for _, x := range origin.PodControllers {
		if x.Name == info.Name && x.Namespace == info.Namespace && x.Type == info.Type && x.ContainerGroup == info.ContainerGroup {
			isChange = false
			break
		}
	}
	return isChange
}

// 获取集群所有负载状态信息
// +状态自动刷新
// cronjob daemonset deploymenet job pod replicaset statefulset
func GetLoadStatuses() (bool, error) {
	num := 0
	var isChange bool
	tmpPods := []PodStatus{}
	tmpPodControllers := []PodController{}

	// pods
	pods, err := origin.ClientSet.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return isChange, err
	}
	if len(pods.Items) > 0 {
		for _, x := range pods.Items {
			tmp := PodStatus{
				Name:      x.GetName(),
				Namespace: x.GetNamespace(),
				Node:      x.Spec.NodeName,
				Ready:     x.Status.Phase,
				Restarts:  fmt.Sprintf("%d", x.Status.ContainerStatuses[0].RestartCount),
				Time:      strings.Replace(fmt.Sprintf("%v", x.Status.StartTime.Sub(time.Now())), "-", "", -1),
			}
			if isPodsChanges(tmp) {
				num++
			}
			tmpPods = append(tmpPods, tmp)
		}
	}

	// daemonset
	dm, err := origin.ClientSet.Extensions().DaemonSets("").List(metav1.ListOptions{})
	if err != nil {
		return isChange, err
	}
	if len(dm.Items) > 0 {
		for _, x := range dm.Items {
			tmp := PodController{
				Type:           "DaemonSet",
				Name:           x.GetName(),
				Namespace:      x.GetNamespace(),
				Tags:           x.Labels,
				ContainerGroup: fmt.Sprintf("%d / %d", x.Status.NumberAvailable, x.Status.DesiredNumberScheduled),
				Time:           strings.Replace(fmt.Sprintf("%v", x.CreationTimestamp.Sub(time.Now())), "-", "", -1),
				Images:         x.Spec.Template.Spec.Containers[0].Image,
			}
			if isPodControllersChanges(tmp) {
				num++
			}
			tmpPodControllers = append(tmpPodControllers, tmp)
		}
	}

	// Deployments
	// w, _ := origin.ClientSet.Extensions().Deployments("").Watch(metav1.ListOptions{})
	// www := w.ResultChan()
	// xxx := <-www
	// xxx.Type

	deploy, err := origin.ClientSet.Extensions().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		return isChange, err
	}
	if len(deploy.Items) > 0 {
		for _, x := range deploy.Items {
			tmp := PodController{
				Type:           "Deployments",
				Name:           x.GetName(),
				Namespace:      x.GetNamespace(),
				Tags:           x.Labels,
				ContainerGroup: fmt.Sprintf("%d / %d", x.Status.AvailableReplicas, x.Status.Replicas),
				Time:           strings.Replace(fmt.Sprintf("%v", x.CreationTimestamp.Sub(time.Now())), "-", "", -1),
				Images:         x.Spec.Template.Spec.Containers[0].Image,
			}
			if isPodControllersChanges(tmp) {
				num++
			}
			tmpPodControllers = append(tmpPodControllers, tmp)
		}
	}

	// job
	job, err := origin.ClientSet.BatchV2alpha1().CronJobs("").List(metav1.ListOptions{})
	if err != nil {
		return isChange, err
	}
	if len(job.Items) > 0 {
		for _, x := range job.Items {
			tmp := PodController{
				Type:           "CronJobs",
				Name:           x.GetName(),
				Namespace:      x.GetNamespace(),
				Tags:           x.Labels,
				ContainerGroup: x.Spec.Schedule,
				Time:           strings.Replace(fmt.Sprintf("%v", x.CreationTimestamp.Sub(time.Now())), "-", "", -1),
				Images:         x.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image,
			}
			if isPodControllersChanges(tmp) {
				num++
			}
			tmpPodControllers = append(tmpPodControllers, tmp)
		}
	}

	// repl
	rep, err := origin.ClientSet.Extensions().ReplicaSets("").List(metav1.ListOptions{})
	if err != nil {
		return isChange, err
	}
	if len(rep.Items) > 0 {
		for _, x := range rep.Items {
			tmp := PodController{
				Type:           "ReplicaSets",
				Name:           x.GetName(),
				Namespace:      x.GetNamespace(),
				Tags:           x.Labels,
				ContainerGroup: fmt.Sprintf("%d / %d", x.Status.AvailableReplicas, x.Status.Replicas),
				Time:           strings.Replace(fmt.Sprintf("%v", x.CreationTimestamp.Sub(time.Now())), "-", "", -1),
				Images:         x.Spec.Template.Spec.Containers[0].Image,
			}
			if isPodControllersChanges(tmp) {
				num++
			}
			tmpPodControllers = append(tmpPodControllers, tmp)
		}
	}

	// statefulset
	sf, err := origin.ClientSet.AppsV1().StatefulSets("").List(metav1.ListOptions{})
	if err != nil {
		return isChange, err
	}
	if len(sf.Items) > 0 {
		for _, x := range sf.Items {
			tmp := PodController{
				Type:           "StatefulSets",
				Name:           x.GetName(),
				Namespace:      x.GetNamespace(),
				Tags:           x.Labels,
				ContainerGroup: fmt.Sprintf("%d / %d", x.Status.ReadyReplicas, x.Status.Replicas),
				Time:           strings.Replace(fmt.Sprintf("%v", x.CreationTimestamp.Sub(time.Now())), "-", "", -1),
				Images:         x.Spec.Template.Spec.Containers[0].Image,
			}
			if isPodControllersChanges(tmp) {
				num++
			}
			tmpPodControllers = append(tmpPodControllers, tmp)
		}
	}
	if num > 0 {
		isChange = true
		// clear and refresh
		origin.Pods = tmpPods
		origin.PodControllers = tmpPodControllers
	} else {
		isChange = false
	}

	return isChange, nil
}

// 获取所有集群信息
// namespace node pv role storageclass
func GetClusterStatuses() error {
	// namespace
	ns, err := origin.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(ns.Items) > 0 {
		ns_rs := ClusterStatus{
			Title: "Namespace",
			Count: len(ns.Items),
			Data:  map[string]string{},
		}
		for _, x := range ns.Items {
			ns_rs.Data[x.GetName()] = x.GetName()
		}
		origin.Cluster = append(origin.Cluster, ns_rs)
	}

	// node
	nodes, err := origin.ClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(nodes.Items) > 0 {
		nodes_rs := ClusterStatus{
			Title: "Node",
			Count: len(nodes.Items),
			Data:  map[string]string{},
		}
		for _, x := range nodes.Items {
			nodes_rs.Data[x.GetName()] = x.GetName()
		}
		origin.Cluster = append(origin.Cluster, nodes_rs)
	}

	// pv
	// pvcs, err := origin.ClientSet.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	pvs, err := origin.ClientSet.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(pvs.Items) > 0 {
		pv_rs := ClusterStatus{
			Title: "Pv",
			Count: len(pvs.Items),
			Data:  map[string]string{},
		}
		for _, x := range pvs.Items {
			pv_rs.Data[x.GetName()] = x.GetName()
		}
		origin.Cluster = append(origin.Cluster, pv_rs)
	}

	// role
	roles, err := origin.ClientSet.RbacV1().Roles("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	clusterroles, err := origin.ClientSet.RbacV1().ClusterRoles().List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(roles.Items) > 0 {
		role_rs := ClusterStatus{
			Title: "Role(** clusterrole * role)",
			Data:  map[string]string{},
		}
		for _, x := range roles.Items {
			role_rs.Data[x.GetName()] = x.GetName()
		}
		for _, y := range clusterroles.Items {
			role_rs.Data[fmt.Sprintf("*%s", y.GetName())] = fmt.Sprintf("*%s", y.GetName())
		}
		role_rs.Count = len(role_rs.Data)
		origin.Cluster = append(origin.Cluster, role_rs)
	}

	// storageclass
	stc, err := origin.ClientSet.StorageV1().StorageClasses().List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(stc.Items) > 0 {
		stc_rs := ClusterStatus{
			Title: "StorageClasses",
			Count: len(stc.Items),
			Data:  map[string]string{},
		}
		for _, x := range stc.Items {
			stc_rs.Data[x.GetName()] = x.GetName()
		}
		origin.Cluster = append(origin.Cluster, stc_rs)
	}
	return nil
}
