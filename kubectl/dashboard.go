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
	if err := g.SetKeybinding("", gocui.KeyCtrl2, gocui.ModNone, ok); err != nil {
		return err
	}
	if err := g.SetKeybinding("ok", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		return err
	}
	return nil
}

func dquit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func dashboard(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("dashboard", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dashboard"
		// v.Frame = true
		// v.Highlight = true
		// v.Editable = true
		if _, err := g.SetCurrentView("dashboard"); err != nil {
			return err
		}
	}
	return nil
}

func ok(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("dashboard"); err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("ok", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Input Port Range(eg: 80,3306,25-100)"
		v.Highlight = true
		v.Editable = true
		// v.Frame = false
		// v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))
		fmt.Fprintln(v, "11111111111111111111111")
		if _, err := g.SetCurrentView("ok"); err != nil {
			return err
		}
	}

	return nil
}
