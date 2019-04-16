package gopacket

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
)

// func Run(in string) {
// 	fmt.Println(in)
// 	tmp := strings.Split(in, " ")
// 	if len(tmp) != 3 {
// 		fmt.Println("命令长度不为三，退出")
// 		return
// 	}

// 	// ok := true
// 	// // 获取退出信号
// 	// c := make(chan os.Signal, 1)
// 	// signal.Notify(c, os.Interrupt, os.Kill)
// 	// num := 0
// 	// go func() {
// 	// 	for {
// 	// 		num++
// 	// 		select {
// 	// 		case s := <-c:
// 	// 			fmt.Printf("\n\033[1;4;31m%s:罒灬罒:小伙子走了哟！\033[0m\n", s)
// 	// 			ok = false
// 	// 			break
// 	// 			os.Exit(1)
// 	// 		}
// 	// 		if !ok {
// 	// 			break
// 	// 		}
// 	// 	}
// 	// }()
// 	utils.WatchDogEasy(tmp[2])
// }

var ined string

func Run(in string) {
	fmt.Println(in)
	tmp := strings.Split(in, " ")
	if len(tmp) != 3 {
		fmt.Println("命令长度不为三，退出")
		return
	}

	ok := true
	// 获取退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	bigChan := make(chan interface{}, 1000)
	// defer close(bigChan)

	go utils.WatchDog(bigChan, tmp[2])

	num := 0
	for {
		num++
		select {
		case s := <-c:
			fmt.Printf("\n\033[1;4;31m%s:罒灬罒:小伙子走了哟！\033[0m\n", s)
			ok = false
			close(bigChan)
			break
		case data, ok := <-bigChan:
			if !ok {
				fmt.Println("is ok?")
				return
			}
			switch data.(type) {
			case *utils.Data:
				json, err := json.Marshal(data.(*utils.Data))
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(string(json))
			}
		}

		if !ok {
			break
		}
		fmt.Println("num", num)
	}
}

func Gopacket(in string) {
	ined = in
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(dlayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func dquit(g *gocui.Gui, v *gocui.View) error {
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

	// 获取退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	bigChan := make(chan interface{}, 1000)
	defer close(bigChan)
	ok := true
	go utils.WatchDog(bigChan, strings.Split(ined, " ")[2])

	num := 0
	for {
		num++
		select {
		case s := <-c:
			fmt.Printf("\n\033[1;4;31m%s:罒灬罒:小伙子走了哟！\033[0m\n", s)
			ok = false
			close(bigChan)
			break
		case data, ok := <-bigChan:
			if !ok {
				fmt.Println("is ok?")
				return nil
			}
			switch data.(type) {
			case *utils.Data:
				// tmp_data := data.(*utils.Data)

				if v, err := g.SetView("v1", 0, 0, maxX/2-1, maxY); err != nil {
					if err != gocui.ErrUnknownView {
						return err
					}
					v.Title = "Num"
					v.Wrap = true
					v.Autoscroll = true
					v.Editable = false

					v.Clear()
					fmt.Fprintln(v, num)

					if _, err = setCurrentViewOnTop(g, "v1"); err != nil {
						return err
					}
				}

				if v, err := g.SetView("v2", maxX/2-1, 0, maxX-1, maxY); err != nil {
					if err != gocui.ErrUnknownView {
						return err
					}
					v.Title = "网络抓包"
					v.Wrap = true
					v.Autoscroll = true
					v.Editable = false

					json, err := json.Marshal(data.(*utils.Data))
					if err != nil {
						fmt.Fprintln(v, err.Error())
						return nil
					}
					fmt.Fprintln(v, string(json))
				}
			}
		}

		if !ok {
			break
		}
		fmt.Println("num", num)
	}

	return nil
}
