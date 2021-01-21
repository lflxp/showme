package pkg

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func delinputIp(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("inputip"); err != nil {
		return err
	}

	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	ips = l
	maxX, _ := g.Size()
	// go ScanIp(v, maxX-4)

	// if _, err = setCurrentViewOnTop(g, "top"); err != nil {
	if vivia, err := g.SetCurrentView("top"); err != nil {
		return err
	} else {
		vivia.Highlight = true
		vivia.Clear()
		go ScanIp(g, vivia, maxX/2-3)
	}

	return nil
}

func delPort(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("port"); err != nil {
		return err
	}

	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	port = l
	maxX, _ := g.Size()
	if v, err := g.SetCurrentView("scanport"); err != nil {
		return err
	} else {
		v.Title = fmt.Sprintf("Current IP:   %s     PORTS RANGE:  %s", selectId, l)
		v.Highlight = true
		// v.Autoscroll = true
		v.Clear()
		go ScanIpPorts(g, v, maxX/2)
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "top" {
		_, err := g.SetCurrentView("scanport")

		return err
	} else if v == nil || v.Name() == "scanport" {
		_, err := g.SetCurrentView("help")
		return err
	} else if v == nil || v.Name() == "help" {
		_, err := g.SetCurrentView("top")
		return err
	} else if v == nil || v.Name() == "gethelp" {
		fmt.Fprintln(v, "Esc push")
		if err := g.DeleteView("gethelp"); err != nil {
			return err
		}

		if _, err := setCurrentViewOnTop(g, "top"); err != nil {
			return err
		}
	} else if v == nil || v.Name() == "help" {
		_, err := g.SetCurrentView("top")

		return err
	}

	// _, err := g.SetCurrentView("top")
	// maxX, _ := g.Size()
	// go GetPacket(v, maxX)
	return nil
}

func dquit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
