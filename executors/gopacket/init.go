// +build gopacket
package gopacket

import (
	"github.com/c-bata/go-prompt"
	"github.com/lflxp/showme/completers"
)

func init() {
	completers.Commands = append(completers.Commands, prompt.Suggest{Text: "gopacket", Description: "网络抓包工具"})
}
