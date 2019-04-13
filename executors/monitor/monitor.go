package monitor

import (
	"fmt"
	"os"
	"os/signal"
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

	fmt.Println(utils.GetCpuTitle())
	fmt.Println(utils.GetCpuColumns())

	for {
		num++
		select {
		case s := <-c:
			fmt.Printf("\n\033[1;4;31m%s:罒灬罒:小伙子走了哟！\033[0m\n", s)
			ok = false
			break
		case <-t.C:
			if num%interval == 0 {
				fmt.Println(utils.GetCpuTitle())
				fmt.Println(utils.GetCpuColumns())
			}
			// fmt.Println(utils.Colorize(utils.GetNowTime(), "red", "black", true, true))
			s, err := utils.CpuPercent()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(s)
			}
		}
		// 终止循环
		if !ok {
			break
		}
	}
}
