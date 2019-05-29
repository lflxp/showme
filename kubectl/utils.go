package kubectl

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

	// pvc
	pvcs, err := origin.ClientSet.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(pvcs.Items) > 0 {
		pvc_rs := ClusterStatus{
			Title: "Pvc",
			Count: len(pvcs.Items),
			Data:  map[string]string{},
		}
		for _, x := range pvcs.Items {
			pvc_rs.Data[x.GetName()] = x.GetName()
		}
		origin.Cluster = append(origin.Cluster, pvc_rs)
	}

	// role
	roles, err := origin.ClientSet.RbacV1().Roles("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(roles.Items) > 0 {
		role_rs := ClusterStatus{
			Title: "Role",
			Count: len(roles.Items),
			Data:  map[string]string{},
		}
		for _, x := range roles.Items {
			role_rs.Data[x.GetName()] = x.GetName()
		}
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
