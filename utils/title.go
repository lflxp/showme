package utils

func GetCpuTitle() string {
	return Colorize("---cpu-usage--- ", "dgreen", "", false, false)
}

func GetCpuColumns() string {
	return Colorize("usr sys idl iow|", "dgreen", "", true, false)
}
