package models

import (
	"encoding/json"

	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LimitRangeList struct {
	apiv1.LimitRangeList
}

type LimitRange struct {
	apiv1.LimitRange
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Max struct {
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
}

type Min struct {
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
}

type Default struct {
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
}

type DefaultRequest struct {
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
}

type Limits struct {
	Type                 string         `json:"type"`
	Max                  Max            `json:"max"`
	Min                  Min            `json:"min"`
	Default              Default        `json:"default"`
	DefaultRequest       DefaultRequest `json:"defaultRequest"`
	MaxLimitRequestRatio string         `json:"maxLimitRequestRatio"`
}

type Spec struct {
	Limits []Limits `json:"limits"`
}

type LimitRanged struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
}

func (this *LimitRanged) Unmarshal(data []byte) error {
	err := json.Unmarshal([]byte(data), this)
	return err
}

func (this *LimitRanged) Resolved() *apiv1.LimitRange {
	//LimitRangeItem
	tmp_lri := []apiv1.LimitRangeItem{}
	for _, data := range this.Spec.Limits {
		tmp := apiv1.LimitRangeItem{}
		if data.Type != "" {
			if data.Type == "Pod" {
				tmp.Type = apiv1.LimitTypePod
			} else if data.Type == "Container" {
				tmp.Type = apiv1.LimitTypeContainer
			} else if data.Type == "PersistentVolumeClaim" {
				tmp.Type = apiv1.LimitTypePersistentVolumeClaim
			}
		}

		if data.Max.Cpu != "" || data.Max.Memory != "" || data.Max.Storage != "" {
			r1 := apiv1.ResourceList{}
			if data.Max.Cpu != "" {
				r1[apiv1.ResourceCPU] = resource.MustParse(data.Max.Cpu)
			}
			if data.Max.Memory != "" {
				r1[apiv1.ResourceMemory] = resource.MustParse(data.Max.Memory)
			}
			if data.Max.Storage != "" {
				r1[apiv1.ResourceStorage] = resource.MustParse(data.Max.Storage)
			}
			tmp.Max = r1
		}

		if data.Min.Cpu != "" || data.Min.Memory != "" || data.Min.Storage != "" {
			r2 := apiv1.ResourceList{}
			if data.Min.Cpu != "" {
				r2[apiv1.ResourceCPU] = resource.MustParse(data.Min.Cpu)
			}
			if data.Min.Memory != "" {
				r2[apiv1.ResourceMemory] = resource.MustParse(data.Min.Memory)
			}
			if data.Min.Storage != "" {
				r2[apiv1.ResourceStorage] = resource.MustParse(data.Min.Storage)
			}
			tmp.Min = r2
		}

		if data.Default.Cpu != "" || data.Default.Memory != "" || data.Default.Storage != "" {
			r3 := apiv1.ResourceList{}
			if data.Default.Cpu != "" {
				r3[apiv1.ResourceCPU] = resource.MustParse(data.Default.Cpu)
			}
			if data.Default.Memory != "" {
				r3[apiv1.ResourceMemory] = resource.MustParse(data.Default.Memory)
			}
			if data.Default.Storage != "" {
				r3[apiv1.ResourceStorage] = resource.MustParse(data.Default.Storage)
			}
			tmp.Default = r3
		}

		if data.DefaultRequest.Cpu != "" || data.DefaultRequest.Memory != "" || data.DefaultRequest.Storage != "" {
			r4 := apiv1.ResourceList{}
			if data.DefaultRequest.Cpu != "" {
				r4[apiv1.ResourceCPU] = resource.MustParse(data.DefaultRequest.Cpu)
			}
			if data.DefaultRequest.Memory != "" {
				r4[apiv1.ResourceMemory] = resource.MustParse(data.DefaultRequest.Memory)
			}
			if data.DefaultRequest.Storage != "" {
				r4[apiv1.ResourceStorage] = resource.MustParse(data.DefaultRequest.Storage)
			}
			tmp.DefaultRequest = r4
		}
		tmp_lri = append(tmp_lri, tmp)
	}

	lr := &apiv1.LimitRange{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LimitRange",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      this.Metadata.Name,
			Namespace: this.Metadata.Namespace,
		},
		Spec: apiv1.LimitRangeSpec{
			Limits: tmp_lri,
		},
	}

	return lr
}
