package kubectl

import (
	"github.com/jroimartin/gocui"
)

func KeyDelete(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, delOtherView); err != nil {
		return err
	}
	return nil
}

func delOtherView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.View("Pod"); err == nil {
		if err = g.DeleteView("Pod"); err != nil {
			return err
		}
	}
	if _, err := g.View("gethelp"); err == nil {
		if err = g.DeleteView("gethelp"); err != nil {
			return err
		}
	}
	if _, err := g.View("msg"); err == nil {
		if err = g.DeleteView("msg"); err != nil {
			return err
		}
	}
	if _, err := g.View("Deployment"); err == nil {
		if err = g.DeleteView("Deployment"); err != nil {
			return err
		}
	}
	if _, err := g.View("delpod"); err == nil {
		if err = g.DeleteView("delpod"); err != nil {
			return err
		}
	}
	if _, err := g.View("Serviceed"); err == nil {
		if err = g.DeleteView("Serviceed"); err != nil {
			return err
		}
	}
	if _, err := setCurrentViewOnTop(g, "bottom"); err != nil {
		return err
	}
	return nil
}

func delOtherViewNoBack(g *gocui.Gui) error {
	if _, err := g.View("Pod"); err == nil {
		if err = g.DeleteView("Pod"); err != nil {
			return err
		}
	}
	if _, err := g.View("gethelp"); err == nil {
		if err = g.DeleteView("gethelp"); err != nil {
			return err
		}
	}
	if _, err := g.View("msg"); err == nil {
		if err = g.DeleteView("msg"); err != nil {
			return err
		}
	}
	if _, err := g.View("Deployment"); err == nil {
		if err = g.DeleteView("Deployment"); err != nil {
			return err
		}
	}
	if _, err := g.View("delpod"); err == nil {
		if err = g.DeleteView("delpod"); err != nil {
			return err
		}
	}
	if _, err := g.View("Serviceed"); err == nil {
		if err = g.DeleteView("Serviceed"); err != nil {
			return err
		}
	}
	return nil
}
