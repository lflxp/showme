package utils

import "github.com/shirou/gopsutil/load"

type MonitorLoad struct {
	//loadavg
	load_1  float64
	load_5  float64
	load_15 float64
}

func (this *MonitorLoad) Get() error {
	l, err := load.Avg()
	if err != nil {
		return err
	}

	this.load_1 = l.Load1
	this.load_5 = l.Load5
	this.load_15 = l.Load15
	return nil
}
