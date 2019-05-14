package suggests

import "github.com/c-bata/go-prompt"

var HttpStaticOptions = []prompt.Suggest{
	{Text: "-path", Description: "静态文件路径，默认： ./"},
	{Text: "-port", Description: "静态文件服务端口,默认：9090"},
}
