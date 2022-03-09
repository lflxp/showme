package suggests

import "github.com/c-bata/go-prompt"

var DashboardOptions = []prompt.Suggest{
	{Text: "--all", Description: "显示所有信息"},
	{Text: "--cpu", Description: "显示cpu信息"},
	{Text: "--mem", Description: "显示内存信息"},
	{Text: "--disk", Description: "显示磁盘信息"},
}
