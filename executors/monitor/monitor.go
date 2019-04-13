package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/lflxp/showme/utils"
)

func Run(cmd string) {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	// 获取退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	ok := true

	interval := 10
	num := 0

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
	title := utils.GetTimeTitle()
	columns := utils.GetTimeColumns()

	if strings.Contains(in, "-l") {
		title += utils.GetLoadTitle()
		columns += utils.GetLoadColumns()
	}

	if strings.Contains(in, "-c") {
		title += utils.GetCpuTitle()
		columns += utils.GetCpuColumns()
	}

	if count%interval == 0 {
		fmt.Println(title)
		fmt.Println(columns)
	}
}

// 抽象命令
// if 顺序决定展示命令
func FilterValue(in string) {

	value, err := utils.TimeNow()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if strings.Contains(in, "-l") {
		tmp_load, err := utils.CpuLoad()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		value += tmp_load
	}

	if strings.Contains(in, "-c") {
		tmp_cpu, err := utils.CpuPercent()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		value += tmp_cpu
	}

	fmt.Println(value)
}
