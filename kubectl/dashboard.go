package kubectl

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils/table"
)

// keybinding

func KeyDashboard(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, searchBuffer); err != nil {
		return err
	}
	if err := g.SetKeybinding("searchBuffer", gocui.KeyEnter, gocui.ModNone, delsearchBuffer); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	return nil
}

func dashboard(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if len(origin.Cluster) > 0 {
		num := len(origin.Cluster)
		len := maxX / num
		for n, x := range origin.Cluster {
			var endX int
			startX := n * len
			startY := 0
			if n == num-1 {
				endX = maxX - 1
			} else {
				endX = (n+1)*len - 1
			}
			// endX = (n + 1) * len
			endY := maxY/4 - 1

			err := StatusTable(g, startX, startY, endX, endY, x)
			if err != nil {
				return err
			}
		}
	}

	if len(origin.ServiceConfig) > 0 {
		num := len(origin.Cluster)
		len := maxX / num
		for n, x := range origin.ServiceConfig {
			var endX int
			startX := n * len
			startY := maxY / 4
			if n == num-1 {
				endX = maxX - 1
			} else {
				endX = (n+1)*len - 1
			}
			// endX = (n + 1) * len
			endY := maxY/2 - 1

			err := StatusTable(g, startX, startY, endX, endY, x)
			if err != nil {
				return err
			}
		}
	}

	// if v, err := g.SetView("bottom", 0, maxY/2, maxX/2-1, maxY-1); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	v.Title = "工作负载状态"
	// 	v.Wrap = true
	// 	// v.Highlight = true
	// 	// v.Autoscroll = true
	// 	v.SelBgColor = gocui.ColorGreen
	// 	v.SelFgColor = gocui.ColorBlack
	// 	// v.Editable = true
	// 	// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
	// 	// uri = fmt.Sprintf("/a%s", time.Now().Format("150405"))
	// 	fmt.Fprintln(v, fmt.Sprintf("URL => 0.0.0.0:%s <= \nPATH: => %s <=", "9999", "/tmp"))
	// 	fmt.Fprintln(v, origin.Cluster)
	// 	if _, err := setCurrentViewOnTop(g, "bottom"); err != nil {
	// 		return err
	// 	}
	// }

	// if err := WorkLoadTable(g, 0, maxY/2, maxX/2-1, maxY-1); err != nil {
	if err := WorkLoadTable(g, 0, maxY/2, maxX-1, maxY*3/4-1); err != nil {
		return err
	}

	// if v, err := g.SetView("pod", maxX/2, maxY/2, maxX-1, maxY-1); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	v.Title = "Pod"
	// 	v.Wrap = true
	// 	// v.Highlight = true
	// 	// v.Autoscroll = true
	// 	v.SelBgColor = gocui.ColorGreen
	// 	v.SelFgColor = gocui.ColorBlack
	// 	// v.Editable = true
	// 	// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
	// 	// uri = fmt.Sprintf("/a%s", time.Now().Format("150405"))
	// 	fmt.Fprintln(v, fmt.Sprintf("URL => 0.0.0.0:%s <= \nPATH: => %s <=", "9999", "/tmp"))
	// 	fmt.Fprintln(v, origin.Cluster)
	// }

	// if err := PodsTable(g, maxX/2, maxY/2, maxX-1, maxY-1); err != nil {
	if err := PodsTable(g, 0, maxY*3/4, maxX-1, maxY-1); err != nil {
		return err
	}

	return nil
}

func RefreshWorkLoad(g *gocui.Gui, startx, starty, endx, endy int) error {
	if v, err := origin.Gui.View("bottom"); err != nil {
		return err
	} else {
		v.Clear()
		v.Autoscroll = true
		v.Wrap = true
		v.Highlight = true
		// v.Editable = true
		num := 0
		tableNow := table.NewTable(endx - startx)

		// tableNow.AddCol("ID").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Type").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("NAME").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Namespace").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Tags").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Ready").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Images").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Time").SetColor("red").SetTextAlign(table.TextRight).SetBgColor("black")
		tableNow.CalColumnWidths()
		for _, value := range origin.PodControllers {
			num++
			if num == 1 {
				tableNow.FprintHeader(v)
			}

			// id := table.NewCol()
			// id.Data = fmt.Sprintf("%d", num)
			// id.TextAlign = table.TextCenter
			// id.Color = "yellow"
			// tableNow.AddRow(0, id)

			tt := table.NewCol()
			tt.Data = fmt.Sprintf("*%s", value.Type)
			tt.TextAlign = table.TextLeft
			tt.Color = "yellow"
			tableNow.AddRow(0, tt)

			name := table.NewCol()
			name.Data = fmt.Sprintf("%s", value.Name)
			name.TextAlign = table.TextLeft
			name.Color = "yellow"
			tableNow.AddRow(1, name)

			ns := table.NewCol()
			ns.Data = fmt.Sprintf("%s", value.Namespace)
			ns.TextAlign = table.TextCenter
			ns.Color = "yellow"
			tableNow.AddRow(2, ns)

			ttags := ""
			for k, v := range value.Tags {
				ttags += fmt.Sprintf("%s:%s ", k, v)
			}
			Tags := table.NewCol()
			Tags.Data = fmt.Sprintf("%s", ttags)
			Tags.TextAlign = table.TextLeft
			Tags.Color = "yellow"
			tableNow.AddRow(3, Tags)

			rd := table.NewCol()
			rd.Data = fmt.Sprintf("%s", value.ContainerGroup)
			rd.TextAlign = table.TextCenter
			rd.Color = "yellow"
			tableNow.AddRow(4, rd)

			image := table.NewCol()
			image.Data = fmt.Sprintf("%s", value.Images)
			image.TextAlign = table.TextLeft
			image.Color = "yellow"
			tableNow.AddRow(5, image)

			time := table.NewCol()
			time.Data = fmt.Sprintf("%s", value.Time)
			time.TextAlign = table.TextRight
			time.Color = "yellow"
			tableNow.AddRow(6, time)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)
	}
	return nil
}

func WorkLoadTable(g *gocui.Gui, startx, starty, endx, endy int) error {
	if v, err := g.SetView("bottom", startx, starty, endx, endy); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("WorkLoad/%d", len(origin.PodControllers))
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		v.MoveCursor(startx, endy, false)

		v.Clear()
		num := 0
		tableNow := table.NewTable(endx - startx)

		// tableNow.AddCol("ID").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Type").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("NAME").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Namespace").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Tags").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Ready").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.AddCol("Images").SetColor("red").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("Time").SetColor("red").SetTextAlign(table.TextRight).SetBgColor("black")
		tableNow.CalColumnWidths()
		for _, value := range origin.PodControllers {
			num++
			if num == 1 {
				tableNow.FprintHeader(v)
			}

			// id := table.NewCol()
			// id.Data = fmt.Sprintf("%d", num)
			// id.TextAlign = table.TextCenter
			// id.Color = "yellow"
			// tableNow.AddRow(0, id)

			tt := table.NewCol()
			tt.Data = fmt.Sprintf("*%s", value.Type)
			tt.TextAlign = table.TextLeft
			tt.Color = "yellow"
			tableNow.AddRow(0, tt)

			name := table.NewCol()
			name.Data = fmt.Sprintf("%s", value.Name)
			name.TextAlign = table.TextLeft
			name.Color = "yellow"
			tableNow.AddRow(1, name)

			ns := table.NewCol()
			ns.Data = fmt.Sprintf("%s", value.Namespace)
			ns.TextAlign = table.TextCenter
			ns.Color = "yellow"
			tableNow.AddRow(2, ns)

			ttags := ""
			for k, v := range value.Tags {
				ttags += fmt.Sprintf("%s:%s ", k, v)
			}
			Tags := table.NewCol()
			Tags.Data = fmt.Sprintf("%s", ttags)
			Tags.TextAlign = table.TextLeft
			Tags.Color = "yellow"
			tableNow.AddRow(3, Tags)

			rd := table.NewCol()
			rd.Data = fmt.Sprintf("%s", value.ContainerGroup)
			rd.TextAlign = table.TextCenter
			rd.Color = "yellow"
			tableNow.AddRow(4, rd)

			image := table.NewCol()
			image.Data = fmt.Sprintf("%s", value.Images)
			image.TextAlign = table.TextLeft
			image.Color = "yellow"
			tableNow.AddRow(5, image)

			time := table.NewCol()
			time.Data = fmt.Sprintf("%s", value.Time)
			time.TextAlign = table.TextRight
			time.Color = "yellow"
			tableNow.AddRow(6, time)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)

		if _, err := setCurrentViewOnTop(g, "bottom"); err != nil {
			return err
		}
	}
	return nil
}

func RefreshPods(g *gocui.Gui, startx, starty, endx, endy int) error {
	if v, err := origin.Gui.View("pod"); err != nil {
		return err
	} else {
		v.Clear()
		v.Autoscroll = true
		v.Wrap = true
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		v.MoveCursor(startx, endy, false)

		num := 0
		tableNow := table.NewTable(endx - startx)

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

			// id := table.NewCol()
			// id.Data = fmt.Sprintf("%d", num)
			// id.TextAlign = table.TextCenter
			// id.Color = "yellow"
			// tableNow.AddRow(0, id)

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

func PodsTable(g *gocui.Gui, startx, starty, endx, endy int) error {
	if v, err := g.SetView("pod", startx, starty, endx, endy); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("Pod/%d", len(origin.Pods))
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		v.MoveCursor(startx, endy, false)

		num := 0
		tableNow := table.NewTable(endx - startx)

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

			// id := table.NewCol()
			// id.Data = fmt.Sprintf("%d", num)
			// id.TextAlign = table.TextCenter
			// id.Color = "yellow"
			// tableNow.AddRow(0, id)

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

func StatusTable(g *gocui.Gui, startx, starty, endx, endy int, data ClusterStatus) error {
	if v, err := g.SetView(data.Title, startx, starty, endx, endy); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("%s/%d", data.Title, data.Count)
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		v.MoveCursor(startx, endy, false)

		num := 0
		tableNow := table.NewTable(endx - startx)

		tableNow.AddCol("NAME").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
		tableNow.CalColumnWidths()

		for _, value := range data.Data {
			num++
			if num == 1 {
				tableNow.FprintHeader(v)
			}

			// id := table.NewCol()
			// id.Data = fmt.Sprintf("%d", num)
			// id.TextAlign = table.TextCenter
			// id.Color = "yellow"
			// tableNow.AddRow(0, id)

			name := table.NewCol()
			name.Data = fmt.Sprintf("*%s", value)
			name.TextAlign = table.TextLeft
			name.Color = "yellow"
			tableNow.AddRow(0, name)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)
	}
	return nil
}
