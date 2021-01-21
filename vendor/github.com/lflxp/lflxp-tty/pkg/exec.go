package pkg

type Lflxp interface {
	Check() error
	Execute() error
}

type Tty struct {
	EnableTLS      bool
	CrtPath        string
	KeyPath        string
	IsProf         bool
	IsXsrf         bool
	IsAudit        bool
	IsPermitWrite  bool
	MaxConnections int64
	IsReconnect    bool
	IsDebug        bool
	Username       string
	Password       string
	Port           string
	Host           string
	Cmds           []string
}

func (this *Tty) Check() error {
	if len(this.Cmds) == 0 {
		this.Cmds = append(this.Cmds, "bash")
	}
	return nil
}

func (this *Tty) Execute() error {
	err := this.Check()
	if err != nil {
		return err
	}

	ServeGin(this)
	return nil
}
