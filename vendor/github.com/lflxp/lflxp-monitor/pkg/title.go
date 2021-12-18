package pkg

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

func GetComTitle() string {
	return Colorize("-------QPS----------TPS------- ", "dgreen", "blue", false, false)
}

func GetComColumns() string {
	return Colorize("  ins   upd   del    sel   iud|", "dgreen", "", true, false)
}

func GetHitTitle() string {
	return Colorize("----KeyBuffer------Index----Qcache---Innodb---(%) ", "dgreen", "blue", false, false)
}

func GetHitColumns() string {
	return Colorize("  read  write    cur  total lorhit readreq  inhit|", "dgreen", "", true, false)
}

func GetInnodbRowsTitle() string {
	return Colorize("---innodb rows status--- ", "dgreen", "blue", false, false)
}

func GetInnodbRowsColumns() string {
	return Colorize("  ins   upd   del   read|", "dgreen", "", true, false)
}

func GetInnodbPagesTitle() string {
	return Colorize("---innodb bp pages status-- ", "dgreen", "blue", false, false)
}

func GetInnodbPagesColumns() string {
	return Colorize("   data   free  dirty flush|", "dgreen", "", true, false)
}

func GetInnodbDataTitle() string {
	return Colorize("-----innodb data status----- ", "dgreen", "blue", false, false)
}

func GetInnodbDataColumns() string {
	return Colorize(" reads writes   read written|", "dgreen", "", true, false)
}

func GetInnodbLogTitle() string {
	return Colorize("--innodb log-- ", "dgreen", "blue", false, false)
}

func GetInnodbLogColumns() string {
	return Colorize("fsyncs written|", "dgreen", "", true, false)
}

func GetInnodbStatusTitle() string {
	return Colorize("  his --log(byte)--  read ---query--- ", "dgreen", "blue", false, false)
}

func GetInnodbStatusColumns() string {
	return Colorize(" list uflush  uckpt  view inside  que|", "dgreen", "", true, false)
}

func GetThreadsTitle() string {
	return Colorize("----------threads--------- ", "dgreen", "blue", false, false)
}

func GetThreadsColumns() string {
	return Colorize(" run  con  cre  cac   "+"%"+"hit|", "dgreen", "", true, false)
}

func GetBytesTitle() string {
	return Colorize("-----bytes---- ", "dgreen", "blue", false, false)
}

func GetBytesColumns() string {
	return Colorize("   recv   send|", "dgreen", "", true, false)
}

func GetSemiTitle() string {
	return Colorize("---avg_wait--tx_times--semi ", "dgreen", "blue", false, false)
}

func GetSemiColumns() string {
	return Colorize("  naw  txaw notx  yes   off|", "dgreen", "", true, false)
}

func GetSlaveTitle() string {
	return Colorize("---------------SlaveStatus------------- ", "dgreen", "blue", false, false)
}

func GetSlaveColumns() string {
	return Colorize("ReadMLP ExecMLP   chkRE   SecBM|", "dgreen", "", true, false)
}
