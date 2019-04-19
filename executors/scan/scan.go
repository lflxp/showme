package scan

import (
	"fmt"
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

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
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
	if v, err := g.SetView("help", 0, 0, 11, maxY-1); err != nil {
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

		if _, err = setCurrentViewOnTop(g, "help"); err != nil {
			return err
		}
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	if v, err := g.SetView("top", 11, 0, maxX, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "top"
		v.Wrap = true
		v.Autoscroll = false
		v.Editable = true

		data, err := utils.ParseIps(ips)
		if err != nil {
			fmt.Fprintln(v, err.Error())
		}
		// else {
		// 	for _, x := range data {
		// 		fmt.Fprintln(v, x)
		// 	}
		// }

		utils.Pings2(data, v)

		if _, err = setCurrentViewOnTop(g, "top"); err != nil {
			return err
		}
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	if v, err := g.SetView("bottom", 11, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "bottom"
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

		if _, err = setCurrentViewOnTop(g, "bottom"); err != nil {
			return err
		}
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	return nil
}
