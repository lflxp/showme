package kubectl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
	"github.com/lflxp/showme/utils/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KeyDeployment(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyF4, gocui.ModNone, Deployment); err != nil {
		return err
	}
	if err := g.SetKeybinding("Deployment", gocui.KeyEnter, gocui.ModNone, getDeployments); err != nil {
		return err
	}
	if err := g.SetKeybinding("msgdeploy", gocui.KeyEnter, gocui.ModNone, deldeploymentmessage); err != nil {
		return err
	}
	if err := g.SetKeybinding("Deployment", gocui.KeyDelete, gocui.ModNone, deleteDeploymentView); err != nil {
		return err
	}
	if err := g.SetKeybinding("deldeployment", gocui.KeyEnter, gocui.ModNone, nextView); err != nil {
		return err
	}
	return nil
}

func deleteDeploymentView(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	rs := strings.Split(strings.Replace(l, ">", "*", 1), "*")
	maxX, maxY := g.Size()
	if v, err := g.SetView("deldeployment", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		origin.CurrentPod = strings.TrimSpace(rs[1])
		origin.DefaultNS = strings.TrimSpace(rs[2])
		v.Title = fmt.Sprintf("确认删除[Deploy] %s:%s?(y/N)", strings.TrimSpace(rs[2]), strings.TrimSpace(rs[1]))
		v.Highlight = true
		v.Editable = true
		// v.Frame = false
		// v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))
		if _, err := g.SetCurrentView("deldeployment"); err != nil {
			return err
		}
	}
	return nil
}

func deldeploymentmessage(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msgdeploy"); err != nil {
		return err
	}
	if _, err := setCurrentViewOnTop(g, "Deployment"); err != nil {
		return err
	}
	return nil
}

func getDeployments(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	ttt := strings.Split(strings.Replace(l, ">", "*", 1), "*")
	if len(ttt) > 1 {
		maxX, maxY := g.Size()
		if v, err := g.SetView("msgdeploy", maxX*8/100, maxY*8/100, maxX*92/100, maxY*92/100); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			namespace := strings.TrimSpace(ttt[2])
			name := strings.TrimSpace(ttt[1])
			v.Title = fmt.Sprintf("Current: %s %s", namespace, name)
			v.Highlight = true
			v.SelFgColor = gocui.ColorMagenta
			v.SelBgColor = gocui.ColorCyan
			v.Editable = true
			v.Wrap = true
			// fmt.Fprintln(v, strings.Trim(l, " "))
			// fmt.Fprintln(v, l)
			// selectId = strings.Trim(l, " ")

			pod, err := origin.ClientSet.Extensions().Deployments(namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				fmt.Fprintln(v, err.Error())
			} else {
				// json格式美化
				b, err := json.MarshalIndent(pod, "", "\t")
				if err != nil {
					fmt.Fprintln(v, err.Error())
				} else {
					fmt.Fprintln(v, utils.Colorize(string(b), "red", "", false, true))
				}
			}

			if _, err := g.SetCurrentView("msgdeploy"); err != nil {
				return err
			}

		}
	}

	return nil
}

func Deployment(g *gocui.Gui, v *gocui.View) error {
	if err = delOtherViewNoBack(g); err != nil {
		return err
	}
	if v, err := g.SetView("Deployment", 0, 0, origin.maxX-1, origin.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		deploy_list, err := origin.ClientSet.Extensions().Deployments("").List(metav1.ListOptions{})
		if err != nil {
			return err
		}

		v.Title = fmt.Sprintf("Deployment/%d", len(deploy_list.Items))
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		// v.MoveCursor(startx, endy, false)
		if _, err := setCurrentViewOnTop(g, "Deployment"); err != nil {
			return err
		}

		tableNow := table.NewTable(origin.maxX - 2)

		// tableNow.AddCol("ID").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("NAME").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Namespace").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Tags").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Ready").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Images").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Time").SetColor("dgreen").SetTextAlign(table.TextRight).SetBgColor("black")
		tableNow.CalColumnWidths()

		for n, value := range deploy_list.Items {
			if n == 0 {
				tableNow.FprintHeader(v)
			}

			name := table.NewCol()
			name.Data = fmt.Sprintf("*%s", value.GetName())
			name.TextAlign = table.TextLeft
			name.Color = "yellow"
			tableNow.AddRow(0, name)

			ns := table.NewCol()
			ns.Data = fmt.Sprintf("*%s", value.GetNamespace())
			ns.TextAlign = table.TextCenter
			ns.Color = "yellow"
			tableNow.AddRow(1, ns)

			ttags := ""
			for k, v := range value.Labels {
				ttags += fmt.Sprintf("%s:%s ", k, v)
			}
			Tags := table.NewCol()
			Tags.Data = fmt.Sprintf("*%s", ttags)
			Tags.TextAlign = table.TextLeft
			Tags.Color = "yellow"
			tableNow.AddRow(2, Tags)

			rd := table.NewCol()
			rd.Data = fmt.Sprintf("%d / %d", value.Status.AvailableReplicas, value.Status.Replicas)
			rd.TextAlign = table.TextCenter
			rd.Color = "yellow"
			tableNow.AddRow(3, rd)

			image := table.NewCol()
			image.Data = fmt.Sprintf("%s", value.Spec.Template.Spec.Containers[0].Image)
			image.TextAlign = table.TextLeft
			image.Color = "yellow"
			tableNow.AddRow(4, image)

			timed := table.NewCol()
			timed.Data = strings.Replace(fmt.Sprintf("%v", value.CreationTimestamp.Sub(time.Now())), "-", "", -1)
			timed.TextAlign = table.TextRight
			timed.Color = "yellow"
			tableNow.AddRow(5, timed)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)
	}
	return nil
}
