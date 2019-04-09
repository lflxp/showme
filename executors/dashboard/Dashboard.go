package dashboard

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type List struct {
	Key   string
	Value string
}

type Info struct {
	Title string
	Data  []List
}

func dquit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func nextViewDashboard(g *gocui.Gui, v *gocui.View) error {
	out, err := g.View("hello")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, time.Now().Format("2006-01-02 15:04:05"))

	if _, err = setCurrentViewOnTop(g, "hello"); err != nil {
		return err
	}

	g.Cursor = true
	return nil
}

func Dashboard() {
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

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextViewDashboard); err != nil {
		log.Println(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// func getInfo() []Info {
// 	var rs []Info

// }

func collect() []string {
	rs := []string{}
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	n, _ := host.Info()
	nv, _ := net.IOCounters(true)
	boottime, _ := host.BootTime()
	btime := time.Unix(int64(boottime), 0).Format("2006-01-02 15:04:05")
	rs = append(rs, fmt.Sprintf("        Mem       : %v MB  Free: %v MB Used:%v Usage:%f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent))
	// fmt.Printf("        Mem       : %v MB  Free: %v MB Used:%v Usage:%f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
	if len(c) > 1 {
		for _, sub_cpu := range c {
			modelname := sub_cpu.ModelName
			cores := sub_cpu.Cores
			// fmt.Printf("        CPU       : %v   %v cores \n", modelname, cores)
			rs = append(rs, fmt.Sprintf("        CPU       : %v   %v cores \n", modelname, cores))
		}
	} else {
		sub_cpu := c[0]
		modelname := sub_cpu.ModelName
		cores := sub_cpu.Cores
		rs = append(rs, fmt.Sprintf("        CPU       : %v   %v cores \n", modelname, cores))
	}
	rs = append(rs, fmt.Sprintf("        Network: %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent))
	rs = append(rs, fmt.Sprintf("        SystemBoot:%v\n", btime))
	rs = append(rs, fmt.Sprintf("        CPU Used    : used %f%% \n", cc[0]))
	rs = append(rs, fmt.Sprintf("        HD        : %v GB  Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent))
	rs = append(rs, fmt.Sprintf("        OS        : %v(%v)   %v  \n", n.Platform, n.PlatformFamily, n.PlatformVersion))
	rs = append(rs, fmt.Sprintf("        Hostname  : %v  \n", n.Hostname))
	return rs
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

// 1.一行四列
// 2.一共九宫格
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
	if v, err := g.SetView("hello", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dashboard"
		v.Wrap = true
		// v.Autoscroll = true
		v.Editable = true

		data := collect()
		for _, x := range data {
			fmt.Fprintln(v, x)
		}

		if _, err = setCurrentViewOnTop(g, "hello"); err != nil {
			return err
		}
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Fprintln(v, fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", m.Total, m.Free, m.UsedPercent))
	}
	return nil
}
