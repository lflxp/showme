package kubectl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
	"github.com/lflxp/showme/utils/table"
	k8stype "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KeyNode(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyF6, gocui.ModNone, Node); err != nil {
		return err
	}
	if err := g.SetKeybinding("Noded", gocui.KeyEnter, gocui.ModNone, getNodes); err != nil {
		return err
	}
	if err := g.SetKeybinding("msgnode", gocui.KeyEnter, gocui.ModNone, delnodemessage); err != nil {
		return err
	}
	if err := g.SetKeybinding("Noded", gocui.KeyDelete, gocui.ModNone, deletenodeView); err != nil {
		return err
	}
	if err := g.SetKeybinding("delnodeView", gocui.KeyEnter, gocui.ModNone, nextView); err != nil {
		return err
	}
	return nil
}

func deletenodeView(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	rs := strings.Split(strings.Replace(l, ">", "*", 1), "*")
	maxX, maxY := g.Size()
	if v, err := g.SetView("delnodeView", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		origin.CurrentNode = strings.TrimSpace(rs[1])
		v.Title = fmt.Sprintf("确认SchedulingDisabled[Node] %s?(y/N)", origin.CurrentNode)
		v.Highlight = true
		v.Editable = true
		// v.Frame = false
		// v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))
		if _, err := g.SetCurrentView("delnodeView"); err != nil {
			return err
		}
	}
	return nil
}

func delnodemessage(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msgnode"); err != nil {
		return err
	}
	if _, err := setCurrentViewOnTop(g, "Noded"); err != nil {
		return err
	}
	return nil
}

func getNodes(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	ttt := strings.Split(strings.Replace(l, ">", "*", 1), "*")
	if len(ttt) > 1 {
		maxX, maxY := g.Size()
		if v, err := g.SetView("msgnode", maxX*8/100, maxY*8/100, maxX*92/100, maxY*92/100); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			name := strings.TrimSpace(ttt[1])
			v.Title = fmt.Sprintf("Current Node: %s", name)
			v.Highlight = true
			v.SelFgColor = gocui.ColorMagenta
			v.SelBgColor = gocui.ColorCyan
			v.Editable = true
			v.Wrap = true
			// fmt.Fprintln(v, strings.Trim(l, " "))
			// fmt.Fprintln(v, l)
			// selectId = strings.Trim(l, " ")

			nn, err := origin.ClientSet.CoreV1().Nodes().Get(name, metav1.GetOptions{})
			if err != nil {
				fmt.Fprintln(v, err.Error())
			} else {
				// json格式美化
				b, err := json.MarshalIndent(nn, "", "\t")
				if err != nil {
					fmt.Fprintln(v, err.Error())
				} else {
					fmt.Fprintln(v, utils.Colorize(string(b), "red", "", false, true))
				}
			}

			if _, err := g.SetCurrentView("msgnode"); err != nil {
				return err
			}

		}
	}

	return nil
}

func Node(g *gocui.Gui, v *gocui.View) error {
	if err = delOtherViewNoBack(g); err != nil {
		return err
	}
	if v, err := g.SetView("Noded", 0, 0, origin.maxX-1, origin.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		node_list, err := origin.ClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			return err
		}

		v.Title = fmt.Sprintf("Node/%d", len(node_list.Items))
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		// v.MoveCursor(startx, endy, false)
		if _, err := setCurrentViewOnTop(g, "Noded"); err != nil {
			return err
		}

		tableNow := table.NewTable(origin.maxX - 1)

		// tableNow.AddCol("ID").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("NAME").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("STATUS").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("ROLES").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("AGE").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("VERSION").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("INTERNAL-IP").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("OS-IMAGE").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("KERNEL-VERSION").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("CONTAINER-RUNTIME").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("LABELS").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.CalColumnWidths()

		for n, value := range node_list.Items {
			if n == 0 {
				tableNow.FprintHeader(v)
			}

			name := table.NewCol()
			name.Data = fmt.Sprintf("*%s", value.GetName())
			name.TextAlign = table.TextLeft
			name.Color = "red"
			tableNow.AddRow(0, name)

			var status string
			for _, x := range value.Status.Conditions {
				if x.Type == k8stype.NodeReady {
					if x.Status == k8stype.ConditionTrue {
						status = "Ready"
					} else {
						status = x.Reason
					}
				}
			}
			ns := table.NewCol()
			ns.Data = fmt.Sprintf("*%s", status)
			ns.TextAlign = table.TextLeft
			ns.Color = "yellow"
			tableNow.AddRow(1, ns)

			roles := table.NewCol()
			roles.Data = fmt.Sprintf("*%s", value.ObjectMeta.Labels["kubernetes.io/role"])
			roles.TextAlign = table.TextLeft
			roles.Color = "yellow"
			tableNow.AddRow(2, roles)

			ages := table.NewCol()
			ages.Data = fmt.Sprintf("%s", strings.Replace(fmt.Sprintf("%v", value.ObjectMeta.CreationTimestamp.Sub(time.Now())), "-", "", -1))
			ages.TextAlign = table.TextLeft
			ages.Color = "yellow"
			tableNow.AddRow(3, ages)

			version := table.NewCol()
			version.Data = value.Status.NodeInfo.KubeProxyVersion
			version.TextAlign = table.TextLeft
			version.Color = "yellow"
			tableNow.AddRow(4, version)

			inip := ""
			for _, x := range value.Status.Addresses {
				if x.Type == k8stype.NodeInternalIP {
					inip = x.Address
				}
			}
			innerip := table.NewCol()
			innerip.Data = inip
			innerip.TextAlign = table.TextLeft
			innerip.Color = "yellow"
			tableNow.AddRow(5, innerip)

			osimage := table.NewCol()
			osimage.Data = value.Status.NodeInfo.OSImage
			osimage.TextAlign = table.TextLeft
			osimage.Color = "yellow"
			tableNow.AddRow(6, osimage)

			kernel := table.NewCol()
			kernel.Data = value.Status.NodeInfo.KernelVersion
			kernel.TextAlign = table.TextLeft
			kernel.Color = "yellow"
			tableNow.AddRow(7, kernel)

			containerrun := table.NewCol()
			containerrun.Data = value.Status.NodeInfo.ContainerRuntimeVersion
			containerrun.TextAlign = table.TextLeft
			containerrun.Color = "yellow"
			tableNow.AddRow(8, containerrun)

			lb := []string{}
			for k, v := range value.ObjectMeta.Labels {
				lb = append(lb, fmt.Sprintf("%s:%s", k, v))
			}

			labels := table.NewCol()
			labels.Data = strings.Join(lb, ",")
			labels.TextAlign = table.TextLeft
			labels.Color = "yellow"
			tableNow.AddRow(9, labels)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)
	}
	return nil
}
