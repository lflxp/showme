package kubectl

import (
	"fmt"

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

// 获取集群所有负载状态信息
// cronjob daemonset deploymenet job pod replicaset statefulset
func GetLoadStatuses() error {
	return nil
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
