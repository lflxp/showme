package models

type NodeLabels struct {
	Node      string            `json:"node"`
	Value     map[string]string `json:"values"`
	Overwrite bool              `json:"overwrite"`
}
