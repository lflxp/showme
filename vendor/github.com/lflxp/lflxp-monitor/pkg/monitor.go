package pkg

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
)

func Run(cmd string) {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	// 获取退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	ok := true

	interval := 20
	num := 0

	// 主机信息
	for _, x := range CollectEasy() {
		fmt.Println(x)
	}
	fmt.Println()

	// print net info
	// xo := MonitorNet{}
	// xo.Get()

	// err := GetHostInfo()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	FilterTitle(cmd, num, interval)

	for {
		num++
		select {
		case s := <-c:
			fmt.Printf("\n\033[1;4;31m%s:罒灬罒:小伙子走了哟！\033[0m\n", s)
			ok = false
			break
		case <-t.C:
			FilterTitle(cmd, num, interval)
			FilterValue(cmd)
		}
		// 终止循环
		if !ok {
			break
		}
	}
}

// 组装标题
func FilterTitle(in string, count, interval int) {
	title := GetTimeTitle()
	columns := GetTimeColumns()

	if strings.Contains(in, "-lazy") {
		title += GetLoadTitle()
		columns += GetLoadColumns()

		title += GetCpuTitle()
		columns += GetCpuColumns()

		title += GetSwapTitle()
		columns += GetSwapColumns()

		title += GetNetTitle(true)
		columns += GetNetColumns(true)
	} else {
		if strings.Contains(in, "-l") {
			title += GetLoadTitle()
			columns += GetLoadColumns()
		}

		if strings.Contains(in, "-c") {
			title += GetCpuTitle()
			columns += GetCpuColumns()
		}

		if strings.Contains(in, "-s") {
			title += GetSwapTitle()
			columns += GetSwapColumns()
		}

		if strings.Contains(in, "-n") {
			title += GetNetTitle(true)
			columns += GetNetColumns(true)
		}

		if strings.Contains(in, "-N") {
			title += GetNetTitle(false)
			columns += GetNetColumns(false)
		}
	}

	if strings.Contains(in, "-d") {
		title += GetDiskTitle()
		columns += GetDiskColumns()
	}

	if count%interval == 0 {
		fmt.Println(title)
		fmt.Println(columns)
	}
}

// 抽象命令
// if 顺序决定展示命令
func FilterValue(in string) {

	value, err := TimeNow()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if strings.Contains(in, "-lazy") {
		tmp_load, err := CpuLoad()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		value += tmp_load

		tmp_cpu, err := CpuPercent()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		value += tmp_cpu

		tmp_swap, err := SwapIO()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		value += tmp_swap

		tmp_net, err := NetInfo(true)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		value += tmp_net
	} else {
		if strings.Contains(in, "-l") {
			tmp_load, err := CpuLoad()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			value += tmp_load
		}

		if strings.Contains(in, "-c") {
			tmp_cpu, err := CpuPercent()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			value += tmp_cpu
		}

		if strings.Contains(in, "-s") {
			tmp_swap, err := SwapIO()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			value += tmp_swap
		}

		if strings.Contains(in, "-n") {
			tmp_net, err := NetInfo(true)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			value += tmp_net
		}

		if strings.Contains(in, "-N") {
			tmp_net, err := NetInfo(false)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			value += tmp_net
		}
	}

	if strings.Contains(in, "-d") {
		tmp_disk, err := DiskInfo()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		value += tmp_disk
	}

	fmt.Println(value)
}
