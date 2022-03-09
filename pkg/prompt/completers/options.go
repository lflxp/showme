package completers

import (
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/lflxp/showme/pkg/prompt/suggests"
)

var optionHelp = []prompt.Suggest{
	{Text: "-h"},
	{Text: "--help"},
}

// 剔除可能重复的全局参数 -- or -
func excludeOptions(args []string) ([]string, bool) {
	l := len(args)
	filtered := make([]string, 0, l)
	skipList := []string{
		"-v",
		"-vvv",
		"-h",
		"--help",
	}

	var skipListArg bool
	for i := 0; i < len(args); i++ {
		if skipListArg {
			skipListArg = false
			continue
		}

		// 剔重复参数
		isExist := false
		for _, s := range skipList {
			if strings.HasPrefix(args[i], s) {
				if strings.Contains(args[i], "=") {
					skipListArg = false
				} else {
					skipListArg = true
				}
				continue
			}
			if args[i] == s {
				isExist = true
			}
		}
		if strings.HasPrefix(args[i], "-") {
			continue
		}

		if !isExist {
			filtered = append(filtered, args[i])
		}
	}
	return filtered, skipListArg
}

// remove exists
func removeGlobalExist(args []string) []prompt.Suggest {
	rs := []prompt.Suggest{}

	for _, a := range suggests.GlobalOptions {
		status := false
		for _, b := range args {
			if a.Text == b {
				status = true
			}
		}
		if !status {
			rs = append(rs, a)
		}
	}
	return rs
}

// long == --
func OptionsCompleters(args []string, long bool) []prompt.Suggest {
	l := len(args)
	if l <= 1 {
		if long {
			return prompt.FilterHasPrefix(optionHelp, "--", false)
		}
		return optionHelp
	}

	var sug []prompt.Suggest
	//TODO:  过滤 - or -- 参数？ 比如 -f --namespace
	commandArgs, _ := excludeOptions(args)
	switch commandArgs[0] {
	case "dashboard":
		sug = suggests.DashboardOptions
	case "monitor":
		sug = suggests.MonitorOptions
	case "mysql":
		sug = suggests.MysqlOptions
	case "tty":
		sug = suggests.TtyOptions
	default:
		sug = optionHelp
	}

	sug = append(sug, removeGlobalExist(args)...)
	if long {
		return prompt.FilterContains(
			prompt.FilterHasPrefix(sug, "--", false),
			strings.TrimLeft(args[l-1], "--"),
			true,
		)
	}
	return prompt.FilterContains(sug, strings.TrimLeft(args[l-1], "-"), true)
}
