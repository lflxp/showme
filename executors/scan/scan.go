package scan

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
)

var ips string

func Scan(in string) {
	ips = strings.Split(in, " ")[1]
	fmt.Println("in", in)
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(dlayout)

	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
	// 	log.Panicln(err)
	// }

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func keybindings(g *gocui.Gui) error {
	// 清空side缓存
	if err := g.SetKeybinding("help", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("top", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("bottom", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		return err
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
	if v == nil || v.Name() == "bottom" {
		_, err := g.SetCurrentView("top")

		return err
	}
	if v == nil || v.Name() == "top" {
		_, err := g.SetCurrentView("help")
		return err
	}
	_, err := g.SetCurrentView("bottom")
	// maxX, _ := g.Size()
	// go GetPacket(v, maxX)
	return err
}

func dquit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func ScanIp(w io.Writer) {
	stop := make(chan string)
	defer close(stop)

	data, err := utils.ParseIps(ips)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	go utils.Pings2(data, stop)

	for {
		select {
		case s, ok := <-stop:
			if !ok {
				break
			}
			fmt.Fprintln(w, s)
		}
	}
}

func ScanIpPorts(w io.Writer) {
	stop := make(chan string)
	defer close(stop)

	data, err := utils.ParseIps(ips)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	go utils.Pings3(data, stop)

	for {
		select {
		case s, ok := <-stop:
			if !ok {
				break
			}
			fmt.Fprintln(w, s)
		}
	}
}

func dlayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// if v, err := g.SetView("hello", maxX/4-7, maxY/2, maxX/4+100, maxY/2+2); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	m, _ := mem.VirtualMemory()
	// 	fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	// }

	// log.Println(data)
	if v, err := g.SetView("help", 0, 0, 15, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "IP列表"
		v.Wrap = true
		v.Autoscroll = false
		v.Editable = true

		data, err := utils.ParseIps(ips)
		if err != nil {
			fmt.Fprintln(v, err.Error())
		} else {
			for _, x := range data {
				fmt.Fprintln(v, x)
			}
		}

		// if _, err = setCurrentViewOnTop(g, "help"); err != nil {
		// 	return err
		// }
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	if v, err := g.SetView("top", 15, 0, maxX-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "top"
		v.Wrap = true
		v.Autoscroll = false
		v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		go ScanIp(v)
		if _, err = setCurrentViewOnTop(g, "top"); err != nil {
			return err
		}

		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	if v, err := g.SetView("bottom", 15, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "bottom"
		v.Wrap = true
		v.Autoscroll = false
		v.Editable = true

		go ScanIpPorts(v)

		// if _, err = setCurrentViewOnTop(g, "bottom"); err != nil {
		// 	return err
		// }
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	return nil
}
