package kubectl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
	"github.com/lflxp/showme/utils/k8s"
	"github.com/lflxp/showme/utils/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KeyPod(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyF3, gocui.ModNone, Pods); err != nil {
		return err
	}
	if err := g.SetKeybinding("Pod", gocui.KeyEnter, gocui.ModNone, getPods); err != nil {
		return err
	}
	if err := g.SetKeybinding("Pod", gocui.KeyCtrlL, gocui.ModNone, getPodLogs); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delpodmessage); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyCtrlL, gocui.ModNone, getPodLogs); err != nil {
		return err
	}
	if err := g.SetKeybinding("Pod", gocui.KeyDelete, gocui.ModNone, deletePodView); err != nil {
		return err
	}
	if err := g.SetKeybinding("delpod", gocui.KeyEnter, gocui.ModNone, nextView); err != nil {
		return err
	}
	return nil
}

func getPodLogs(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "msg" {
		v.Clear()
		v.Title = fmt.Sprintf("Current Logs R: %s %s", origin.DefaultNS, origin.CurrentPod)
		v.Autoscroll = false
		podstring, err := k8s.GetPodLogByPodIdByNum(origin.DefaultNS, origin.CurrentPod, 200)
		if err != nil {
			panic(err)
		} else {
			// json格式美化
			fmt.Fprintln(v, podstring)
		}
	} else {
		var l string
		var err error

		_, cy := v.Cursor()
		if l, err = v.Line(cy); err != nil {
			l = ""
		}

		ttt := strings.Split(strings.Replace(l, ">", "*", 1), "*")
		if len(ttt) > 1 {
			maxX, maxY := g.Size()
			if v, err := g.SetView("msg", maxX*8/100, maxY*8/100, maxX*92/100, maxY*92/100); err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}

				origin.DefaultNS = strings.TrimSpace(ttt[2])
				origin.CurrentPod = strings.TrimSpace(ttt[1])
				v.Title = fmt.Sprintf("Current Logs: %s %s", origin.DefaultNS, origin.CurrentPod)
				v.Highlight = true
				v.SelFgColor = gocui.ColorWhite
				v.SelBgColor = gocui.ColorCyan
				v.Editable = true
				v.Wrap = true
				v.Autoscroll = true
				// fmt.Fprintln(v, strings.Trim(l, " "))
				// fmt.Fprintln(v, l)
				// selectId = strings.Trim(l, " ")

				podstring, err := k8s.GetPodLogByPodIdByNum(origin.DefaultNS, origin.CurrentPod, 200)
				if err != nil {
					// panic(err)
					fmt.Fprintln(v, utils.Colorize(err.Error(), "green", "", false, true))
				} else {
					// json格式美化
					fmt.Fprintln(v, podstring)
				}

				if _, err := g.SetCurrentView("msg"); err != nil {
					return err
				}

			}
		}
	}

	return nil
}

func deletePodView(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	rs := strings.Split(strings.Replace(l, ">", "*", 1), "*")
	maxX, maxY := g.Size()
	if v, err := g.SetView("delpod", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		origin.CurrentPod = strings.TrimSpace(rs[1])
		origin.DefaultNS = strings.TrimSpace(rs[2])
		v.Title = fmt.Sprintf("确认删除[POD] %s:%s?(y/N)", strings.TrimSpace(rs[2]), strings.TrimSpace(rs[1]))
		v.Highlight = true
		v.Editable = true
		// v.Frame = false
		// v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))
		if _, err := g.SetCurrentView("delpod"); err != nil {
			return err
		}
	}
	return nil
}

func delpodmessage(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := setCurrentViewOnTop(g, "Pod"); err != nil {
		return err
	}
	return nil
}

func getPods(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	ttt := strings.Split(strings.Replace(l, ">", "*", 1), "*")
	if len(ttt) > 1 {
		maxX, maxY := g.Size()
		if v, err := g.SetView("msg", maxX*8/100, maxY*8/100, maxX*92/100, maxY*92/100); err != nil {
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

			pod, err := origin.ClientSet.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
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

			if _, err := g.SetCurrentView("msg"); err != nil {
				return err
			}

		}
	}

	return nil
}

func Pods(g *gocui.Gui, v *gocui.View) error {
	if err = delOtherViewNoBack(g); err != nil {
		return err
	}
	if v, err := g.SetView("Pod", 0, 0, origin.maxX-1, origin.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		pod_list, err := origin.ClientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		v.Title = fmt.Sprintf("Pod/%d", len(pod_list.Items))
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		// v.MoveCursor(startx, endy, false)
		if _, err := setCurrentViewOnTop(g, "Pod"); err != nil {
			return err
		}

		tableNow := table.NewTable(origin.maxX - 1)

		// tableNow.AddCol("ID").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("NAME").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Namespace").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Node").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Ready").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Restarts").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Time").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.CalColumnWidths()

		for n, value := range pod_list.Items {
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

			node := table.NewCol()
			node.Data = fmt.Sprintf("*%s", value.Spec.NodeName)
			node.TextAlign = table.TextCenter
			node.Color = "yellow"
			tableNow.AddRow(2, node)

			rd := table.NewCol()
			rd.Data = fmt.Sprintf("%s", value.Status.Phase)
			rd.TextAlign = table.TextCenter
			rd.Color = "yellow"
			tableNow.AddRow(3, rd)

			rs := table.NewCol()
			rs.Data = fmt.Sprintf("%s", fmt.Sprintf("%d", value.Status.ContainerStatuses[0].RestartCount))
			rs.TextAlign = table.TextCenter
			rs.Color = "yellow"
			tableNow.AddRow(4, rs)

			timed := table.NewCol()
			timed.Data = fmt.Sprintf("%s", strings.Replace(fmt.Sprintf("%v", value.Status.StartTime.Sub(time.Now())), "-", "", -1))
			timed.TextAlign = table.TextLeft
			timed.Color = "yellow"
			tableNow.AddRow(5, timed)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)
	}
	return nil
}
