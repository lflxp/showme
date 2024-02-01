package pkg

import (
	"strings"

	"github.com/shirou/gopsutil/v3/load"
)

type MonitorLoad struct {
	//loadavg
	load_1  float64
	load_5  float64
	load_15 float64
}

func (this *MonitorLoad) Get() (string, error) {
	l, err := load.Avg()
	if err != nil {
		return "", err
	}

	this.load_1 = l.Load1
	this.load_5 = l.Load5
	this.load_15 = l.Load15
	data_detail := ""

	if this.load_1 >= 10.0 {
		data_detail += Colorize(strings.Repeat(" ", 5-len(floatToString(this.load_1, 2)))+floatToString(this.load_1, 2), "red", "", false, true)
	} else {
		data_detail += Colorize(strings.Repeat(" ", 5-len(floatToString(this.load_1, 2)))+floatToString(this.load_1, 2), "green", "", false, false)
	}

	if this.load_1 >= 10.0 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(this.load_5, 2)))+floatToString(this.load_5, 2), "red", "", false, true)
	} else {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(this.load_5, 2)))+floatToString(this.load_5, 2), "green", "", false, false)
	}

	if this.load_1 >= 10.0 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(this.load_15, 2)))+floatToString(this.load_15, 2), "red", "", false, true) + Colorize("|", "dgreen", "", false, false)
	} else {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(this.load_15, 2)))+floatToString(this.load_15, 2), "green", "", false, false) + Colorize("|", "dgreen", "", false, false)
	}

	return data_detail, nil
}

func NewLoad() *MonitorLoad {
	return &MonitorLoad{}
}

func CpuLoad() (string, error) {
	data := NewLoad()
	rs, err := data.Get()
	return rs, err
}
