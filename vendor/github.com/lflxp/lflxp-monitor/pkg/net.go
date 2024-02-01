package pkg

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	net1 "github.com/shirou/gopsutil/v3/net"
)

var beforeNet *MonitorNet

func init() {
	beforeNet, err = NewNet()
}

func GetIps() []string {
	rs := []string{}
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		rs = append(rs, err.Error())
		return rs
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				rs = append(rs, ipnet.IP.String())
			}
		}
	}
	return rs
}

func GetCurrentInterfaceCommands() ([]prompt.Suggest, error) {
	var rs []prompt.Suggest
	data, err := net1.Interfaces()
	if err != nil {
		return rs, err
	}

	if len(data) > 0 {
		rs = []prompt.Suggest{}
		for _, x := range data {
			ips := []string{}
			for _, y := range x.Addrs {
				ips = append(ips, y.Addr)
			}
			rs = append(rs, prompt.Suggest{
				Text:        x.Name,
				Description: fmt.Sprintf("%s %s", strings.Join(ips, ","), x.HardwareAddr),
			})
		}
	}
	return rs, nil
}

type MonitorNet struct {
	BytesSent   uint64 `json:"bytesSent"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytesRecv"`   // number of bytes received
	PacketsSent uint64 `json:"packetsSent"` // number of packets sent
	PacketsRecv uint64 `json:"packetsRecv"` // number of packets received
	Errin       uint64 `json:"errin"`       // total number of errors while receiving
	Errout      uint64 `json:"errout"`      // total number of errors while sending
	Dropin      uint64 `json:"dropin"`      // total number of incoming packets which were dropped
	Dropout     uint64 `json:"dropout"`     // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
	Fifoin      uint64 `json:"fifoin"`      // total number of FIFO buffers errors while receiving
	Fifoout     uint64 `json:"fifoout"`     // total number of FIFO buffers errors while sendin
}

func (this *MonitorNet) Get() error {
	data, err := net1.IOCounters(false)
	if err != nil {
		return err
	}

	this.BytesSent = data[0].BytesSent
	this.BytesRecv = data[0].BytesRecv
	this.PacketsRecv = data[0].PacketsRecv
	this.PacketsSent = data[0].PacketsSent
	this.Errin = data[0].Errin
	this.Errout = data[0].Errout
	this.Dropin = data[0].Dropin
	this.Dropout = data[0].Dropout
	this.Fifoin = data[0].Fifoin
	this.Fifoout = data[0].Fifoout
	return nil
}

func (this *MonitorNet) PrintIps() error {
	data, err := net1.Interfaces()
	if err != nil {
		return err
	}

	for _, x := range data {
		if len(x.Addrs) > 0 {
			fmt.Println(fmt.Sprintf("Flags %s Addr %s MTU %d Name %s Addrs %s", strings.Join(x.Flags, ","), x.HardwareAddr, x.MTU, x.Name))
			for _, y := range x.Addrs {
				fmt.Println("addrs", y.Addr)
			}
		} else {
			fmt.Println(fmt.Sprintf("Flags %s Addr %s MTU %d Name %s", strings.Join(x.Flags, ","), x.HardwareAddr, x.MTU, x.Name))
		}
	}
	return nil
}

func NewNet() (*MonitorNet, error) {
	data := &MonitorNet{}
	err = data.Get()
	return data, err
}

func NetInfo(detail bool) (string, error) {
	var rs string
	after, err := NewNet()
	if err != nil {
		return rs, err
	}

	netIn := float64(after.BytesRecv-beforeNet.BytesRecv) / 0.99
	netOut := float64(after.BytesSent-beforeNet.BytesSent) / 0.99

	if netIn/1024/1024 >= 1.0 {
		// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(netIn/1024/1024, 1)))+floatToString(netIn/1024/1024, 1)+"m", "red", "", false, true)
		rs += Colorize(parseRepeatSpace(floatToString(netIn/1024/1024, 1), 6)+"m", "red", "", false, true)
	} else if netIn/1024 < 1.0 {
		// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(netIn))))+strconv.Itoa(int(netIn)), "white", "", false, false)
		rs += Colorize(parseRepeatSpace(strconv.Itoa(int(netIn)), 7), "white", "", false, false)
	} else if netIn/1024/1024 < 1.0 && netIn/1024 >= 1.0 {
		// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(netIn)/1024)))+strconv.Itoa(int(netIn)/1024)+"k", "yellow", "", false, false)
		rs += Colorize(parseRepeatSpace(strconv.Itoa(int(netIn)/1024), 6)+"k", "yellow", "", false, false)
	}

	if netOut/1024/1024 >= 1.0 {
		// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(netOut)/1024/1024, 1)))+floatToString(float64(netOut)/1024/1024, 1)+"m", "red", "", false, true)
		rs += Colorize(parseRepeatSpace(floatToString(float64(netOut)/1024/1024, 1), 6)+"m", "red", "", false, true)
	} else if netOut/1024 < 1.0 {
		// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(netOut))))+strconv.Itoa(int(netOut)), "white", "", false, false)
		rs += Colorize(parseRepeatSpace(strconv.Itoa(int(netOut)), 7), "white", "", false, false)
	} else if netOut/1024/1024 < 1.0 && netOut/1024 >= 1.0 {
		// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(netOut)/1024)))+strconv.Itoa(int(netOut)/1024)+"k", "yellow", "", false, false)
		rs += Colorize(parseRepeatSpace(strconv.Itoa(int(netOut)/1024), 6)+"k", "yellow", "", false, false)
	}

	if detail == false {
		packetsIn := float64(after.PacketsRecv - beforeNet.PacketsRecv)
		packetsOut := float64(after.PacketsSent - beforeNet.PacketsSent)

		if packetsIn/1000/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(packetsIn/1000/1000, 1)))+floatToString(packetsIn/1000/1000, 1)+"m", "red", "", false, true)
			rs += Colorize(parseRepeatSpace(floatToString(packetsIn/1000/1000, 1), 6)+"m", "red", "", false, true)
		} else if packetsIn/1000 < 1.0 {
			// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(packetsIn))))+strconv.Itoa(int(packetsIn)), "white", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(packetsIn)), 7), "white", "", false, false)
		} else if packetsIn/1000/1000 < 1.0 && packetsIn/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(packetsIn)/1000)))+strconv.Itoa(int(packetsIn)/1000)+"k", "yellow", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(packetsIn)/1000), 6)+"k", "yellow", "", false, false)
		}

		if packetsOut/1000/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(packetsOut)/1000/1000, 1)))+floatToString(float64(packetsOut)/1000/1000, 1)+"m", "red", "", false, true)
			rs += Colorize(parseRepeatSpace(floatToString(float64(packetsOut)/1000/1000, 1), 6)+"m", "red", "", false, true)
		} else if packetsOut/1000 < 1.0 {
			// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(packetsOut))))+strconv.Itoa(int(packetsOut)), "white", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(packetsOut)), 7), "white", "", false, false)
		} else if packetsOut/1000/1000 < 1.0 && packetsOut/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(packetsOut)/1000)))+strconv.Itoa(int(packetsOut)/1000)+"k", "yellow", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(packetsOut)/1000), 6)+"k", "yellow", "", false, false)
		}

		errIn := float64(after.Errin - beforeNet.Errin)
		errOut := float64(after.Errout - beforeNet.Errout)

		if errIn/1000/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(errIn/1000/1000, 1)))+floatToString(errIn/1000/1000, 1)+"m", "red", "", false, true)
			rs += Colorize(parseRepeatSpace(floatToString(errIn/1000/1000, 1), 6)+"m", "red", "", false, true)
		} else if errIn/1000 < 1.0 {
			// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(errIn))))+strconv.Itoa(int(errIn)), "white", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(errIn)), 7), "white", "", false, false)
		} else if errIn/1000/1000 < 1.0 && errIn/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(errIn)/1000)))+strconv.Itoa(int(errIn)/1000)+"k", "yellow", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(errIn)/1000), 6)+"k", "yellow", "", false, false)
		}

		if errOut/1000/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(errOut)/1000/1000, 1)))+floatToString(float64(errOut)/1000/1000, 1)+"m", "red", "", false, true)
			rs += Colorize(parseRepeatSpace(floatToString(float64(errOut)/1000/1000, 1), 6)+"m", "red", "", false, true)
		} else if errOut/1000 < 1.0 {
			// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(errOut))))+strconv.Itoa(int(errOut)), "white", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(errOut)), 7), "white", "", false, false)
		} else if errOut/1000/1000 < 1.0 && errOut/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(errOut)/1000)))+strconv.Itoa(int(errOut)/1000)+"k", "yellow", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(errOut)/1000), 7)+"k", "yellow", "", false, false)
		}

		dropIn := float64(after.Dropin - beforeNet.Dropin)
		dropOut := float64(after.Dropout - beforeNet.Dropout)

		if dropIn/1000/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(dropIn/1000/1000, 1)))+floatToString(dropIn/1000/1000, 1)+"m", "red", "", false, true)
			rs += Colorize(parseRepeatSpace(floatToString(dropIn/1000/1000, 1), 6)+"m", "red", "", false, true)
		} else if dropIn/1000 < 1.0 {
			// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(dropIn))))+strconv.Itoa(int(dropIn)), "white", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(dropIn)), 7), "white", "", false, false)
		} else if dropIn/1000/1000 < 1.0 && dropIn/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(dropIn)/1000)))+strconv.Itoa(int(dropIn)/1000)+"k", "yellow", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(dropIn)/1000), 6)+"k", "yellow", "", false, false)
		}

		if dropOut/1000/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(dropOut)/1000/1000, 1)))+floatToString(float64(dropOut)/1000/1000, 1)+"m", "red", "", false, true)
			rs += Colorize(parseRepeatSpace(floatToString(float64(dropOut)/1000/1000, 1), 6)+"m", "red", "", false, true)
		} else if dropOut/1000 < 1.0 {
			// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(dropOut))))+strconv.Itoa(int(dropOut)), "white", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(dropOut)), 7), "white", "", false, false)
		} else if dropOut/1000/1000 < 1.0 && dropOut/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(dropOut)/1000)))+strconv.Itoa(int(dropOut)/1000)+"k", "yellow", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(dropOut)/1000), 6)+"k", "yellow", "", false, false)
		}

		ffIn := float64(after.Fifoin - beforeNet.Fifoin)
		ffOut := float64(after.Fifoout - beforeNet.Fifoout)

		if ffIn/1000/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(ffIn/1000/1000, 1)))+floatToString(ffIn/1000/1000, 1)+"m", "red", "", false, true)
			rs += Colorize(parseRepeatSpace(floatToString(ffIn/1000/1000, 1), 6)+"m", "red", "", false, true)
		} else if ffIn/1000 < 1.0 {
			// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(ffIn))))+strconv.Itoa(int(ffIn)), "white", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(ffIn)), 7), "white", "", false, false)
		} else if ffIn/1000/1000 < 1.0 && ffIn/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(ffIn)/1000)))+strconv.Itoa(int(ffIn)/1000)+"k", "yellow", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(ffIn)/1000), 6)+"k", "yellow", "", false, false)
		}

		if ffOut/1000/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(ffOut)/1000/1000, 1)))+floatToString(float64(ffOut)/1000/1000, 1)+"m", "red", "", false, true)
			rs += Colorize(parseRepeatSpace(floatToString(float64(ffOut)/1000/1000, 1), 6)+"m", "red", "", false, true)
		} else if ffOut/1000 < 1.0 {
			// rs += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(ffOut))))+strconv.Itoa(int(ffOut)), "white", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(ffOut)), 7), "white", "", false, false)
		} else if ffOut/1000/1000 < 1.0 && ffOut/1000 >= 1.0 {
			// rs += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(ffOut)/1000)))+strconv.Itoa(int(ffOut)/1000)+"k", "yellow", "", false, false)
			rs += Colorize(parseRepeatSpace(strconv.Itoa(int(ffOut)/1000), 6)+"k", "yellow", "", false, false)
		}
	}

	rs += Colorize("|", "dgreen", "", false, false)
	beforeNet = after
	return rs, err
}
