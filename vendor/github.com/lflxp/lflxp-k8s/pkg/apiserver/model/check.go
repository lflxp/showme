package model

import (
	"context"
	"errors"

	"github.com/lflxp/tools/sdk/clientgo"

	log "github.com/go-eden/slf4go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type GetGVR struct {
	Group     string                     `json:"group"`
	Version   string                     `json:"version"`
	Resource  string                     `json:"resource"`
	Namespace string                     `json:"namespace"`
	Name      string                     `json:"name"`
	Data      *unstructured.Unstructured `json:"data"`
	PatchData string                     `json:"patchdata"`
}

func (g *GetGVR) GetStruct() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    g.Group,
		Version:  g.Version,
		Resource: g.Resource,
	}
}

func (g *GetGVR) List() (list *unstructured.UnstructuredList, err error) {
	cli, err := clientgo.InitClientDynamic()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	if g.Namespace != "" {
		list, err = cli.Resource(g.GetStruct()).Namespace(g.Namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		list, err = cli.Resource(g.GetStruct()).List(ctx, metav1.ListOptions{Limit: 500})
		if err != nil {
			return nil, err
		}
	}

	// target := &v1.Deployment{}

	// toBytes, err := list.MarshalJSON()
	// if err != nil {
	// 	log.Errorf("Error marshalling namespace %s: %v", data.Namespace, err)
	// 	httpclient.SendErrorMessage(c, 500, "error marshalling namespace", err.Error())
	// 	return
	// }

	// if err := json.Unmarshal(toBytes, &target.Object); err != nil {
	// 	log.Errorf("Error unmarshalling namespace %s: %v")
	// 	httpclient.SendErrorMessage(c, 500, "error Unmarshal", err.Error())
	// 	return
	// }

	return list, nil
}

func (g *GetGVR) Get() (data *unstructured.Unstructured, err error) {
	cli, err := clientgo.InitClientDynamic()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	if g.Namespace != "" {
		data, err = cli.Resource(g.GetStruct()).Namespace(g.Namespace).Get(ctx, g.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		data, err = cli.Resource(g.GetStruct()).Get(ctx, g.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (g *GetGVR) Post() (*unstructured.Unstructured, error) {
	gvr := g.GetStruct()
	cli, err := clientgo.InitClientDynamic()
	if err != nil {
		return nil, err
	}

	var list *unstructured.Unstructured

	if g.Namespace != "" {
		list, err = cli.Resource(gvr).Namespace(g.Namespace).Create(context.TODO(), g.Data, metav1.CreateOptions{})
	} else {
		list, err = cli.Resource(gvr).Create(context.TODO(), g.Data, metav1.CreateOptions{})
	}

	return list, err
}

func (g *GetGVR) Patch() (*unstructured.Unstructured, error) {
	if g.PatchData == "" {
		log.Error("PatchData is empty")
		return nil, errors.New("PatchData is empty")
	}

	gvr := g.GetStruct()
	cli, err := clientgo.InitClientDynamic()
	if err != nil {
		return nil, err
	}

	var list *unstructured.Unstructured

	if g.Namespace != "" {
		list, err = cli.Resource(gvr).Namespace(g.Namespace).Patch(context.TODO(), g.Name, types.MergePatchType, []byte(g.PatchData), metav1.PatchOptions{})
	} else {
		list, err = cli.Resource(gvr).Patch(context.TODO(), g.Name, types.MergePatchType, []byte(g.PatchData), metav1.PatchOptions{})
	}

	return list, err
}

func (g *GetGVR) Put() (*unstructured.Unstructured, error) {
	gvr := g.GetStruct()
	cli, err := clientgo.InitClientDynamic()
	if err != nil {
		return nil, err
	}

	var list *unstructured.Unstructured

	if g.Namespace != "" {
		list, err = cli.Resource(gvr).Namespace(g.Namespace).Update(context.TODO(), g.Data, metav1.UpdateOptions{})
	} else {
		list, err = cli.Resource(gvr).Update(context.TODO(), g.Data, metav1.UpdateOptions{})
	}

	return list, err
}

func (g *GetGVR) Delete() error {
	gvr := g.GetStruct()
	cli, err := clientgo.InitClientDynamic()
	if err != nil {
		return err
	}

	if g.Namespace != "" {
		err = cli.Resource(gvr).Namespace(g.Namespace).Delete(context.TODO(), g.Name, metav1.DeleteOptions{})
	} else {
		err = cli.Resource(gvr).Delete(context.TODO(), g.Name, metav1.DeleteOptions{})
	}

	return err
}
