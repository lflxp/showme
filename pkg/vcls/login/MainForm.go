//go:build govcl

package login

import (
	"github.com/ying32/govcl/vcl"
)

type TMainForm struct {
	*vcl.TForm
	Button1 *vcl.TButton

	//::private::
	TMainFormFields
}

var MainForm *TMainForm

// Loaded in bytes.
// vcl.Application.CreateForm(&MainForm)

func NewMainForm(owner vcl.IComponent) (root *TMainForm) {
	vcl.CreateResForm(owner, &root)
	return
}

var MainFormBytes = []byte("\x54\x50\x46\x30\x0B\x54\x44\x65\x73\x69\x67\x6E\x46\x6F\x72\x6D\x08\x4D\x61\x69\x6E\x46\x6F\x72\x6D\x04\x4C\x65\x66\x74\x02\x08\x06\x48\x65\x69\x67\x68\x74\x03\x55\x02\x03\x54\x6F\x70\x02\x08\x05\x57\x69\x64\x74\x68\x03\xB1\x03\x07\x43\x61\x70\x74\x69\x6F\x6E\x06\x09\xE4\xB8\xBB\xE7\xAA\x97\xE5\x8F\xA3\x0C\x43\x6C\x69\x65\x6E\x74\x48\x65\x69\x67\x68\x74\x03\x55\x02\x0B\x43\x6C\x69\x65\x6E\x74\x57\x69\x64\x74\x68\x03\xB1\x03\x05\x43\x6F\x6C\x6F\x72\x07\x09\x63\x6C\x42\x74\x6E\x46\x61\x63\x65\x0A\x46\x6F\x6E\x74\x2E\x43\x6F\x6C\x6F\x72\x07\x0C\x63\x6C\x57\x69\x6E\x64\x6F\x77\x54\x65\x78\x74\x0B\x46\x6F\x6E\x74\x2E\x48\x65\x69\x67\x68\x74\x02\xF3\x09\x46\x6F\x6E\x74\x2E\x4E\x61\x6D\x65\x06\x06\x54\x61\x68\x6F\x6D\x61\x08\x50\x6F\x73\x69\x74\x69\x6F\x6E\x07\x0E\x70\x6F\x53\x63\x72\x65\x65\x6E\x43\x65\x6E\x74\x65\x72\x00\x07\x54\x42\x75\x74\x74\x6F\x6E\x07\x42\x75\x74\x74\x6F\x6E\x31\x04\x4C\x65\x66\x74\x03\x91\x01\x06\x48\x65\x69\x67\x68\x74\x02\x25\x03\x54\x6F\x70\x03\x07\x01\x05\x57\x69\x64\x74\x68\x02\x5E\x07\x43\x61\x70\x74\x69\x6F\x6E\x06\x07\x42\x75\x74\x74\x6F\x6E\x31\x08\x54\x61\x62\x4F\x72\x64\x65\x72\x02\x00\x00\x00\x00")

// 注册窗口资源
var _ = vcl.RegisterFormResource(MainForm, &MainFormBytes)
