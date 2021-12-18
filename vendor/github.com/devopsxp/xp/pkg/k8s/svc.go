package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/lflxp/lflxp-kubectl/models"
	"github.com/mattbaird/jsonpatch"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GetServiceList() (*apiv1.ServiceList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	return config, err
}

func GetServiceListByNamespace(namespace string) (*apiv1.ServiceList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	return config, err
}

func GetServiceByName(namespace, name string) (*apiv1.Service, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	config, err := cli.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return config, err
}

func DeleteService(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func CreateService(info *models.Service) (*apiv1.Service, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	label := map[string]string{
		info.Tag[0]: info.Tag[1],
	}

	serviceSpec := apiv1.ServiceSpec{Selector: label}

	// Port Ready
	ports := []apiv1.ServicePort{}

	for _, x := range info.ServicePods {
		t1 := apiv1.ServicePort{}
		if x.Portname != "" {
			t1.Name = x.Portname
		}

		switch x.Protocol {
		case "TCP":
			t1.Protocol = apiv1.ProtocolTCP
		case "UDP":
			t1.Protocol = apiv1.ProtocolUDP
		}

		iport, err := strconv.ParseInt(x.Port, 10, 64)
		if err != nil {
			return nil, err
		}
		t1.Port = int32(iport)

		if x.TargetPort != "" {
			itarget, err := strconv.ParseInt(x.TargetPort, 10, 64)
			if err != nil {
				return nil, err
			}
			t1.TargetPort = intstr.IntOrString{
				Type:   intstr.Int,
				StrVal: x.TargetPort,
				IntVal: int32(itarget),
			}
		}

		if x.NodePort != "" {
			inode, err := strconv.ParseInt(x.NodePort, 10, 64)
			if err != nil {
				return nil, err
			}
			t1.NodePort = int32(inode)
		}
		ports = append(ports, t1)
	}

	serviceSpec.Ports = ports

	if info.ClusetIP != "" {
		serviceSpec.ClusterIP = info.ClusetIP
	}

	if info.LoadIp != "" {
		serviceSpec.LoadBalancerIP = info.LoadIp
	}

	switch info.Type {
	case "NodePort":
		serviceSpec.Type = apiv1.ServiceTypeNodePort
	case "ClusterIP":
		serviceSpec.Type = apiv1.ServiceTypeClusterIP
	case "LoadBalancer":
		serviceSpec.Type = apiv1.ServiceTypeLoadBalancer
	case "ExternalName":
		serviceSpec.Type = apiv1.ServiceTypeExternalName
	}

	tmp := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.Name,
			Namespace: info.Namespace,
			Labels:    label,
		},
		Spec: serviceSpec,
	}

	if info.Ingressip != "" {
		loadingress := apiv1.LoadBalancerIngress{IP: info.Ingressip}
		tmp.Status.LoadBalancer.Ingress = []apiv1.LoadBalancerIngress{loadingress}
	}

	data, err := cli.CoreV1().Services(info.Namespace).Create(context.TODO(), tmp, metav1.CreateOptions{})
	return data, err
}

func PatchServicePort(info *models.Service) (*apiv1.Service, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	origin, err := cli.CoreV1().Services(info.Namespace).Get(context.TODO(), info.Name, metav1.GetOptions{})
	if err != nil {
		return origin, err
	}

	oJson, err := json.Marshal(origin)
	if err != nil {
		return nil, err
	}
	if info.Method == "ADD" {
		for _, x := range info.ServicePods {
			t1 := apiv1.ServicePort{}
			if x.Portname != "" {
				t1.Name = x.Portname
			}

			switch x.Protocol {
			case "TCP":
				t1.Protocol = apiv1.ProtocolTCP
			case "UDP":
				t1.Protocol = apiv1.ProtocolUDP
			}

			iport, err := strconv.ParseInt(x.Port, 10, 64)
			if err != nil {
				return nil, err
			}
			t1.Port = int32(iport)

			if x.TargetPort != "" {
				itarget, err := strconv.ParseInt(x.TargetPort, 10, 64)
				if err != nil {
					return nil, err
				}
				t1.TargetPort = intstr.IntOrString{
					Type:   intstr.Int,
					StrVal: x.TargetPort,
					IntVal: int32(itarget),
				}
			}

			if x.NodePort != "" {
				inode, err := strconv.ParseInt(x.NodePort, 10, 64)
				if err != nil {
					return nil, err
				}
				t1.NodePort = int32(inode)
			}
			origin.Spec.Ports = append(origin.Spec.Ports, t1)
		}
	} else if info.Method == "DELETE" {
		portsTmp := []apiv1.ServicePort{}
		for _, x := range origin.Spec.Ports {
			rs := false
			for _, y := range info.ServicePods {
				if fmt.Sprintf("%d", x.Port) == y.Port {
					rs = true
				}
			}
			if rs == false {
				portsTmp = append(portsTmp, x)
			}
		}
		origin.Spec.Ports = portsTmp
	}

	if info.ClusetIP != "" {
		origin.Spec.ClusterIP = info.ClusetIP
	}

	if info.LoadIp != "" {
		origin.Spec.LoadBalancerIP = info.LoadIp
	}

	if info.Type == "" {
		info.Type = "NodePort"
	}

	switch info.Type {
	case "NodePort":
		origin.Spec.Type = apiv1.ServiceTypeNodePort
	case "ClusterIP":
		origin.Spec.Type = apiv1.ServiceTypeClusterIP
	case "LoadBalancer":
		origin.Spec.Type = apiv1.ServiceTypeLoadBalancer
	case "ExternalName":
		origin.Spec.Type = apiv1.ServiceTypeExternalName
	}

	mJson, err := json.Marshal(origin)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.CreatePatch(oJson, mJson)
	if err != nil {
		return nil, err
	}

	pb, err := json.MarshalIndent(patch, "", "  ")
	if err != nil {
		return nil, err
	}

	final, err := cli.CoreV1().Services(info.Namespace).Patch(context.TODO(), info.Name, types.JSONPatchType, pb, metav1.PatchOptions{})
	return final, err
}
