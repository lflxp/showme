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
