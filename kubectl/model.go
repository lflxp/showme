package kubectl

import (
	"time"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils/k8s"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	err    error
	origin *BasicKubectl
)

func ManualInit() {
	origin = NewBasicKubectl()
	err = origin.NewGui()
	if err != nil {
		log.Panicln(err)
	}
	defer origin.Gui.Close()

	// init clienetset
	k8s.InitClientSet()
	origin.ClientSet, err = k8s.GetClientSet()
	if err != nil {
		log.Panicln(err)
	}
	origin.DefaultNS = "default"

	// get basic info
	err = GetClusterStatuses()
	if err != nil {
		log.Error(err.Error())
	}

	err = GetServiceConfigStatus()
	if err != nil {
		log.Error(err.Error())
	}

	err = GetLoadStatuses()
	if err != nil {
		log.Error(err.Error())
	}

	// init gocui

	origin.Gui.Highlight = true
	origin.Gui.Cursor = true
	origin.Gui.SelFgColor = gocui.ColorGreen
	origin.Gui.SetManagerFunc(dashboard)

	d := time.Duration(5 * time.Second)
	t := time.NewTicker(d)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				err = GetLoadStatuses()
				if err != nil {
					log.Error(err.Error())
				}
				maxX, maxY := origin.Gui.Size()

				if err := RefreshWorkLoad(origin.Gui, 0, maxY/2, maxX-1, maxY*3/4-1); err != nil {
					log.Error(err.Error())
				}
				if err := RefreshPods(origin.Gui, 0, maxY*3/4, maxX-1, maxY-1); err != nil {
					log.Error(err.Error())
				}
				origin.Gui.Update(func(g *gocui.Gui) error { return nil })
				// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				// fmt.Fprintln(v, )
			}
		}
	}()

	if err := KeyDashboard(origin.Gui); err != nil {
		log.Panicln(err.Error())
	}

	if err := origin.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err.Error())
	}
}

func NewBasicKubectl() *BasicKubectl {
	return &BasicKubectl{}
}

// cluster
type ClusterStatus struct {
	Title string
	Data  map[string]string
	Count int
}

// load status
type PodController struct {
	Type           string
	Name           string
	Namespace      string
	Tags           map[string]string
	ContainerGroup string
	Time           string
	Images         string
}

type PodStatus struct {
	Name      string
	Namespace string
	Node      string
	Ready     v1.PodPhase
	Restarts  string
	Time      string
}

// Global Values
type BasicKubectl struct {
	// gocui
	Gui *gocui.Gui
	// kubectl
	ClientSet      *kubernetes.Clientset
	DefaultNS      string   // current namespace
	Helps          []string // F1 View show help message
	Cluster        []ClusterStatus
	ServiceConfig  []ClusterStatus
	PodControllers []PodController
	Pods           []PodStatus
	BeforeSearch   string // before search view name
}

func (this *BasicKubectl) NewGui() error {
	this.Gui, err = gocui.NewGui(gocui.OutputNormal)
	return err
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func dquit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "bottom" {
		_, err := setCurrentViewOnTop(g, "pod")

		return err
	} else if v == nil || v.Name() == "Namespace" {
		_, err := setCurrentViewOnTop(g, "Node")
		return err
	} else if v == nil || v.Name() == "Node" {
		if _, err := g.View("Pv"); err != nil {
			if err != nil {
				_, err := setCurrentViewOnTop(g, "Role(** clusterrole * role)")
				return err
			}
		} else {
			_, err := setCurrentViewOnTop(g, "Pv")
			return err
		}
	} else if v == nil || v.Name() == "Pv" {
		_, err := setCurrentViewOnTop(g, "Role(** clusterrole * role)")
		return err
	} else if v == nil || v.Name() == "Role(** clusterrole * role)" {
		if _, err := g.View("StorageClasses"); err != nil {
			if err != nil {
				_, err := setCurrentViewOnTop(g, "Service")
				return err
			}
		} else {
			_, err := setCurrentViewOnTop(g, "StorageClasses")
			return err
		}
	} else if v == nil || v.Name() == "StorageClasses" {
		_, err := setCurrentViewOnTop(g, "Service")
		return err
	} else if v == nil || v.Name() == "Service" {
		if _, err := g.View("Ingress"); err != nil {
			if err != nil {
				if _, err := g.View("Pvc"); err != nil {
					if err != nil {
						_, err := setCurrentViewOnTop(g, "Configmap")
						return err
					}
				} else {
					_, err := setCurrentViewOnTop(g, "Pvc")
					return err
				}
			}
		} else {
			_, err := setCurrentViewOnTop(g, "Ingress")
			return err
		}
	} else if v == nil || v.Name() == "Ingress" {
		if _, err := g.View("Pvc"); err != nil {
			if err != nil {
				_, err := setCurrentViewOnTop(g, "Configmap")
				return err
			}
		} else {
			_, err := setCurrentViewOnTop(g, "Pvc")
			return err
		}
	} else if v == nil || v.Name() == "Pvc" {
		_, err := setCurrentViewOnTop(g, "Configmap")
		return err
	} else if v == nil || v.Name() == "Configmap" {
		_, err := setCurrentViewOnTop(g, "Secrets")
		return err
	} else if v == nil || v.Name() == "Secrets" {
		_, err := setCurrentViewOnTop(g, "bottom")
		return err
	} else if v == nil || v.Name() == "pod" {
		_, err := setCurrentViewOnTop(g, "Namespace")
		return err
	}
	return nil
}