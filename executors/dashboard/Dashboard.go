package dashboard

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/shirou/gopsutil/mem"
)

func Dashboard() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(dlayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func dlayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/4-7, maxY/2, maxX/4+100, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		m, _ := mem.VirtualMemory()
		fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	return nil
}

func dquit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
