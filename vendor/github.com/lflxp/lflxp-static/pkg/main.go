package pkg

import "errors"

type Lflxp interface {
	Check() error
	Execute() error
}

type Apis struct {
	Port     string
	Path     string
	Types    string
	IsVideo  bool
	PageSize int
}

func (this *Apis) Check() error {
	if this.Path == "" || this.Types == "" {
		return errors.New("path or types is wrong")
	}
	return nil
}

func (this *Apis) Execute() error {
	HttpStaticServeForCorba(this)
	return nil
}
