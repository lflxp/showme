package scan

import (
	"fmt"
	"io"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
	"github.com/lflxp/showme/utils/table"
)

func ScanIp(g *gocui.Gui, w *gocui.View, width int) error {
	stop := make(chan string)
	defer close(stop)

	data, err := utils.ParseIps(ips)
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
		go utils.Pings2(data, stop, v)
	}

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
		go utils.ScanPort2H(selectId, port, stop, v)
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
	return nil
}
