// +build linux

package utils

import (
	"encoding/json"

	"github.com/shirou/gopsutil/host"
)

type MonitorHost struct {
	Ip []string `json:ip`
	host.InfoStat
}

func NewHost() (*MonitorHost, error) {
	data := &MonitorHost{}
	err := data.Set()
	return data, err
}

func (this *MonitorHost) Set() error {
	err := this.SetHostname()
	if err != nil {
		return err
	}
	err = this.SetIps()
	if err != nil {
		return err
	}
	return nil
}

func (this *MonitorHost) SetHostname() error {
	n, err := host.Info()
	if err != nil {
		return err
	}
	this.BootTime = n.BootTime
	this.HostID = n.HostID
	this.Hostname = n.Hostname
	this.OS = n.OS
	return nil
}

func (this *MonitorHost) SetIps() error {
	this.Ip = GetIps()
	return nil
}

func (this *MonitorHost) String() string {
	bytes, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
