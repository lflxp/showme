package suggests

import "github.com/c-bata/go-prompt"

var Commands = []prompt.Suggest{
	{Text: "dashboard", Description: "当前计算机状态总览"},
	{Text: "gocui", Description: "https://github.com/jroimartin/gocui"},
	{Text: "monitor", Description: "monitoring Linux/Unix or MacOs status runtime"},
}
