package kubectl

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
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
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delpodmessage); err != nil {
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
	if v, err := g.View("Pod"); err == nil {
		v.Clear()
		v.Title = fmt.Sprintf("Pod/%d", len(origin.Pods))
		v.Highlight = true
		v.Editable = true

		num := 0
		tableNow := table.NewTable(origin.maxX - 1)

		// tableNow.AddCol("ID").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("NAME").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Namespace").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Node").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Ready").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Restarts").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Time").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.CalColumnWidths()

		for _, value := range origin.Pods {
			num++
			if num == 1 {
				tableNow.FprintHeader(v)
			}

			name := table.NewCol()
			name.Data = fmt.Sprintf("*%s", value.Name)
			name.TextAlign = table.TextLeft
			name.Color = "yellow"
			tableNow.AddRow(0, name)

			ns := table.NewCol()
			ns.Data = fmt.Sprintf("*%s", value.Namespace)
			ns.TextAlign = table.TextCenter
			ns.Color = "yellow"
			tableNow.AddRow(1, ns)

			node := table.NewCol()
			node.Data = fmt.Sprintf("*%s", value.Node)
			node.TextAlign = table.TextCenter
			node.Color = "yellow"
			tableNow.AddRow(2, node)

			rd := table.NewCol()
			rd.Data = fmt.Sprintf("%s", value.Ready)
			rd.TextAlign = table.TextCenter
			rd.Color = "yellow"
			tableNow.AddRow(3, rd)

			rs := table.NewCol()
			rs.Data = fmt.Sprintf("%s", value.Restarts)
			rs.TextAlign = table.TextCenter
			rs.Color = "yellow"
			tableNow.AddRow(4, rs)

			time := table.NewCol()
			time.Data = fmt.Sprintf("%s", value.Time)
			time.TextAlign = table.TextLeft
			time.Color = "yellow"
			tableNow.AddRow(5, time)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)
	} else if v, err := g.SetView("Pod", 0, 0, origin.maxX-1, origin.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("Pod/%d", len(origin.Pods))
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		// v.MoveCursor(startx, endy, false)
		if _, err := setCurrentViewOnTop(g, "Pod"); err != nil {
			return err
		}

		num := 0
		tableNow := table.NewTable(origin.maxX - 1)

		// tableNow.AddCol("ID").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("NAME").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Namespace").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Node").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Ready").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Restarts").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Time").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.CalColumnWidths()

		for _, value := range origin.Pods {
			num++
			if num == 1 {
				tableNow.FprintHeader(v)
			}

			name := table.NewCol()
			name.Data = fmt.Sprintf("*%s", value.Name)
			name.TextAlign = table.TextLeft
			name.Color = "yellow"
			tableNow.AddRow(0, name)

			ns := table.NewCol()
			ns.Data = fmt.Sprintf("*%s", value.Namespace)
			ns.TextAlign = table.TextCenter
			ns.Color = "yellow"
			tableNow.AddRow(1, ns)

			node := table.NewCol()
			node.Data = fmt.Sprintf("*%s", value.Node)
			node.TextAlign = table.TextCenter
			node.Color = "yellow"
			tableNow.AddRow(2, node)

			rd := table.NewCol()
			rd.Data = fmt.Sprintf("%s", value.Ready)
			rd.TextAlign = table.TextCenter
			rd.Color = "yellow"
			tableNow.AddRow(3, rd)

			rs := table.NewCol()
			rs.Data = fmt.Sprintf("%s", value.Restarts)
			rs.TextAlign = table.TextCenter
			rs.Color = "yellow"
			tableNow.AddRow(4, rs)

			time := table.NewCol()
			time.Data = fmt.Sprintf("%s", value.Time)
			time.TextAlign = table.TextLeft
			time.Color = "yellow"
			tableNow.AddRow(5, time)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)
	}
	return nil
}
