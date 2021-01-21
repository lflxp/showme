package pkg

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

var ips string
var selectId string
var port string

func Scan(in string) {
	tmp := strings.Split(in, " ")
	if len(tmp) <= 1 {
		fmt.Println("输入错误，未指定IP地址")
		return
	}
	ips = tmp[1]
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(dlayout)

	// auto refresh every second
	d := time.Duration(time.Second)
	t := time.NewTicker(d)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				g.Update(func(g *gocui.Gui) error { return nil })
				// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				// fmt.Fprintln(v, )
			}
		}
	}()

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

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}
