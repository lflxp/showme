package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	types     []string
	pageSize  int
	isvideo   bool
	path      string
	port      string
	closeChan chan os.Signal
	uri       string
	initnum   int
)

const maxUploadSize = 2000 * 1024 * 2014 // 2 MB
const uploadPath = "/tmp"

func init() {
	initnum = 0
	// path = utils.GetCurrentDirectory()
	// port = "9090"
	closeChan = make(chan os.Signal)
}

func HttpStaticServeForCorba(data *Apis) {
	// httpstatic -port 9090 -path ./
	port = data.Port
	path = data.Path
	isvideo = data.IsVideo
	pageSize = data.PageSize
	types = strings.Split(data.Types, ",")

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(dlayout)

	d := time.Duration(time.Second)
	t := time.NewTicker(d)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				g.Update(func(g *gocui.Gui) error { return nil })
				// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				// fmt.Fprintln(v, )
			}
		}
	}()

	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
	// 	log.Panicln(err)
	// }

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func HttpStaticServe(in string) {
	// httpstatic -port 9090 -path ./
	tmp := strings.Split(in, " ")
	for n, x := range tmp {
		if x == "-port" {
			port = tmp[n+1]
		}
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(dlayout)

	d := time.Duration(time.Second)
	t := time.NewTicker(d)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				g.Update(func(g *gocui.Gui) error { return nil })
				// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				// fmt.Fprintln(v, )
			}
		}
	}()

	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
	// 	log.Panicln(err)
	// }

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func keybindings(g *gocui.Gui) error {
	// 清空side缓存
	// if err := g.SetKeybinding("help", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		return err
	}
	return nil
}

func dquit(g *gocui.Gui, v *gocui.View) error {
	closeChan <- syscall.SIGINT
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func dlayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("history", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "访问记录"
		v.Wrap = true
		v.Frame = false
		// v.Highlight = true
		v.Autoscroll = true
		v.SelFgColor = gocui.ColorYellow
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// uri = fmt.Sprintf("/a%s", time.Now().Format("150405"))
	}
	ips := GetIPs()
	if v, err := g.SetView("top", maxX/2-80, maxY/2, maxX/2+80, maxY/2+2*len(ips)+4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if isvideo {
			v.Title = "视频服务器地址"
		} else {
			v.Title = "静态服务器地址"
		}

		v.Wrap = true
		// v.Highlight = true
		// v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// uri = fmt.Sprintf("/a%s", time.Now().Format("150405"))

		urls := []string{}

		for _, ip := range ips {
			urls = append(urls, fmt.Sprintf("UploadURL: => %s:%s <= PATH: => /upload <= ", ip, port))
		}
		for _, ip := range ips {
			urls = append(urls, fmt.Sprintf("DownURL: => %s:%s <= PATH: => / <= ", ip, port))
		}
		dir, _ := os.Getwd()
		urls = append(urls, fmt.Sprintf("CurrentPWD: => %s <= ", dir))
		urls = append(urls, fmt.Sprintf("UploadDIR: => %s <= ", dir))
		urls = append(urls, "curl -X POST http://127.0.0.1:9090/upload -F \"file=@/Users/lxp/123.mp4\" -H \"Content-Type:multipart/form-data\"")
		fmt.Fprintln(v, strings.Join(urls, "\n"))
		go serverGin(g)

		if _, err = setCurrentViewOnTop(g, "top"); err != nil {
			return err
		}
	}
	return nil
}
