package k8s

import (
	"fmt"

	"github.com/lflxp/showme/utils/k8s/models"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPvcList() (*apiv1.PersistentVolumeClaimList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	return config, err
}

func GetPvcByName(namespace, name string) (*apiv1.PersistentVolumeClaim, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().PersistentVolumeClaims(namespace).Get(name, metav1.GetOptions{})
	return config, err
}

/*
// Resource names must be not more than 63 characters, consisting of upper- or lower-case alphanumeric characters,
// with the -, _, and . characters allowed anywhere, except the first or last character.
// The default convention, matching that for annotations, is to use lower-case names, with dashes, rather than
// camel case, separating compound words.
// Fully-qualified resource typenames are constructed from a DNS-style subdomain, followed by a slash `/` and a name.
const (
	// CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage ResourceName = "storage"
	// Local ephemeral storage, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	// The resource name for ResourceEphemeralStorage is alpha and it can change across releases.
	ResourceEphemeralStorage ResourceName = "ephemeral-storage"
)

const (
	// Default namespace prefix.
	ResourceDefaultNamespacePrefix = "kubernetes.io/"
	// Name prefix for huge page resources (alpha).
	ResourceHugePagesPrefix = "hugepages-"
)

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[ResourceName]resource.Quantity

[root@k8s-master1 ~]# kubectl describe pvc glusterfs-heketi-mysql1
Name:          glusterfs-heketi-mysql1
Namespace:     default
StorageClass:  heketi-kubernetes
Status:        Bound
Volume:        pvc-a8df095b-78ea-11e8-8a36-000c29015d43
Labels:        <none>
Annotations:   pv.kubernetes.io/bind-completed=yes
               pv.kubernetes.io/bound-by-controller=yes
               volume.beta.kubernetes.io/storage-class=heketi-kubernetes
               volume.beta.kubernetes.io/storage-provisioner=kubernetes.io/glusterfs
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:      2Gi
Access Modes:  RWX
Events:        <none>

**/

func CreatePvc(data *models.DynamicPVC) (*apiv1.PersistentVolumeClaim, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	accessModes := []apiv1.PersistentVolumeAccessMode{}
	for _, x := range data.Region {
		switch x {
		case "ReadWriteOnce":
			accessModes = append(accessModes, apiv1.ReadWriteOnce)
		case "ReadOnlyMany":
			accessModes = append(accessModes, apiv1.ReadOnlyMany)
		case "ReadWriteMany":
			accessModes = append(accessModes, apiv1.ReadWriteMany)
		default:
			log.Warning(fmt.Sprintf("%s is not in my range", x))
		}
	}

	resourceRequire := apiv1.ResourceRequirements{}
	if data.MaxSize != 0 {
		resourceRequire.Limits = apiv1.ResourceList{
			apiv1.ResourceStorage: resource.MustParse(fmt.Sprintf("%d%s", data.MaxSize, data.Unit)),
		}
	}
	resourceRequire.Requests = apiv1.ResourceList{
		apiv1.ResourceStorage: resource.MustParse(fmt.Sprintf("%d%s", data.Size, data.Unit)),
	}

	pvc := &apiv1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &data.Storage,
			Resources:        resourceRequire,
		},
	}

	config, err := cli.CoreV1().PersistentVolumeClaims(data.Namespace).Create(pvc)
	return config, err
}

func DeletePvc(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().PersistentVolumeClaims(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}
