package kubectl

import (
	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils/k8s"
	log "github.com/sirupsen/logrus"
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

	// init gocui

	origin.Gui.Highlight = true
	origin.Gui.Cursor = true
	origin.Gui.SelFgColor = gocui.ColorGreen
	origin.Gui.SetManagerFunc(dashboard)

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

// Global Values
type BasicKubectl struct {
	// gocui
	Gui *gocui.Gui
	// kubectl
	ClientSet *kubernetes.Clientset
	DefaultNS string   // current namespace
	Helps     []string // F1 View show help message
	Cluster   []ClusterStatus
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
		_, err := setCurrentViewOnTop(g, "Namespace")

		return err
	} else if v == nil || v.Name() == "Namespace" {
		_, err := setCurrentViewOnTop(g, "Node")
		return err
	} else if v == nil || v.Name() == "Node" {
		_, err := setCurrentViewOnTop(g, "Role")
		return err
	} else if v == nil || v.Name() == "Role" {
		_, err := setCurrentViewOnTop(g, "bottom")
		return err
	}
	return nil
}
