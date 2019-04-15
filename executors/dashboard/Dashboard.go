package dashboard

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	net1 "github.com/shirou/gopsutil/net"
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
	out.Clear()
	data := collect()
	for _, x := range data {
		fmt.Fprintln(v, x)
	}

	fmt.Fprintln(out, utils.Colorize(time.Now().Format("2006-01-02 15:04:05"), "white", "green", false, true))

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
	load, _ := load.Avg()
	// cc, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	n, _ := host.Info()
	nv, _ := net1.IOCounters(true)
	boottime, _ := host.BootTime()
	btime := time.Unix(int64(boottime), 0).Format("2006-01-02 15:04:05")
	rs = append(rs, fmt.Sprintf("%s: %v MB  Free: %v MB Used:%v Usage:%f%%\n", utils.Colorize("        Mem       ", "white", "red", true, true), v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent))
	// fmt.Printf("        Mem       : %v MB  Free: %v MB Used:%v Usage:%f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
	if len(c) > 1 {
		for _, sub_cpu := range c {
			modelname := sub_cpu.ModelName
			cores := sub_cpu.Cores
			// fmt.Printf("        CPU       : %v   %v cores \n", modelname, cores)
			rs = append(rs, fmt.Sprintf("%s: %v   %v cores \n", utils.Colorize("        CPU       ", "white", "red", true, true), modelname, cores))
		}
	} else {
		sub_cpu := c[0]
		modelname := sub_cpu.ModelName
		cores := sub_cpu.Cores
		rs = append(rs, fmt.Sprintf("%s: %v   %v cores \n", utils.Colorize("        CPU       ", "white", "red", true, true), modelname, cores))
	}
	rs = append(rs, fmt.Sprintf("%s: %.2f %.2f %.2f \n", utils.Colorize("        LOAD      ", "white", "red", true, true), load.Load1, load.Load5, load.Load15))
	rs = append(rs, fmt.Sprintf("%s: %v bytes / %v bytes\n", utils.Colorize("        Network   ", "white", "red", true, true), nv[0].BytesRecv, nv[0].BytesSent))
	rs = append(rs, fmt.Sprintf("%s:%v\n", utils.Colorize("        SystemBoot", "white", "red", true, true), btime))
	// rs = append(rs, fmt.Sprintf("        CPU Used    : used %f%% \n", cc[0]))
	rs = append(rs, fmt.Sprintf("%s: %v GB  Free: %v GB Usage:%f%%\n", utils.Colorize("        HD        ", "white", "red", true, true), d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent))
	rs = append(rs, fmt.Sprintf("%s :%v %v(%v)   %v  \n", utils.Colorize("        OS        ", "white", "red", true, true), n.OS, n.Platform, n.PlatformFamily, n.PlatformVersion))
	rs = append(rs, fmt.Sprintf("%s: %v  \n", utils.Colorize("        Hostname  ", "white", "red", true, true), n.Hostname))
	rs = append(rs, fmt.Sprintf("%s: %s", utils.Colorize("        IpLists   ", "white", "red", true, true), strings.Join(utils.GetIps(), ",")))
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
