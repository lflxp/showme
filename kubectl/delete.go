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

// refresh dashboard
func delOtherView(g *gocui.Gui, v *gocui.View) error {
	// refresh dashboard view
	if len(origin.Cluster) > 0 {
		for _, x := range origin.Cluster {
			if _, err := g.View(x.Title); err == nil {
				if err = g.DeleteView(x.Title); err != nil {
					return err
				}
			}
		}
	}
	if _, err := g.View("bottom"); err == nil {
		if err = g.DeleteView("bottom"); err != nil {
			return err
		}
	}
	if _, err := g.View("pod"); err == nil {
		if err = g.DeleteView("pod"); err != nil {
			return err
		}
	}
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
	if _, err := g.View("Noded"); err == nil {
		if err = g.DeleteView("Noded"); err != nil {
			return err
		}
	}
	if err := dashboard(g); err != nil {
		return err
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
	if _, err := g.View("Noded"); err == nil {
		if err = g.DeleteView("Noded"); err != nil {
			return err
		}
	}
	return nil
}
