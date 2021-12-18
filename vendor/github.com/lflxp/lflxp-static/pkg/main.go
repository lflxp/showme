package pkg

import "errors"

type Lflxp interface {
	Check() error
	Execute() error
}

type Apis struct {
	StaticPort string
	Port       string
	Path       string
	Types      string
	IsVideo    bool
	PageSize   int
	Raw        bool
}

func (this *Apis) Check() error {
	if this.Path == "" || this.Types == "" {
		return errors.New("path or types is wrong")
	}
	return nil
}

func (this *Apis) Execute() error {
	err := this.Check()
	if err != nil {
		return err
	}

	HttpStaticServeForCorba(this)
	return nil
}
