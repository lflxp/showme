package main

import (
	"fmt"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/lflxp/showme/completers"
	"github.com/lflxp/showme/executors"
)

// 实时左标显示
var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix, LivePrefixState.IsEnable
}

// 功能筛选信息
func completer(in prompt.Document) []prompt.Suggest {
	// 常规过滤
	// 如果输入值为空 返回空字符串
	if in.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	// 	获取所有输入字符串并以空格分割
	args := strings.Split(in.TextBeforeCursor(), " ")
	// 获取当前输入字符
	current := in.GetWordBeforeCursor()

	// If PIPE is in text before the cursor, returns empty suggestions.
	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}

	// If word before the cursor starts with "-", returns CLI flag options.
	if strings.HasPrefix(current, "-") {
		return completers.OptionsCompleters(args, strings.HasPrefix(current, "--"))
	}

	// 功能列表 排除包含“-”的字符，遇到-则返回空交由下面函数处理
	// if suggests, found := completers.GlobalOptionFunc(in); found {
	// 	return suggests
	// }

	// 输入即取消提示
	// 非常规过滤

	return completers.FirstCommandFunc(in, args)
}

// 执行命令过程
func executor(in string) {
	fmt.Println("Your input: " + in)
	if in == "" {
		LivePrefixState.IsEnable = false
		LivePrefixState.LivePrefix = in
		return
	} else {
		thisisit, status := executors.ParseExecutors(in)
		if status {
			thisisit()
		}
	}
	LivePrefixState.LivePrefix = in + "> "
	LivePrefixState.IsEnable = true
	// fmt.Println("executor")
}

func main() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("live-prefix-example"),
	)
	p.Run()
}
