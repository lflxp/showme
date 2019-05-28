package kubectl

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// keybinding

func KeyDashboard(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		return err
	}
	if err := g.SetKeybinding("top", gocui.KeyEnter, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("bottom", gocui.KeyEnter, gocui.ModNone, nextView); err != nil {
		return err
	}
	return nil
}

func dashboard(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("top", 0, 0, maxX-1, maxY/4-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dashboard"
		// v.Frame = true
		v.Highlight = true
		v.Editable = true
		fmt.Fprintln(v, "Hello")
		fmt.Fprintln(v, "world")
		if _, err := setCurrentViewOnTop(g, "top"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("bottom", 0, maxY/4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "静态服务器地址"
		v.Wrap = true
		// v.Highlight = true
		// v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// uri = fmt.Sprintf("/a%s", time.Now().Format("150405"))
		fmt.Fprintln(v, fmt.Sprintf("URL => 0.0.0.0:%s <= \nPATH: => %s <=", "9999", "/tmp"))
	}
	// if v, err := g.SetView("oook", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}

	// 	v.Title = "Input Port Range(eg: 80,3306,25-100)"
	// 	v.Highlight = true
	// 	v.Editable = true
	// 	v.Frame = true
	// 	v.SelBgColor = gocui.ColorYellow
	// 	v.SelFgColor = gocui.ColorRed
	// 	// fmt.Fprintln(v, strings.Trim(l, " "))
	// 	// fmt.Fprintln(v, l)
	// 	// selectId = strings.Trim(l, " ")
	// 	// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))
	// 	fmt.Fprintln(v, "11111111111111111111111")
	// 	if _, err := g.SetViewOnBottom("oook"); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
