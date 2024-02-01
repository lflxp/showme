package pkg

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/host"
)

func GetHostInfo() error {
	data, err := host.SensorsTemperatures()
	if err != nil {
		return err
	}

	for _, x := range data {
		fmt.Println(fmt.Sprintf("%v: %.2f", Colorize(x.SensorKey, "white", "red", false, false), x.Temperature))
	}
	return nil
}
