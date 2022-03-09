package suggests

import "github.com/c-bata/go-prompt"

func DetailSuggest(in prompt.Document) []prompt.Suggest {
	rs := []prompt.Suggest{
		{Text: "-v", Description: "详细信息 【informational】"},
		{Text: "-vvv", Description: "DEBUG模式，最详细信息"},
	}
	return prompt.FilterHasPrefix(rs, in.GetWordBeforeCursor(), true)
}

var GlobalOptions = []prompt.Suggest{
	{Text: "-v", Description: "详细信息 【informational】"},
	{Text: "-vvv", Description: "DEBUG模式，最详细信息"},
	{Text: "-h", Description: "查看帮助简写"},
	{Text: "--help", Description: "查看帮助"},
}
