package pkg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

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
	if v, err := g.SetView("help", 0, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Processing"
	}
	if v, err := g.SetView("top", 0, 0, maxX/2-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "IP List Result(Help: F1)"
		v.Wrap = true
		v.Highlight = true
		// v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		go ScanIp(g, v, maxX/2-3)

		if _, err = setCurrentViewOnTop(g, "top"); err != nil {
			return err
		}
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	if v, err := g.SetView("scanport", maxX/2, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Port Scan Result"
		v.Highlight = true
		v.Editable = true
		v.Autoscroll = true
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))

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
		// v.Frame = false
		// v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorRed
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
			v.SelFgColor = gocui.ColorYellow
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

func gethelp(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("gethelp", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Keybindings(F1: quit)"
		fmt.Fprintln(v, Colorize("Tab: Next View", "yellow", "", false, true))
		fmt.Fprintln(v, Colorize("Enter: Select IP/Commit Input", "yellow", "", false, true))
		fmt.Fprintln(v, Colorize("F5: Input New Scan IP or Port range", "yellow", "", false, true))
		fmt.Fprintln(v, Colorize("↑ ↓: Move View", "yellow", "", false, true))
		fmt.Fprintln(v, Colorize("^c: Exit", "yellow", "", false, true))
		fmt.Fprintln(v, Colorize("F1: Help", "yellow", "", false, true))
		fmt.Fprintln(v, Colorize("Space: search result with ip view and port view", "yellow", "", false, true))

		if _, err := setCurrentViewOnTop(g, "gethelp"); err != nil {
			return err
		}
	}

	return nil
}

func inputIpPorts(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("inputIpPorts", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = fmt.Sprintf("Current IP %s Input Port Range:", selectId)
		v.Highlight = true
		v.Editable = true
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		fmt.Fprintf(v, port)
		if _, err := g.SetCurrentView("inputIpPorts"); err != nil {
			return err
		}
	}

	return nil
}

func delinputIpPorts(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("inputIpPorts"); err != nil {
		return err
	}

	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	port = l
	maxX, _ := g.Size()
	if vivia, err := g.SetCurrentView("scanport"); err != nil {
		return err
	} else {
		vivia.Title = fmt.Sprintf("Current IP:   %s     PORTS RANGE:  %s", selectId, port)
		vivia.Highlight = true
		vivia.Clear()
		go ScanIpPorts(g, vivia, maxX/2)
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
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		fmt.Fprintf(v, ips)
		if _, err := g.SetCurrentView("inputip"); err != nil {
			return err
		}
	}

	return nil
}

func searchIp(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("searchIp", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Search"
		v.Highlight = true
		v.Editable = true
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		if _, err := g.SetCurrentView("searchIp"); err != nil {
			return err
		}
	}

	return nil
}

func searchPorts(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("searchPorts", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Search"
		v.Highlight = true
		v.Editable = true
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		if _, err := g.SetCurrentView("searchPorts"); err != nil {
			return err
		}
	}

	return nil
}

func delsearchIp(g *gocui.Gui, v *gocui.View) error {
	// getline
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	// getResult
	if v, err := g.View("top"); err != nil {
		return err
	} else {
		// v.Autoscroll = true
		// v.Highlight = true
		// v.Wrap = true
		// deleteOrigin
		if err := g.DeleteView("searchIp"); err != nil {
			return err
		}
		// color for search word and replace before result
		tmprs := v.BufferLines()
		v.Clear()
		front := []string{}
		backend := []string{}
		title := ""
		countTitle := 0
		for _, x := range tmprs {
			x = strings.Replace(x, ">", " ", -1)
			match, _ := regexp.MatchString(l, x)
			if match {
				// fmt.Fprintln(v, Colorize(strings.Replace(x, " ", ">", 3), "yellow", "", false, true))
				front = append(front, Colorize(strings.Replace(x, " ", ">", 3), "yellow", "", false, true))
			} else {
				match1, _ := regexp.MatchString("ID", x)
				if match1 && countTitle == 0 {
					// fmt.Fprintln(v, Colorize(x, "dgreen", "black", false, false))
					title = Colorize(x, "dgreen", "black", false, false)
					countTitle++
				} else {
					// fmt.Fprintln(v, strings.Replace(x, ">", " ", -1))
					backend = append(backend, x)
				}
			}
		}
		// print
		fmt.Fprintln(v, title)
		for _, o := range front {
			fmt.Fprintln(v, o)
		}
		for _, t := range backend {
			fmt.Fprintln(v, t)
		}
		// setTop
		if _, err := g.SetCurrentView("top"); err != nil {
			return err
		}
	}

	return nil
}

func delsearchPorts(g *gocui.Gui, v *gocui.View) error {
	// getline
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	// getResult
	if v, err := g.View("scanport"); err != nil {
		return err
	} else {
		// v.Autoscroll = true
		// v.Highlight = true
		// v.Wrap = true
		// deleteOrigin
		if err := g.DeleteView("searchPorts"); err != nil {
			return err
		}
		// color for search word and replace before result
		tmprs := v.BufferLines()
		v.Clear()
		front := []string{}
		backend := []string{}
		title := ""
		countTitle := 0
		for _, x := range tmprs {
			x = strings.Replace(x, ">", " ", -1)
			match, _ := regexp.MatchString(l, x)
			if match {
				// fmt.Fprintln(v, Colorize(strings.Replace(x, " ", ">", 3), "yellow", "", false, true))
				front = append(front, Colorize(strings.Replace(x, " ", ">", 3), "yellow", "", false, true))
			} else {
				match1, _ := regexp.MatchString("ID", x)
				if match1 && countTitle == 0 {
					// fmt.Fprintln(v, Colorize(x, "dgreen", "black", false, false))
					title = Colorize(x, "dgreen", "black", false, false)
					countTitle++
				} else {
					// fmt.Fprintln(v, strings.Replace(x, ">", " ", -1))
					backend = append(backend, x)
				}
			}
		}
		// print
		fmt.Fprintln(v, title)
		for _, o := range front {
			fmt.Fprintln(v, o)
		}
		for _, t := range backend {
			fmt.Fprintln(v, t)
		}
		// setTop
		if _, err := g.SetCurrentView("scanport"); err != nil {
			return err
		}
	}

	return nil
}
