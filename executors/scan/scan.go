package scan

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
	"github.com/lflxp/showme/utils/table"
)

var ips string
var selectId string

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
	if err := g.SetKeybinding("top", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("top", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("top", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if v, err := g.SetCurrentView("bottom"); err != nil {
		return err
	} else {
		v.Highlight = true
		v.Autoscroll = true

		fmt.Fprintln(v, selectId)
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

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		// fmt.Fprintln(v, strings.Trim(l, " "))
		fmt.Fprintln(v, l)
		selectId = strings.Trim(l, " ")
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

func ScanIp(w io.Writer, width int) {
	stop := make(chan string)
	defer close(stop)

	data, err := utils.ParseIps(ips)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	go utils.Pings2(data, stop)

	num := 0
	for {
		select {
		case s, ok := <-stop:
			num++
			if !ok {
				break
			}
			tableNow := table.NewTable(width)

			tableNow.AddCol("ID").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
			tableNow.AddCol("IPAddress").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
			tableNow.AddCol("RTT").SetColor("dgreen").SetTextAlign(table.TextCenter).SetBgColor("black")
			tableNow.CalColumnWidths()
			if num%50 == 0 || num == 1 {
				tableNow.FprintHeader(w)
			}

			tmp := strings.Split(s, ":")
			id := table.NewCol()
			id.Data = fmt.Sprintf("%d", num)
			id.TextAlign = table.TextCenter
			id.Color = "yellow"
			tableNow.AddRow(0, id)

			ip := table.NewCol()
			ip.Data = tmp[0]
			ip.TextAlign = table.TextCenter
			ip.Color = "green"
			tableNow.AddRow(1, ip)

			rt := table.NewCol()
			rt.Data = strings.Trim(tmp[1], "\n")
			rt.TextAlign = table.TextCenter
			rt.Color = "red"
			tableNow.AddRow(2, rt)

			// fmt.Fprintln(w, s)
			tableNow.Fprint(w)
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
				fmt.Println("scanipport error")
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
		v.Highlight = true
		v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		go ScanIp(v, maxX-18)
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
