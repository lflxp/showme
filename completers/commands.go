package completers

import (
	"fmt"
	"log"

	"github.com/c-bata/go-prompt"
	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
)

var Commands = []prompt.Suggest{
	{Text: "dashboard", Description: "computer configeration"},
	{Text: "gocui", Description: "https://github.com/jroimartin/gocui"},
	{Text: "monitor", Description: "monitoring Linux/Unix or MacOs status runtime"},
	{Text: "scan", Description: "ip && port scaning online"},
	{Text: "mysql", Description: "monitor mysql info"},
	{Text: "kubectl", Description: "kubectl 可视化管理界面"},
	{Text: "help", Description: "List All Menu"},
}

func GetHelp() []string {
	rs := []string{}
	for _, x := range Commands {
		rs = append(rs, fmt.Sprintf("%s: %s", utils.Colorize(x.Text, "white", "green", false, true), x.Description))
	}
	return rs
}

func Help() {
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
	if v, err := g.SetView("help", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "功能菜单一览"
		v.Wrap = true
		v.Autoscroll = true
		v.Editable = false

		data := GetHelp()
		for _, x := range data {
			fmt.Fprintln(v, x)
		}

		if _, err = setCurrentViewOnTop(g, "help"); err != nil {
			return err
		}
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	return nil
}
