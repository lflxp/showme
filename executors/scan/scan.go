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
var port string

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
	// if err := g.SetKeybinding("help", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding("top", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("scanport", gocui.KeyF5, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("scanport", gocui.KeyTab, gocui.ModNone, toTop); err != nil {
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
	return nil
}

func inputIp(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("inputip", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "输入IP范围(eg: 10-192.168.1-50.256)"
		v.Highlight = true
		v.Editable = true
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))
		if _, err := g.SetCurrentView("inputip"); err != nil {
			return err
		}
	}

	return nil
}

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
		go ScanIp(vivia, maxX-4)
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
	maxX, maxY := g.Size()
	if v, err := g.SetView("scanport", maxX/4, maxY/4, maxX*3/4, maxY*3/4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "端口扫描"
		v.Highlight = true
		v.Editable = true
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))

		if v, err := g.SetCurrentView("scanport"); err != nil {
			return err
		} else {
			v.Highlight = true
			// v.Autoscroll = true
			v.Clear()
			go ScanIpPorts(v, maxX/2)
		}
	}

	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("port", maxX/2-30, maxY/2+3, maxX/2+30, maxY/2+5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Editable = true

		if _, err := g.SetCurrentView("port"); err != nil {
			return err
		}
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

func toTop(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("scanport"); err != nil {
		return err
	}

	if v == nil || v.Name() == "scanport" {
		_, err := g.SetCurrentView("top")
		return err
	}

	// if _, err = setCurrentViewOnTop(g, "top"); err != nil {
	// 	return err
	// }
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "top" {
		_, err := g.SetCurrentView("top")

		return err
	} else if v == nil || v.Name() == "scanport" {
		_, err := g.SetCurrentView("scanport")

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

	ttt := strings.Split(l, "|")
	if len(ttt) > 1 {
		maxX, maxY := g.Size()
		if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Highlight = true
			v.Title = "Your Selectd"
			// v.Editable = true
			// fmt.Fprintln(v, strings.Trim(l, " "))
			// fmt.Fprintln(v, l)
			// selectId = strings.Trim(l, " ")

			selectId = ttt[1]
			fmt.Fprintln(v, selectId)
			if _, err := g.SetCurrentView("msg"); err != nil {
				return err
			}

		}
	}

	return nil
}

func getPort(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("port", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Input Port Range(eg: 80,3306,25-100)"
		v.Highlight = true
		v.Editable = true
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))
		if _, err := g.SetCurrentView("port"); err != nil {
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
			ip.Data = fmt.Sprintf("|%s|", tmp[0])
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

func ScanIpPorts(w io.Writer, width int) {
	stop := make(chan string)
	defer close(stop)

	go utils.ScanPort2H(selectId, port, stop)

	num := 0
	for {
		select {
		case s, ok := <-stop:
			num++
			if !ok {
				fmt.Println("scanipport error")
				break
			}
			if !strings.Contains(s, "range") {
				tablePort := table.NewTable(width)
				tablePort.AddCol("ID").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
				tablePort.AddCol("IPAddress").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
				tablePort.AddCol("Port").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
				tablePort.AddCol("Status").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("black")
				tablePort.CalColumnWidths()
				if num%50 == 0 || num == 1 {
					tablePort.FprintHeader(w)
				}

				tmp := strings.Split(s, ":")
				id := table.NewCol()
				id.Data = fmt.Sprintf("%d", num)
				id.TextAlign = table.TextCenter
				id.Color = "yellow"
				tablePort.AddRow(0, id)

				ip := table.NewCol()
				ip.Data = tmp[0]
				ip.TextAlign = table.TextCenter
				ip.Color = "green"
				tablePort.AddRow(1, ip)

				rt := table.NewCol()
				rt.Data = tmp[1]
				rt.TextAlign = table.TextCenter
				rt.Color = "red"
				tablePort.AddRow(2, rt)

				st := table.NewCol()
				st.Data = "Active"
				st.TextAlign = table.TextCenter
				st.Color = "yellow"
				tablePort.AddRow(3, st)

				tablePort.Fprint(w)
			}
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
	if v, err := g.SetView("help", 0, maxY-7, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keybindings"
		fmt.Fprintln(v, "Tab: Next View/Refresh IP or Port")
		fmt.Fprintln(v, "Enter: Select IP/Commit Input")
		fmt.Fprintln(v, "F5: Input New IP range/Refresh IP or Port")
		fmt.Fprintln(v, "↑ ↓: Move View")
		fmt.Fprintln(v, "^c: Exit")
	}
	if v, err := g.SetView("top", 0, 0, maxX-1, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "IP扫描列表"
		v.Wrap = true
		v.Highlight = true
		// v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		go ScanIp(v, maxX-4)

		if _, err = setCurrentViewOnTop(g, "top"); err != nil {
			return err
		}
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	// if v, err := g.SetView("bottom", 0, maxY/2, maxX-1, maxY-1); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	v.Title = "PortScan Result"
	// 	v.Wrap = true
	// 	// v.Autoscroll = false
	// 	v.Editable = true

	// 	// go ScanIpPorts(v)

	// 	// if _, err = setCurrentViewOnTop(g, "bottom"); err != nil {
	// 	// 	return err
	// 	// }
	// 	// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
	// 	// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	// }
	return nil
}
