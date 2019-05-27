package models

type ConfigmapsJson struct {
	Data        map[string]string `json:"data"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}
