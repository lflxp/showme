package pkg

import "errors"

type Lflxp interface {
	Check() error
	Execute() error
}

type Apis struct {
	Host  string // http bind host
	Port  string // http bind port
	Stats bool   // is output db stats
}

func (this *Apis) Check() error {
	if this.Host == "" || this.Port == "" {
		return errors.New("host or port is none")
	}
	return nil
}

func (this *Apis) Execute() error {
	err := this.Check()
	if err != nil {
		return err
	}

	Api(this)
	return nil
}
