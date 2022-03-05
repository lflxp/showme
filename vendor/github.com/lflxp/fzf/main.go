package main

import (
	fzf "github.com/lflxp/fzf/src"
	"github.com/lflxp/fzf/src/protector"
)

var version string = "0.29"
var revision string = "devel"

func main() {
	protector.Protect()
	fzf.Run(fzf.ParseOptions(), version, revision)
}
