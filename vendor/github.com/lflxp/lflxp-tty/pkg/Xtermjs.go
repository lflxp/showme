package pkg

import (
	"sync"
)

// TODO: 标准化前端输入类型
// xtermjs eventListener
const (
	Input          = '0'
	Ping           = '1'
	ResizeTerminal = '2'
	Heartbeat      = '3'
)

// 后端控制前端
const (
	Output         = '0'
	Pong           = '1'
	SetWindowTitle = '2'
	SetPreferences = '3'
	SetReconnect   = '4'
)

type Options struct {
	PermitWrite      bool
	MaxConnections   int64
	CloseSignal      int
	Audit            bool
	Xsrf             bool
	EnableTLS        bool
	CrtPath, KeyPath string
	IsReconnect      bool
	IsDebug          bool
}

// 原本是命令端http server管理，这里后期可以改成gin server管理
// 单线程等待
type Server struct {
	wg *sync.WaitGroup
}

func (this *Server) StartGo() {
	this.wg.Add(1)
}

func (this *Server) DoneGo() {
	this.wg.Done()
}

func (this *Server) WaitGo() {
	this.wg.Wait()
}

func NewServer() *Server {
	return &Server{
		wg: new(sync.WaitGroup),
	}
}

// xtermjs前端配置
type XtermJs struct {
	Title       string
	Server      *Server
	Options     Options
	Connections *int64   // 统计连接数
	XsrfToken   sync.Map // xsrftoken存储
	Cmds        []string // 命令集
}
