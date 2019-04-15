package utils

func GetCpuTitle() string {
	return Colorize("----cpu-usage--- ", "dgreen", "", false, false)
}

func GetCpuColumns() string {
	return Colorize(" usr sys idl iow|", "dgreen", "", true, false)
}

func GetTimeTitle() string {
	return Colorize("-------- ", "dgreen", "", false, false)
}

func GetTimeColumns() string {
	return Colorize("  time  |", "dgreen", "", true, false)
}

func GetLoadTitle() string {
	return Colorize("-----load-avg---- ", "dgreen", "", false, false)
}

func GetLoadColumns() string {
	return Colorize("  1m    5m   15m |", "dgreen", "", true, false)
}

func GetSwapTitle() string {
	return Colorize("---swap--- ", "dgreen", "", false, false)
}

func GetSwapColumns() string {
	return Colorize("   si   so|", "dgreen", "", true, false)
}

func GetNetTitle(detail bool) string {
	var rs string
	if detail {
		rs = Colorize("----net(A)---- ", "dgreen", "", false, false)
	} else {
		rs = Colorize("------------------------------net(Detail)----------------------------- ", "dgreen", "", false, false)
	}
	return rs
}

func GetNetColumns(detail bool) string {
	var rs string
	if detail {
		rs = Colorize("   recv   send|", "dgreen", "", true, false)
	} else {
		rs = Colorize("   recv   send   psin   psot  errin  errot   dpin  dpout   ffin  ffout|", "dgreen", "", true, false)
	}
	return rs
}

func GetDiskTitle() string {
	return Colorize("------------------------io-usage---------------------- ", "dgreen", "", false, false)
}

func GetDiskColumns() string {
	return Colorize(" readc writec    srkB    swkB queue  await svctm %util|", "dgreen", "", true, false)
}
