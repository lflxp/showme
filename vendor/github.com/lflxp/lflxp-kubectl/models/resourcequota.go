package models

// import (
// 	"gitlab.yc/lixueping/k8s/pkg"
// 	apiv1 "k8s.io/api/core/v1"
// )

type Resource struct {
	Cpu            string `json:"cpu"`
	Memory         string `json:"memory"`
	LimitsCpu      string `json:"limitsCpu"`
	LimitsMemory   string `json:"limitsMemory"`
	RequestsCpu    string `json:"requestsCpu"`
	RequestsMemory string `json:"requestsMemory"`
	NvidiaGpu      string `json:"gpu"`
}

type Storage struct {
	RequestsStorage                    string `json:"requestsStorage"`
	Persistentvolumeclaims             string `json:"persistentvolumeclaims"`
	StorageClassRequestsStorage        string `json:"requestsStorageClass"`
	StorageClassPersistentvolumeclaims string `json:"persistentvolumeclaimsClass"`
}

type Objected struct {
	Configmaps             string `json:"configmaps"`
	Pods                   string `json:"pods"`
	ReplicationControllers string `json:"replicationcontrollers"`
	ResourceQuotas         string `json:"resourcequotas"`
	Services               string `json:"services"`
	ServicesLoadbalances   string `json:"servicesLoadbalances"`
	ServicesNodeports      string `json:"servicesNodeports"`
	Secrets                string `json:"secrets"`
	Persistentvolumeclaims string `json:"persistentvolumeclaims"`
}

type Domain struct {
	Terminating    string `json:"Terminating"`
	NotTerminating string `json:"NotTerminating"`
	BestEffort     string `json:"BestEffort"`
	NotBestEffort  string `json:"NotBestEffort"`
}

/*
// A ResourceQuotaScope defines a filter that must match each object tracked by a quota
type ResourceQuotaScope string

const (
	// Match all pod objects where spec.activeDeadlineSeconds
	ResourceQuotaScopeTerminating ResourceQuotaScope = "Terminating"
	// Match all pod objects where !spec.activeDeadlineSeconds
	ResourceQuotaScopeNotTerminating ResourceQuotaScope = "NotTerminating"
	// Match all pod objects that have best effort quality of service
	ResourceQuotaScopeBestEffort ResourceQuotaScope = "BestEffort"
	// Match all pod objects that do not have best effort quality of service
	ResourceQuotaScopeNotBestEffort ResourceQuotaScope = "NotBestEffort"
)
*/
// type Scope string

// type Spec struct {
// 	Resource
// 	Storage
// 	Objected
// 	Scope
// }

type ResourceQuota struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Resource  Resource `json:"resource"`
	Storage   Storage  `json:"storage"`
	Objected  Objected `json:"object"`
	Domain    Domain   `json:"domain"`
}
