package completers

import (
	"fmt"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/lflxp/showme/suggests"
	"github.com/lflxp/showme/utils"
)

// 解析函数 判断最新参数是否含有-字符
func getPreviousOption(d prompt.Document) (cmd, option string, found bool) {
	args := strings.Split(d.TextBeforeCursor(), " ")
	l := len(args)
	if l >= 2 {
		option = args[l-2]
	}
	if strings.HasPrefix(option, "-") {
		return args[0], option, true
	}
	return "", "", false
}

// 全局固定命令
func GlobalOptionFunc(d prompt.Document) ([]prompt.Suggest, bool) {
	cmd, option, found := getPreviousOption(d)
	if !found {
		return []prompt.Suggest{}, false
	}

	switch cmd {
	case "dashboard":
		// 带 - 参数的命令提示
		// 命令输入大于等于两个
		switch option {
		case "-v", "-vvv":
			return prompt.FilterHasPrefix(
				suggests.DetailSuggest(d),
				d.GetWordBeforeCursor(),
				true,
			), true
		}
	}
	return []prompt.Suggest{}, false
}

// 用户自定义命令
func FirstCommandFunc(d prompt.Document, args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(Commands, args[0], true)
	}

	first := args[0]
	switch first {
	case "dashboard":
		second := args[1]
		if len(args) == 2 {
			subcommands := []prompt.Suggest{
				{Text: "help"},
				{Text: "show", Description: "console for show me"},
				{Text: "helloworld", Description: "dashboard for tcell cellviews.go"},
			}
			return prompt.FilterHasPrefix(subcommands, second, true)
		}
	case "gocui":
		second := args[1]
		if len(args) == 2 {
			subcommands := []prompt.Suggest{
				{Text: "help", Description: "帮助文档"},
				{Text: "active", Description: "界面布局"},
			}
			return prompt.FilterHasPrefix(subcommands, second, true)
		}
	case "gopacket":
		second := args[1]
		if len(args) == 2 {
			subcommands := []prompt.Suggest{
				{Text: "interface", Description: "制定监听网卡"},
				{Text: "screen", Description: "gocui 可视化"},
			}
			return prompt.FilterHasPrefix(subcommands, second, true)
		}

		third := args[2]
		if len(args) == 3 {
			switch second {
			case "interface", "in", "screen":
				interfaces, err := utils.GetCurrentInterfaceCommands()
				if err != nil {
					fmt.Println(err.Error())
				}
				return prompt.FilterContains(interfaces, third, true)
			}
		}
	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}
