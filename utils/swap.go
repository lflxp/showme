package utils

// swap in mem functions
import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/mem"
)

var beforeSwap *MonitorSwap

func init() {
	beforeSwap, err = NewSwap()
}

type MonitorSwap struct {
	swap_in  uint64
	swap_out uint64
}

func (this *MonitorSwap) Get() error {
	data, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	this.swap_in = data.Sin
	this.swap_out = data.Sout

	return nil
}

func NewSwap() (*MonitorSwap, error) {
	data := &MonitorSwap{}
	err := data.Get()
	return data, err
}

func SwapIO() (string, error) {
	var rs string
	after, err := NewSwap()
	if err != nil {
		return rs, err
	}
	si := after.swap_in - beforeSwap.swap_in
	so := after.swap_out - beforeSwap.swap_out

	in := strings.Repeat(" ", 5-len(fmt.Sprintf("%d", si))) + fmt.Sprintf("%d", si)
	out := strings.Repeat(" ", 5-len(fmt.Sprintf("%d", so))) + fmt.Sprintf("%d", so)

	if si > 0 {
		rs += Colorize(in, "red", "", false, true)
	} else {
		rs += Colorize(in, "white", "", false, false)
	}

	if so > 0 {
		rs += Colorize(out, "red", "", false, true)
	} else {
		rs += Colorize(out, "white", "", false, false)
	}

	rs += Colorize("|", "dgreen", "", false, false)
	beforeSwap = after
	return rs, nil
}
