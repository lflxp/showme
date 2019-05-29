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

	if v, err := g.SetView("bottom", 0, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "工作负载状态"
		v.Wrap = true
		// v.Highlight = true
		// v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// uri = fmt.Sprintf("/a%s", time.Now().Format("150405"))
		fmt.Fprintln(v, fmt.Sprintf("URL => 0.0.0.0:%s <= \nPATH: => %s <=", "9999", "/tmp"))
		fmt.Fprintln(v, origin.Cluster)
		if _, err := setCurrentViewOnTop(g, "bottom"); err != nil {
			return err
		}
	}

	return nil
}

func ServiceConfigStatusTable(g *gocui.Gui, startx, starty, endx, endy int, data ClusterStatus) error {
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
