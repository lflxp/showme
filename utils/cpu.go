package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
)

var before *MonitorCpu

func init() {
	before, err = NewMonitorCpu()
}

type MonitorCpu struct {
	//cpu
	cpu_cores      float64
	cpu_stolen     float64
	cpu_usr        float64
	cpu_nice       float64
	cpu_sys        float64
	cpu_idl        float64
	cpu_iow        float64
	cpu_irq        float64
	cpu_softirq    float64
	cpu_steal      float64
	cpu_guest      float64
	cpu_guest_nice float64
}

// get total cpu total
func (this *MonitorCpu) Get() error {
	t, err := cpu.Times(false)
	if err != nil {
		return err
	}

	c, _ := cpu.Info()
	this.cpu_cores, _ = strconv.ParseFloat(fmt.Sprintf("%d", c[0].Cores), 64)
	this.cpu_stolen = t[0].Stolen
	this.cpu_usr = t[0].User
	this.cpu_nice = t[0].Nice
	this.cpu_sys = t[0].System
	this.cpu_idl = t[0].Idle
	this.cpu_iow = t[0].Iowait
	this.cpu_irq = t[0].Irq
	this.cpu_softirq = t[0].Softirq
	this.cpu_steal = t[0].Steal
	this.cpu_guest = t[0].Guest
	this.cpu_guest_nice = t[0].GuestNice
	return nil
}

func NewMonitorCpu() (*MonitorCpu, error) {
	data := &MonitorCpu{}
	err := data.Get()
	return data, err
}

func CpuPercent() (string, error) {
	var rs string
	after, err := NewMonitorCpu()
	if err != nil {
		return rs, err
	}

	cpu_total1 := before.cpu_usr + before.cpu_nice + before.cpu_sys + before.cpu_idl + before.cpu_iow + before.cpu_irq + before.cpu_softirq
	cpu_total2 := after.cpu_usr + after.cpu_nice + after.cpu_sys + after.cpu_idl + after.cpu_iow + after.cpu_irq + after.cpu_softirq

	usr := (after.cpu_usr - before.cpu_usr) * 100 / (cpu_total2 - cpu_total1)
	sys := (after.cpu_sys - before.cpu_sys) * 100 / (cpu_total2 - cpu_total1)
	idl := (after.cpu_idl - before.cpu_idl) * 100 / (cpu_total2 - cpu_total1)
	iow := (after.cpu_iow - before.cpu_iow) * 100 / (cpu_total2 - cpu_total1)
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
