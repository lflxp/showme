package models

type DynamicPVC struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Region    []string `json:"region"`
	Size      int64    `json:"size"`
	Storage   string   `json:"storage"`
	Unit      string   `json:"unit"`
	MaxSize   int64    `json:"maxsize"`
}
