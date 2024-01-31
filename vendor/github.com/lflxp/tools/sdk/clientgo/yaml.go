package clientgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"go/ast"
	"go/token"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yaml2 "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
)

func InstallYaml(dataRaw []byte) error {
	var decUnstructured = yaml2.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	// 1. Prepare a RESTMapper to find GVR
	dc := InitClientDiscovery()
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := InitClientDynamic()
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(dataRaw, nil, obj)
	if err != nil {
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// 6. Marshal object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// 7. Create or Update the object with SSA
	//     types.ApplyPatchType indicates SSA.
	//     FieldManager specifies the field owner ID.
	_, err = dr.Patch(context.Background(), obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "lflxp-k8s-controller",
	})

	return err
}

func UnInstallYaml(dataRaw []byte) error {
	var decUnstructured = yaml2.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	// 1. Prepare a RESTMapper to find GVR
	dc := InitClientDiscovery()
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := InitClientDynamic()
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(dataRaw, nil, obj)
	if err != nil {
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// 6. find object name
	tmpName := obj.GetName()

	// 7. Delete the object with SSA
	err = dr.Delete(context.TODO(), tmpName, metav1.DeleteOptions{})

	return err
}

func InstallYamlFilename(clientSet *kubernetes.Clientset, dynamicClient dynamic.Interface, ns string, filename string) error {
	f, err := os.Open(filename)
	slog.Info("=====================>" + filename)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	d := yaml.NewYAMLOrJSONDecoder(f, 4096)
	dc := clientSet.Discovery()

	restMapperRes, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	restMapper := restmapper.NewDiscoveryRESTMapper(restMapperRes)

	for {
		ext := runtime.RawExtension{}

		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			slog.Error(err.Error())
		}

		// runtime.Object
		obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		fmt.Printf("mapping:%+v\n", mapping)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		// runtime.Object转换为unstructed
		unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		// fmt.Printf("unstructuredObj: %+v", unstructuredObj)

		var unstruct unstructured.Unstructured

		unstruct.Object = unstructuredObj

		if ns == "" {
			res, err := dynamicClient.Resource(mapping.Resource).Create(context.TODO(), &unstruct, metav1.CreateOptions{})
			if err != nil {
				slog.Error(err.Error())
				return err
			}

			GuessType(res)
		} else {
			res, err := dynamicClient.Resource(mapping.Resource).Namespace(ns).Create(context.TODO(), &unstruct, metav1.CreateOptions{})
			if err != nil {
				slog.Error(err.Error())
				return err
			}
			GuessType(res)
		}

	}

	return nil
}

func GuessType(obj interface{}) {
	fset := token.NewFileSet()
	ast.Print(fset, obj)
}

func UninstallYaml(clientSet *kubernetes.Clientset, dynamicClient dynamic.Interface, ns string, filename string) error {
	f, err := os.Open(filename)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	d := yaml.NewYAMLOrJSONDecoder(f, 4096)
	dc := clientSet.Discovery()

	restMapperRes, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	restMapper := restmapper.NewDiscoveryRESTMapper(restMapperRes)

	for {
		ext := runtime.RawExtension{}

		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			slog.Error(err.Error())
		}

		// runtime.Object
		obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		// fmt.Printf("mapping:%+v\n", mapping)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		// runtime.Object转换为unstructed
		unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		// fmt.Printf("unstructuredObj: %+v", unstructuredObj)

		var unstruct unstructured.Unstructured

		unstruct.Object = unstructuredObj

		tmpMetadata := unstructuredObj["metadata"].(map[string]interface{})
		tmpName := tmpMetadata["name"].(string)
		tmpKind := unstructuredObj["kind"].(string)
		slog.Info("deleting resource name: " + tmpName + ", kind: " + tmpKind + ", ns: " + ns)

		if ns == "" {
			err := dynamicClient.Resource(mapping.Resource).Delete(context.TODO(), tmpName, metav1.DeleteOptions{})
			if err != nil {
				slog.Error(err.Error())
				return err
			}
		} else {
			err := dynamicClient.Resource(mapping.Resource).Namespace(ns).Delete(context.TODO(), tmpName, metav1.DeleteOptions{})
			if err != nil {
				slog.Error(err.Error())
				return err
			}
		}

	}

	return nil
}
