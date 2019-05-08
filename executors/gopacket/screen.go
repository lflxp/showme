package gopacket

// import (
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"strings"

// 	"github.com/jroimartin/gocui"
// 	"github.com/lflxp/showme/utils"
// 	"github.com/lflxp/showme/utils/table"
// )

// var device string

// // 去除字符串重复
// func quc(in string) string {
// 	tmp := strings.Split(in, " ")
// 	rs := []string{}
// 	for _, x := range tmp {
// 		if x != " " {
// 			rs = append(rs, x)
// 		}
// 	}
// 	return strings.Join(rs, " ")
// }

// func nextView(g *gocui.Gui, v *gocui.View) error {
// 	if v == nil || v.Name() == "side" {
// 		rcur, err := g.SetCurrentView("right")
// 		rcur.Clear()

// 		tmp := v.BufferLines()
// 		fmt.Fprintln(rcur, "Count: "+utils.Colorize(fmt.Sprintf("%d", len(tmp)), "red", "", true, true))
// 		if len(tmp) > 50 {
// 			for i := len(tmp) - 50; i < len(tmp); i++ {
// 				fmt.Fprintln(rcur, quc(tmp[i]))
// 			}
// 		} else {
// 			for _, x := range tmp {
// 				fmt.Fprintln(rcur, quc(x))
// 			}
// 		}

// 		return err
// 	}
// 	if v == nil || v.Name() == "right" {
// 		_, err := g.SetCurrentView("main")
// 		return err
// 	}
// 	_, err := g.SetCurrentView("side")
// 	// maxX, _ := g.Size()
// 	// go GetPacket(v, maxX)
// 	return err
// }

// func clearView(g *gocui.Gui, v *gocui.View) error {
// 	fmt.Println("ccccc view")
// 	v.Clear()
// 	return nil
// }

// func cursorDown(g *gocui.Gui, v *gocui.View) error {
// 	if v != nil {
// 		cx, cy := v.Cursor()
// 		if err := v.SetCursor(cx, cy+1); err != nil {
// 			ox, oy := v.Origin()
// 			if err := v.SetOrigin(ox, oy+1); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func cursorUp(g *gocui.Gui, v *gocui.View) error {
// 	if v != nil {
// 		ox, oy := v.Origin()
// 		cx, cy := v.Cursor()
// 		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
// 			if err := v.SetOrigin(ox, oy-1); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func getLine(g *gocui.Gui, v *gocui.View) error {
// 	var l string
// 	var err error

// 	_, cy := v.Cursor()
// 	if l, err = v.Line(cy); err != nil {
// 		l = ""
// 	}

// 	maxX, maxY := g.Size()
// 	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		fmt.Fprintln(v, l)
// 		if _, err := g.SetCurrentView("msg"); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func delMsg(g *gocui.Gui, v *gocui.View) error {
// 	if err := g.DeleteView("msg"); err != nil {
// 		return err
// 	}
// 	if _, err := g.SetCurrentView("side"); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func quit(g *gocui.Gui, v *gocui.View) error {
// 	return gocui.ErrQuit
// }

// func keybindings(g *gocui.Gui) error {
// 	// 清空side缓存
// 	if err := g.SetKeybinding("side", gocui.KeyCtrl2, gocui.ModNone, clearView); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("side", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("main", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("right", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
// 		return err
// 	}

// 	if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
// 		return err
// 	}
// 	if err := g.SetKeybinding("main", gocui.KeyCtrlW, gocui.ModNone, saveVisualMain); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func saveMain(g *gocui.Gui, v *gocui.View) error {
// 	f, err := ioutil.TempFile("", "gocui_demo_")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	p := make([]byte, 5)
// 	v.Rewind()
// 	for {
// 		n, err := v.Read(p)
// 		if n > 0 {
// 			if _, err := f.Write(p[:n]); err != nil {
// 				return err
// 			}
// 		}
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func saveVisualMain(g *gocui.Gui, v *gocui.View) error {
// 	f, err := ioutil.TempFile("", "gocui_demo_")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	vb := v.ViewBuffer()
// 	if _, err := io.Copy(f, strings.NewReader(vb)); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func layout(g *gocui.Gui) error {
// 	maxX, maxY := g.Size()
// 	if v, err := g.SetView("main", -1, -1, maxX/3, maxY/3); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		// b, err := ioutil.ReadFile("README.md")
// 		// if err != nil {
// 		// 	panic(err)
// 		// }
// 		// fmt.Fprintf(v, "%s", b)
// 		fmt.Fprintln(v, utils.Colorize("Shortcut keys", "dgreen", "", true, true))
// 		fmt.Fprintln(v, utils.Colorize("Tab      ->          wipe cache & change next view", "green", "", false, true))
// 		fmt.Fprintln(v, utils.Colorize("ArrowDown & ArrowUp On side view  ->  move cursor", "green", "", false, true))
// 		fmt.Fprintln(v, utils.Colorize("CtrlC  -> QUIT", "green", "", false, true))
// 		fmt.Fprintln(v, utils.Colorize("Enter  -> get current line info into input", "green", "", false, true))
// 		fmt.Fprintln(v, utils.Colorize("CtrlS  -> Save Side View", "green", "", false, true))
// 		fmt.Fprintln(v, utils.Colorize("CtrlW  -> Save Visual Side View", "green", "", false, true))
// 		// v.Editable = true
// 		v.Frame = true
// 		v.Wrap = true
// 		v.SelFgColor = gocui.ColorGreen
// 	}
// 	if v, err := g.SetView("right", maxX/3, -1, maxX, maxY/3); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}

// 		v.Editable = true
// 		v.Frame = true
// 		v.Wrap = true
// 		v.SelFgColor = gocui.ColorGreen
// 		fmt.Fprintln(v, "Waiting for Tabs it!")
// 		list, err := utils.ParseIps("10.1.1.1-200")
// 		if err != nil {
// 			fmt.Fprintln(v, err.Error())
// 		} else {
// 			for _, x := range list {
// 				fmt.Fprintln(v, x)
// 			}
// 		}
// 	}
// 	if v, err := g.SetView("side", -1, maxY/3, maxX, maxY); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		v.Highlight = true
// 		v.Autoscroll = true
// 		v.SelBgColor = gocui.ColorGreen
// 		v.SelFgColor = gocui.ColorBlack
// 		// fmt.Fprintln(v, "Item 1")
// 		// fmt.Fprintln(v, "Item 2")
// 		// fmt.Fprintln(v, "Item 3")
// 		// fmt.Fprint(v, "\rWill be")
// 		// fmt.Fprint(v, "deleted\rItem 4\nItem 5")
// 		// go table.TableTest(v, maxX)
// 		// go GetPacket(v, maxX)
// 		go GetPacketOrigin(v, maxX)
// 		if _, err := g.SetCurrentView("side"); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func GetPacketOrigin(w io.Writer, width int) {
// 	bigChan := make(chan interface{})
// 	defer close(bigChan)

// 	go utils.WatchDog(bigChan, device)

// 	num := 0
// 	for {
// 		select {
// 		case s, ok := <-bigChan:
// 			if !ok {
// 				return
// 			}
// 			if num > 10000 {
// 				return
// 			}
// 			num++

// 			tableNow := table.NewTable(width)
// 			tableNow.ShowHeader = true

// 			tableNow.AddCol("ID").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 			tableNow.AddCol("SrcIp").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 			tableNow.AddCol("SrcMac").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 			tableNow.AddCol("DstIp").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 			tableNow.AddCol("DstMac").SetColor("blue").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 			tableNow.AddCol("SrcPort").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 			tableNow.AddCol("DstPort").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 			tableNow.AddCol("Protocol").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 			if num%50 == 0 || num == 1 {
// 				tableNow.FprintHeader(w)
// 			}

// 			i := s.(*utils.Data)
// 			id := table.NewCol()
// 			id.Data = fmt.Sprintf("%d", num)
// 			id.TextAlign = table.TextCenter
// 			id.Color = "red"
// 			tableNow.AddRow(0, id)

// 			si := table.NewCol()
// 			si.Data = i.SrcIp
// 			si.TextAlign = table.TextCenter
// 			si.Color = "dgreen"
// 			tableNow.AddRow(1, si)

// 			sm := table.NewCol()
// 			sm.Data = i.SrcMac
// 			sm.TextAlign = table.TextCenter
// 			sm.Color = "blue"
// 			tableNow.AddRow(2, sm)

// 			di := table.NewCol()
// 			di.Data = i.DstIp
// 			di.TextAlign = table.TextCenter
// 			di.Color = "purple"
// 			tableNow.AddRow(3, di)

// 			dm := table.NewCol()
// 			dm.Data = i.DstMac
// 			dm.TextAlign = table.TextCenter
// 			dm.Color = "blue"
// 			tableNow.AddRow(4, dm)

// 			sp := table.NewCol()
// 			sp.Data = i.SrcPort
// 			sp.TextAlign = table.TextCenter
// 			sp.Color = "dgreen"
// 			tableNow.AddRow(5, sp)

// 			dp := table.NewCol()
// 			dp.Data = i.DstPort
// 			dp.TextAlign = table.TextCenter
// 			dp.Color = "purple"
// 			tableNow.AddRow(6, dp)

// 			po := table.NewCol()
// 			po.Data = i.Protocol
// 			po.TextAlign = table.TextCenter
// 			po.Color = "yellow"
// 			tableNow.AddRow(7, po)

// 			tableNow.CalColumnWidths()
// 			tableNow.Fprint(w)
// 			// tableNow.FprintOrderDesc(w)
// 		}
// 	}
// }

// func GetPacket(w io.Writer, width int) {
// 	// func GetPacket(w *gocui.View, width int) {
// 	// tableNow := table.NewTable(width)
// 	// tableNow.ShowHeader = true

// 	// tableNow.AddCol("ID").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	// tableNow.AddCol("SrcIp").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	// tableNow.AddCol("SrcMac").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	// tableNow.AddCol("DstIp").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	// tableNow.AddCol("DstMac").SetColor("blue").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	// tableNow.AddCol("SrcPort").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	// tableNow.AddCol("DstPort").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	// tableNow.AddCol("Protocol").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")

// 	bigChan := make(chan interface{})
// 	defer close(bigChan)
// 	// fmt.Println(table)
// 	// data := utils.WatchDogData(100, device)

// 	// for num, i := range data {
// 	// 	id := table.NewCol()
// 	// 	id.Data = fmt.Sprintf("%d", num)
// 	// 	id.TextAlign = table.TextCenter
// 	// 	id.Color = "red"
// 	// 	tableNow.AddRow(0, id)

// 	// 	si := table.NewCol()
// 	// 	si.Data = i.SrcIp
// 	// 	si.TextAlign = table.TextCenter
// 	// 	si.Color = "blue"
// 	// 	tableNow.AddRow(1, si)

// 	// 	sm := table.NewCol()
// 	// 	sm.Data = i.SrcMac
// 	// 	sm.TextAlign = table.TextCenter
// 	// 	sm.Color = "blue"
// 	// 	tableNow.AddRow(2, sm)

// 	// 	di := table.NewCol()
// 	// 	di.Data = i.DstIp
// 	// 	di.TextAlign = table.TextCenter
// 	// 	di.Color = "blue"
// 	// 	tableNow.AddRow(3, di)

// 	// 	dm := table.NewCol()
// 	// 	dm.Data = i.DstMac
// 	// 	dm.TextAlign = table.TextCenter
// 	// 	dm.Color = "blue"
// 	// 	tableNow.AddRow(4, dm)

// 	// 	sp := table.NewCol()
// 	// 	sp.Data = i.SrcPort
// 	// 	sp.TextAlign = table.TextCenter
// 	// 	sp.Color = "blue"
// 	// 	tableNow.AddRow(5, sp)

// 	// 	dp := table.NewCol()
// 	// 	dp.Data = i.DstPort
// 	// 	dp.TextAlign = table.TextCenter
// 	// 	dp.Color = "blue"
// 	// 	tableNow.AddRow(6, dp)

// 	// 	po := table.NewCol()
// 	// 	po.Data = i.Protocol
// 	// 	po.TextAlign = table.TextCenter
// 	// 	po.Color = "blue"
// 	// 	tableNow.AddRow(7, po)

// 	// 	// tt := table.NewCol()
// 	// 	// tt.Data = i.Type
// 	// 	// tt.TextAlign = table.TextCenter
// 	// 	// tt.Color = "blue"
// 	// 	// tableNow.AddRow(6, tt)
// 	// }

// 	// tableNow.CalColumnWidths()
// 	// tableNow.Fprint(w)

// 	go utils.WatchDog(bigChan, device)

// 	tableNow := table.NewTable(width)
// 	tableNow.ShowHeader = true

// 	tableNow.AddCol("ID").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	tableNow.AddCol("SrcIp").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	tableNow.AddCol("SrcMac").SetColor("red").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	tableNow.AddCol("DstIp").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	tableNow.AddCol("DstMac").SetColor("blue").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	tableNow.AddCol("SrcPort").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	tableNow.AddCol("DstPort").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	tableNow.AddCol("Protocol").SetColor("green").SetTextAlign(table.TextCenter).SetBgColor("dgreen")
// 	tableNow.CalColumnWidths()

// 	num := 0
// 	for {
// 		select {
// 		case s := <-bigChan:
// 			num++
// 			// 清空缓存 保持在2000左右
// 			// if num > 2000 {
// 			// 	for n, ts := range tableNow.Rows {
// 			// 		tableNow.Rows[n] = []*table.Col{ts[0]}
// 			// 	}
// 			// }

// 			i := s.(*utils.Data)
// 			id := table.NewCol()
// 			id.Data = fmt.Sprintf("%d", num)
// 			id.TextAlign = table.TextCenter
// 			id.Color = "red"
// 			tableNow.AddRowByIndex(0, id)

// 			si := table.NewCol()
// 			si.Data = i.SrcIp
// 			si.TextAlign = table.TextCenter
// 			si.Color = "blue"
// 			tableNow.AddRowByIndex(1, si)

// 			sm := table.NewCol()
// 			sm.Data = i.SrcMac
// 			sm.TextAlign = table.TextCenter
// 			sm.Color = "blue"
// 			tableNow.AddRowByIndex(2, sm)

// 			di := table.NewCol()
// 			di.Data = i.DstIp
// 			di.TextAlign = table.TextCenter
// 			di.Color = "blue"
// 			tableNow.AddRowByIndex(3, di)

// 			dm := table.NewCol()
// 			dm.Data = i.DstMac
// 			dm.TextAlign = table.TextCenter
// 			dm.Color = "blue"
// 			tableNow.AddRowByIndex(4, dm)

// 			sp := table.NewCol()
// 			sp.Data = i.SrcPort
// 			sp.TextAlign = table.TextCenter
// 			sp.Color = "blue"
// 			tableNow.AddRowByIndex(5, sp)

// 			dp := table.NewCol()
// 			dp.Data = i.DstPort
// 			dp.TextAlign = table.TextCenter
// 			dp.Color = "blue"
// 			tableNow.AddRowByIndex(6, dp)

// 			po := table.NewCol()
// 			po.Data = i.Protocol
// 			po.TextAlign = table.TextCenter
// 			po.Color = "blue"
// 			tableNow.AddRowByIndex(7, po)

// 			tableNow.Fprint(w)
// 			// tableNow.FprintOrderDesc(w)
// 		}
// 	}
// }

// func Screen(in string) {
// 	device = in
// 	g, err := gocui.NewGui(gocui.OutputNormal)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	defer g.Close()

// 	g.Cursor = true
// 	g.SelFgColor = gocui.ColorGreen

// 	g.SetManagerFunc(layout)

// 	if err := keybindings(g); err != nil {
// 		log.Panicln(err)
// 	}

// 	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
// 		log.Panicln(err)
// 	}
// }
