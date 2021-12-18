package pkg

import (
	"fmt"
	"io"
	"strings"

	"github.com/jroimartin/gocui"
)

func ScanIp(g *gocui.Gui, w *gocui.View, width int) error {
	stop := make(chan string)
	defer close(stop)

	data, err := ParseIps(ips)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return err
	}
	// empty process tab
	if v, err := g.View("help"); err != nil {
		return err
	} else {
		v.Clear()
		v.Autoscroll = true
		v.Wrap = true
		v.Highlight = true
		// v.Editable = true
		go Pings2(data, stop, v)
	}

	num := 0
	for {
		select {
		case s, ok := <-stop:
			num++
			if !ok {
				break
			}
			tableNow := NewTable(width)

			tableNow.AddCol("ID").SetColor("dgreen").SetTextAlign(TextCenter).SetBgColor("black")
			tableNow.AddCol("IPAddress").SetColor("dgreen").SetTextAlign(TextCenter).SetBgColor("black")
			tableNow.AddCol("RTT").SetColor("dgreen").SetTextAlign(TextCenter).SetBgColor("black")
			tableNow.CalColumnWidths()
			if num%50 == 0 || num == 1 {
				tableNow.FprintHeader(w)
			}

			tmp := strings.Split(s, ":")
			id := NewCol()
			id.Data = fmt.Sprintf("%d", num)
			id.TextAlign = TextCenter
			id.Color = "yellow"
			tableNow.AddRow(0, id)

			ip := NewCol()
			ip.Data = fmt.Sprintf("|%s|", tmp[0])
			ip.TextAlign = TextCenter
			ip.Color = "green"
			tableNow.AddRow(1, ip)

			rt := NewCol()
			rt.Data = strings.Trim(tmp[1], "\n")
			rt.TextAlign = TextCenter
			rt.Color = "red"
			tableNow.AddRow(2, rt)

			// fmt.Fprintln(w, s)
			tableNow.Fprint(w)
		}
	}
}

func ScanIpPorts(g *gocui.Gui, w io.Writer, width int) error {
	stop := make(chan string)
	defer close(stop)

	if v, err := g.View("help"); err != nil {
		return err
	} else {
		v.Clear()
		v.Autoscroll = true
		v.Highlight = true
		v.Wrap = true
		go ScanPort2H(selectId, port, stop, v)
	}

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
				tablePort := NewTable(width)
				tablePort.AddCol("ID").SetColor("red").SetTextAlign(TextCenter).SetBgColor("black")
				tablePort.AddCol("IPAddress").SetColor("red").SetTextAlign(TextCenter).SetBgColor("black")
				tablePort.AddCol("Port").SetColor("red").SetTextAlign(TextCenter).SetBgColor("black")
				tablePort.AddCol("Status").SetColor("red").SetTextAlign(TextCenter).SetBgColor("black")
				tablePort.CalColumnWidths()
				if num%50 == 0 || num == 1 {
					tablePort.FprintHeader(w)
				}

				tmp := strings.Split(s, ":")
				id := NewCol()
				id.Data = fmt.Sprintf("%d", num)
				id.TextAlign = TextCenter
				id.Color = "yellow"
				tablePort.AddRow(0, id)

				ip := NewCol()
				ip.Data = tmp[0]
				ip.TextAlign = TextCenter
				ip.Color = "green"
				tablePort.AddRow(1, ip)

				rt := NewCol()
				rt.Data = tmp[1]
				rt.TextAlign = TextCenter
				rt.Color = "red"
				tablePort.AddRow(2, rt)

				st := NewCol()
				st.Data = "Active"
				st.TextAlign = TextCenter
				st.Color = "yellow"
				tablePort.AddRow(3, st)

				tablePort.Fprint(w)
			}
		}
	}
	return nil
}
