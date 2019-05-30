package kubectl

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils/table"
)

func KeyPod(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, Pods); err != nil {
		return err
	}
	return nil
}

func Pods(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.SetView("Pod", 0, 0, origin.maxX-1, origin.maxY-1); err != nil {
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
		tableNow.AddCol("NAME").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Namespace").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Node").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Ready").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Restarts").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Time").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
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
			ns.Data = fmt.Sprintf("%s", value.Namespace)
			ns.TextAlign = table.TextCenter
			ns.Color = "yellow"
			tableNow.AddRow(1, ns)

			node := table.NewCol()
			node.Data = fmt.Sprintf("%s", value.Node)
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
