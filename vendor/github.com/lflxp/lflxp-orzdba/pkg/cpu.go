package pkg

import (
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
)

var before []cpu.TimesStat

func init() {
	before, err = cpu.Times(false)
}

func CpuPercent() (string, error) {
	var rs string
	// after, err := NewMonitorCpu()
	after, err := cpu.Times(false)
	if err != nil {
		return rs, err
	}

	// cpu_total1 := before.cpu_usr + before.cpu_nice + before.cpu_sys + before.cpu_idl + before.cpu_iow + before.cpu_irq + before.cpu_softirq
	cpu_total1 := before[0].User + before[0].Nice + before[0].System + before[0].Idle + before[0].Iowait + before[0].Irq + before[0].Softirq
	// cpu_total2 := after.cpu_usr + after.cpu_nice + after.cpu_sys + after.cpu_idl + after.cpu_iow + after.cpu_irq + after.cpu_softirq
	cpu_total2 := after[0].User + after[0].Nice + after[0].System + after[0].Idle + after[0].Iowait + after[0].Irq + after[0].Softirq

	usr := (after[0].User - before[0].User) * 100 / (cpu_total2 - cpu_total1)
	sys := (after[0].System - before[0].System) * 100 / (cpu_total2 - cpu_total1)
	idl := (after[0].Idle - before[0].Idle) * 100 / (cpu_total2 - cpu_total1)
	iow := (after[0].Iowait - before[0].Iowait) * 100 / (cpu_total2 - cpu_total1)
	// usr
	if usr > 10.0 {
		rs += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(int(usr))))+strconv.Itoa(int(usr))+" ", "red", "", false, true)
	} else {
		rs += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(int(usr))))+strconv.Itoa(int(usr))+" ", "green", "", false, false)
	}

	if sys > 10.0 {
		rs += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(int(sys))))+strconv.Itoa(int(sys))+" ", "red", "", false, true)
	} else {
		rs += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(int(sys))))+strconv.Itoa(int(sys))+" ", "white", "", false, false)
	}

	if idl < 30.0 {
		rs += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(int(idl))))+strconv.Itoa(int(idl))+" ", "red", "", false, true)
	} else {
		rs += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(int(idl))))+strconv.Itoa(int(idl))+" ", "white", "", false, false)
	}

	if iow > 10.0 {
		rs += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(int(iow))))+strconv.Itoa(int(iow))+" ", "red", "", false, true)
	} else {
		rs += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(int(iow))))+strconv.Itoa(int(iow))+" ", "green", "", false, false)
	}

	rs += Colorize("|", "dgreen", "", false, false)
	before = after
	return rs, nil
}
