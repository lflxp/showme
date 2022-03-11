//go:build govcl

package login

import (
	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
)

func Run() {

	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&MainForm)
	vcl.Application.CreateForm(&LoginForm)

	// 这里可以决定是不是显示主窗口
	vcl.Application.SetShowMainForm(false)
	// 这里显示出第二个窗口
	LoginForm.Show()

	vcl.Application.Run()

}
