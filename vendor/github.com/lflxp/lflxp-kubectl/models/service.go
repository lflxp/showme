package models

type ServicePods struct {
	Index      int64  `json:"index"`
	NodePort   string `json:"nodeport"`
	Port       string `json:"port"`
	Portname   string `json:"portname"`
	Protocol   string `json:"protocol"`
	TargetPort string `json:"targetport"`
}

type Service struct {
	Name        string        `json:"name"`
	Namespace   string        `json:"namespace"`
	Type        string        `json:"type"`
	ServicePods []ServicePods `json:"servicePods"`
	Tag         []string      `json:"tag"`
	ClusetIP    string        `json:"clusterip"`
	Ingressip   string        `json:"ingressip"`
	LoadIp      string        `json:"loadip"`
	Method      string        `json:"method"` //  patch add or delete
}
