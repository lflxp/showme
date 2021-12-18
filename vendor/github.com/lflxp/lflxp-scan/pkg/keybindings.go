package pkg

import "github.com/jroimartin/gocui"

func keybindings(g *gocui.Gui) error {
	// 清空side缓存
	// if err := g.SetKeybinding("help", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding("top", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("scanport", gocui.KeyF5, gocui.ModNone, inputIpPorts); err != nil {
		return err
	}
	if err := g.SetKeybinding("inputIpPorts", gocui.KeyEnter, gocui.ModNone, delinputIpPorts); err != nil {
		return err
	}
	if err := g.SetKeybinding("scanport", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("top", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("top", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("scanport", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("scanport", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("top", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, gethelp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, getPort); err != nil {
		return err
	}
	if err := g.SetKeybinding("port", gocui.KeyEnter, gocui.ModNone, delPort); err != nil {
		return err
	}
	if err := g.SetKeybinding("top", gocui.KeyF5, gocui.ModNone, inputIp); err != nil {
		return err
	}
	if err := g.SetKeybinding("inputip", gocui.KeyEnter, gocui.ModNone, delinputIp); err != nil {
		return err
	}
	if err := g.SetKeybinding("gethelp", gocui.KeyF1, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("top", gocui.KeySpace, gocui.ModNone, searchIp); err != nil {
		return err
	}
	if err := g.SetKeybinding("searchIp", gocui.KeyEnter, gocui.ModNone, delsearchIp); err != nil {
		return err
	}
	if err := g.SetKeybinding("scanport", gocui.KeySpace, gocui.ModNone, searchPorts); err != nil {
		return err
	}
	if err := g.SetKeybinding("searchPorts", gocui.KeyEnter, gocui.ModNone, delsearchPorts); err != nil {
		return err
	}
	return nil
}
