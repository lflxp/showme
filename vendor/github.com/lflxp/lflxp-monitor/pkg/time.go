package pkg

func TimeNow() (string, error) {
	return Colorize(GetNowTime(), "yellow", "", false, false) + Colorize("|", "dgreen", "", false, false), nil
}
