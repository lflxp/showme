package search

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"gitee.com/lflxp/tcell/v2"
	"github.com/briandowns/spinner"
	"github.com/creack/pty"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

// 1. 画搜索框 || 继承fzf tui 纯搜索
// 2. 在input中输入搜索内容
// 3. 下拉列表显示搜索内容
// 4. 选择并显示到console
// 5. 搜索内容真彩色
// 6. 搜索内容背景色和道标
// 7. 搜索内容上下移动+限制显示数量
// 8. dashboard 色塊 github.com/lucasb-eyer/go-colorful
// 9. 进度条 github.com/briandowns/spinner

//************************************************************************************************
//   __test() {
// 	local cmd="${FZF_TEST_COMMAND:-"cd ~/code/go && list"}"
// 	setopt localoptions pipefail no_aliases 2> /dev/null
// 	eval "$cmd"
// 	local ret=$?
// 	echo
// 	return $ret
//   }

//   # test
//   fzf-test() {
// 	LBUFFER="${LBUFFER}$(__test)"
// 	local ret=$?
// 	zle reset-prompt
// 	return $ret
//   }
//   zle     -N   fzf-test
//   bindkey '^[e' fzf-test
//************************************************************************************************

//****************************************************************
// # 更多配置项
// \033[0m 关闭所有属性
// \033[1m 设置高亮度
// \033[4m 下划线
// \033[5m 闪烁
// \033[7m 反显
// \033[8m 消隐
// \033[30m -- \033[37m 设置前景色
// \033[40m -- \033[47m 设置背景色
// \033[nA 光标上移n行
// \033[nB 光标下移n行
// \033[nC 光标右移n行
// \033[nD 光标左移n行
// \033[y;xH设置光标位置
// \033[2J 清屏
// \033[K 清除从光标到行尾的内容
// \033[s 保存光标位置
// \033[u 恢复光标位置
// \033[?25l 隐藏光标
// \033[?25h 显示光标
//********************************

type TuiScreen struct {
	screen       tcell.Screen
	EnableMouse  bool // 开启鼠标可用
	EnablePaste  bool // 开启粘贴板可用
	boxStyle     tcell.Style
	X            int
	Y            int
	Input        []rune   // 键盘输入结果
	Files        []string // 文件查询总结果
	Total        int
	SearchNum    int
	SearchResult []string // 搜索临时结果
	// 二维空间定位
	top        int // 起点x轴坐标
	left       int // 起点y轴坐标
	width      int // 宽度
	height     int // 高度
	MaxContent int // 最多显示条数
	CursorPos  int // 搜索光标位置 0 < CursorPos <= MaxContent
	CurContent string
	origState  *term.State
	ttyin      *os.File
	queued     strings.Builder
	offset     int // 偏移量 相对于顶部
	fd         int
	Unsearch   []string
	finish     chan int // 扫描结束
}

func (t *TuiScreen) init() (err error) {
	t.screen, err = tcell.NewScreen()
	// t.screen, err = tcell.NewConsoleScreen()
	// t.screen, err = tcell.NewTerminfoScreen()
	if err != nil {
		return
	}

	err = t.screen.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "init failed: %v\n", err)
		return
	}

	if t.EnableMouse {
		t.screen.EnableMouse()
	}

	if t.EnablePaste {
		t.screen.EnablePaste()
	}

	t.screen.SetStyle(tcell.StyleDefault)
	t.X = -1
	t.Y = -1
	if t.finish == nil {
		t.finish = make(chan int)
	}
	t.Files = []string{}

	t.boxStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	t.Input = []rune{}

	// t.ttyin = TtyIn()

	// if err := t.initPlatform(); err != nil {
	// 	errorExit(err.Error())
	// }
	// t.updateTerminalSize()
	// t.csi("J")
	// t.makeSpace()
	// t.csi(fmt.Sprintf("%dA", t.MaxContent))
	// t.csi("G")
	// t.csi("K")

	// t.SearchResult = t.Files
	// t.Total = len(files)
	// t.MaxContent = 10
	t.CursorPos = 0
	t.SearchResult = []string{}

	// t.out = os.Stdin
	t.fd = syscall.Stdout
	t.width, t.height = t.screen.Size()
	t.MaxContent = t.height - 4

	return
}

func (r *TuiScreen) Fd() int {
	return int(r.ttyin.Fd())
}

func ExecCommandString(cmd string) (string, error) {
	pipeline := exec.Command("/bin/sh", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	pipeline.Stdout = &out
	pipeline.Stderr = &stderr
	err := pipeline.Run()
	if err != nil {
		return stderr.String(), err
	}
	return out.String(), nil
}

func CommandPty(cmd string) error {
	// Create arbitrary command.
	c := exec.Command("/bin/sh", "-c", cmd)

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err

	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)

			}

		}

	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)

	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil

}

func (t *TuiScreen) GetFiles(dirPath string) error {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	if t.Files == nil {
		t.Files = []string{}
	}

	for _, dir := range dir {
		if dir.IsDir() {
			t.Files = append(t.Files, fmt.Sprintf("Dir %s", dir.Name()))
		} else {
			t.Files = append(t.Files, fmt.Sprintf("File %s", dir.Name()))
		}
		t.Total = t.Total + 1
	}
	return nil
}

// 快速判断字符是否在字符数组中
func In(target string, source []string) bool {
	sort.Strings(source)
	index := sort.SearchStrings(source, target)
	if index < len(source) && source[index] == target {
		return true
	}
	return false
}

func (t *TuiScreen) GetCommand() {
	t.Files = append(t.Files, "CMD ps -ef|showme|awk '{print $2}'|xargs kill -9 ")
	t.Files = append(t.Files, "CMD showme martix")
	t.Files = append(t.Files, "CMD showme cmd")
	t.Files = append(t.Files, "CMD showme static")
	t.Files = append(t.Files, "CMD showme tty -w")
	t.Files = append(t.Files, "CMD showme watch")
	t.Files = append(t.Files, "CMD fzf --height=40% --preview 'cat {}'")
}

func (t *TuiScreen) Gopsutil() error {
	// 获取环境变量
	for _, x := range os.Environ() {
		t.Files = append(t.Files, fmt.Sprintf("ENV %s", x))
	}
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorln("os.Hostname", err.Error())
		return err
	}
	t.Files = append(t.Files, fmt.Sprintf("HOST %s", hostname))
	// 获取内存信息
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Errorln("mem.VirtualMemory", err.Error())
		return err
	}
	t.Files = append(t.Files, fmt.Sprintf("MEM %v", v))
	t.Files = append(t.Files, fmt.Sprintf("MEM Total: %v", v.Total))
	t.Files = append(t.Files, fmt.Sprintf("MEM Available: %v", v.Available))
	t.Files = append(t.Files, fmt.Sprintf("MEM UsedPercent: %v", v.UsedPercent))
	t.Files = append(t.Files, fmt.Sprintf("MEM FREE: %v", v.Free))
	t.Files = append(t.Files, fmt.Sprintf("MEM SwapTotal: %v", v.SwapTotal))
	t.Files = append(t.Files, fmt.Sprintf("MEM SwapFree: %v", v.SwapFree))
	// 获取CPU信息
	physicalCnt, err := cpu.Counts(false)
	if err != nil {
		log.Errorln("cpu.Counts.false", err.Error())
		return err
	}
	logicalCnt, err := cpu.Counts(true)
	if err != nil {
		log.Errorln("cpu.Counts", err.Error())
		return err
	}
	t.Files = append(t.Files, fmt.Sprintf("CPU PHYSICAL %d", physicalCnt))
	t.Files = append(t.Files, fmt.Sprintf("CPU LOGICAL %d", logicalCnt))

	info, err := cpu.Info()
	if err != nil {
		log.Errorln("cpu.Info", err.Error())
		return err
	}
	for i, c := range info {
		t.Files = append(t.Files, fmt.Sprintf("CPU Num %d %v", i, c))
	}

	times, err := cpu.Times(true)
	if err != nil {
		log.Errorln("cpu.Times", err.Error())
		return err
	}
	for i, c := range times {
		t.Files = append(t.Files, fmt.Sprintf("CPU TIMES %d %v", i, c))
	}

	// 磁盘
	mapStat, err := disk.IOCounters()
	if err != nil {
		log.Errorln("disk.IOCounters", err.Error())
		return err
	}

	for name, stat := range mapStat {
		t.Files = append(t.Files, fmt.Sprintf("DISK %s %v", name, stat))
	}

	infos, err := disk.Partitions(true)
	if err != nil {
		log.Errorln("disk.Partitions", err.Error())
		return err
	}
	for i, c := range infos {
		t.Files = append(t.Files, fmt.Sprintf("DISK PARTITIONS %d %v", i, c))
		tmp, err := disk.Usage(c.Mountpoint)
		if err != nil {
			// log.Debug("disk.infos %s %s", c.Mountpoint, err.Error())
		} else {
			t.Files = append(t.Files, fmt.Sprintf("DISK USAGE %v", tmp))
		}
	}
	// 主机信息
	timestamp, err := host.BootTime()
	if err != nil {
		log.Errorln("host.BootTime()", err.Error())
		return err
	}
	t.Files = append(t.Files, fmt.Sprintf("HOST BootTime %s Unix %v", time.Unix(int64(timestamp), 0).Local().Format("2006-01-02 15:04:05"), timestamp))

	version, err := host.KernelVersion()
	if err != nil {
		log.Errorln("kernelVersion", err.Error())
		return err
	}
	t.Files = append(t.Files, fmt.Sprintf("HOST VERSION %s", version))

	platform, family, version, err := host.PlatformInformation()
	if err != nil {
		log.Errorln("PlatformInformation", err.Error())
		return err
	}
	t.Files = append(t.Files, fmt.Sprintf("HOST PLATFORM %s", platform))
	t.Files = append(t.Files, fmt.Sprintf("HOST FAMILY %s", family))
	t.Files = append(t.Files, fmt.Sprintf("HOST VERSION %s", version))

	users, err := host.Users()
	if err != nil {
		log.Errorln("host.Users", err.Error())
		return err
	}

	// 终端用户
	for i, u := range users {
		t.Files = append(t.Files, fmt.Sprintf("USER TERMINIAL %d %v", i, u))
	}

	// 内存
	swapMemory, err := mem.SwapMemory()
	if err != nil {
		log.Errorln("swapMemory", err.Error())
		return err
	}
	t.Files = append(t.Files, fmt.Sprintf("MEM SWAP %v", swapMemory))

	// 进程
	processes, err := process.Processes()
	if err != nil {
		log.Errorln("process0", err.Error())
		return err
	}

	var rootProcess *process.Process
	for _, p := range processes {
		if p.Pid == 0 {
			rootProcess = p
			break
		}
	}

	if rootProcess != nil {
		rootP, err := rootProcess.Children()
		if err != nil {
			log.Errorln("process", err.Error())
			return err
		}
		for i, pp := range rootP {
			t.Files = append(t.Files, fmt.Sprintf("PROCESS %d %v", i, pp))
		}
	}

	// services, _ := winservices.ListServices()

	// for _, service := range services {
	// 	newservice, _ := winservices.NewService(service.Name)
	// 	newservice.GetServiceDetail()
	// 	fmt.Println("Name:", newservice.Name, "Binary Path:", newservice.Config.BinaryPathName, "State: ", newservice.Status.State)
	// }

	cmds, err := ExecCommandString("compgen -c")
	if err != nil {
		log.Errorln("compgen", err.Error())
		return err
	}
	for _, c := range strings.Split(cmds, "\n") {
		t.Files = append(t.Files, fmt.Sprintf("CMD %s", c))
	}

	// 历史命令
	hcmds, err := ExecCommandString("cat ~/.zsh_history")
	if err != nil {
		log.Errorln("fc -rl 1", err.Error())
		return err
	}
	// log.Error("hcmds", hcmds)
	for _, c := range strings.Split(hcmds, "\n") {
		if strings.Contains(c, ";") {
			tp := strings.Split(c, ";")
			t.Files = append(t.Files, fmt.Sprintf("HISTORY %s", strings.Join(tp[1:], ";")))
		} else {
			t.Files = append(t.Files, fmt.Sprintf("HISTORY %s", c))
		}
	}
	return nil
}

// 递归获取当前目录所有文件
func (t *TuiScreen) GetAllFiles(dirPth string, istop bool) error {
	var dirs []string
	// fmt.Printf("dirPath %s\n", dirPth)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			if !In(fi.Name(), t.Unsearch) {
				dirs = append(dirs, fi.Name())
				// GetAllFiles(dirPth+PthSep+fi.Name(), suff)
				// 获取文件夹
				if istop {
					t.Files = append(t.Files, fmt.Sprintf("Dir %s", fi.Name()))
				} else {
					t.Files = append(t.Files, fmt.Sprintf("Dir %s%s%s", strings.Replace(dirPth, "./", "", 1), PthSep, fi.Name()))
				}

				t.SearchResult = t.Files
				t.Total = len(t.Files)
				// t.ShowResultFiles(true)
			}
		} else {
			// 过滤指定格式
			// for _, x := range suff {
			//      ok := strings.HasSuffix(fi.Name(), x)
			//      if ok {
			//              files = append(files, fi.Name())
			//      }
			// }
			if istop {
				t.Files = append(t.Files, fmt.Sprintf("File %s", fi.Name()))
			} else {
				t.Files = append(t.Files, fmt.Sprintf("File %s%s%s", strings.Replace(dirPth, "./", "", 1), PthSep, fi.Name()))
			}
			t.SearchResult = t.Files
			t.Total = len(t.Files)
			// t.ShowResultFiles(true)
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		// fmt.Printf("table is %s%s%s\n", dirPth, PthSep, table)
		err = t.GetAllFiles(fmt.Sprintf("%s%s%s", dirPth, PthSep, table), false)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		// for _, temp1 := range temp {
		// 	t.Files = append(t.Files, fmt.Sprintf("%s%s%s", table, PthSep, temp1))
		// }
	}

	return nil
}

func (t *TuiScreen) ShowAllFiles() {
	// t.screen.SetContent(1, 14, ' ', []rune("\x1b[1000D"), tcell.StyleDefault)
	// t.screen.SetContent(1, 15, ' ', []rune("\x1b[13A"), tcell.StyleDefault)
	if len(t.Files) > 0 {
		index := 0
		t.SearchNum = 0
		for _, x := range t.Files {
			// t.screen.SetContent(1, index+3, tcell.RuneHLine, []rune(x), tcell.StyleDefault)
			if index == 0 {
				t.printString(1, index+3, x, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSlateGrey), true)
			} else {
				if t.MaxContent > index {
					t.printString(1, index+3, x, tcell.StyleDefault, false)
				}
			}

			index++
			t.SearchNum = t.SearchNum + 1
		}
		t.ShowCount()
	}
}

// 移动历史记录
func (t *TuiScreen) ShowResultFiles(up bool) {
	if len(t.Input) == 0 && t.Total > 0 {
		t.SearchResult = t.Files
	}
	// t.Clear()
	t.screen.Clear()
	// 判断当前关闭范围
	// if t.CursorPos <= t.MaxContent-1 {
	// 	// 如果搜索结果小于最大展示数
	// 	if len(t.SearchResult) >= t.MaxContent {
	// 		for i, x := range t.SearchResult {
	// 			if i == t.CursorPos && i < t.MaxContent {
	// 				t.printString(1, i+3, x, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSlateGrey).Bold(true), true)
	// 				t.offset = i
	// 			} else if i < t.MaxContent && i != t.CursorPos {
	// 				t.printString(1, i+3, x, tcell.StyleDefault, false)
	// 			} else {
	// 				break
	// 			}
	// 		}
	// 	} else {
	// 		for i, x := range t.SearchResult {
	// 			if i == t.CursorPos && i < len(t.SearchResult) {
	// 				t.printString(1, i+3, x, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSlateGrey).Bold(true), true)
	// 				t.offset = i
	// 			} else if i < len(t.SearchResult) && i != t.CursorPos {
	// 				t.printString(1, i+3, x, tcell.StyleDefault, false)
	// 			} else {
	// 				break
	// 			}
	// 		}
	// 	}

	// } else {
	// t.Clear()
	// t.screen.Sync()
	// 判断光标偏移量
	// var more int
	if t.offset > 0 && up {
		t.offset = t.offset - 1
		// more = t.MaxContent - t.offset - 1
	} else if t.offset >= 0 && t.offset < t.MaxContent-1 && !up {
		t.offset = t.offset + 1
		// more = t.MaxContent - t.offset - 1
	}
	t.ShowCount()
	for i, x := range t.SearchResult[t.CursorPos-t.offset:] {
		if len(t.SearchResult[t.CursorPos-t.offset:]) > t.MaxContent {

			// TODO: 任意位置光标中间挪动
			if i == t.offset {
				t.printString(1, i+3, x, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSlateGrey).Bold(true), true)
			} else if i < t.MaxContent {
				t.printString(1, i+3, x, tcell.StyleDefault, false)
			} else {
				break
			}

		} else {
			if i == t.offset {
				t.printString(1, i+3, x, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSlateGrey).Bold(true), true)
			} else {
				t.printString(1, i+3, x, tcell.StyleDefault, false)
			}
		}
	}
	// }
	t.screen.Sync()

	t.screen.SetCell(2, 1, tcell.StyleDefault, t.Input...)
	// t.screen.SetCursorStyle(tcell.CursorStyleBlinkingUnderline)
	t.SetCursor(runewidth.StringWidth(string(t.Input)), 1, "x")
	t.ShowCount()
}

func (t *TuiScreen) ShowFiles() {
	search := string(t.Input)
	if len(t.Files) > 0 {
		index := 0
		// 清空结果缓存
		t.SearchNum = 0
		t.SearchResult = []string{}
		for _, x := range t.Files {
			// fmt.Println(x)
			// tmp := []rune(fmt.Sprintf("|%s", x))
			// t.screen.SetCell(1, index+2, tcell.StyleDefault, tmp...)
			if strings.Contains(x, search) {
				t.SearchResult = append(t.SearchResult, x)
				if t.CursorPos == index {
					t.printString(1, index+3, x, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSlateGrey), true)
				} else {
					if t.MaxContent > index {
						t.printString(1, index+3, x, tcell.StyleDefault, false)
					}
				}

				index++
				t.SearchNum = t.SearchNum + 1
				t.ShowCount()
			}

		}
	}
}

func (t *TuiScreen) FuzzySearch() {
	search := string(t.Input)
	if len(t.Files) > 0 {
		t.SearchResult = fuzzy.Find(search, t.Files)
		t.SearchNum = len(t.SearchResult)
		t.CursorPos = 0
		t.offset = 0
		t.ShowResultFiles(true)
	}
}

func (t *TuiScreen) ShowCount() {
	t.screen.SetContent(1, 2, tcell.RuneBoard, []rune(fmt.Sprintf(" %d/%d (cursprPos %d) MaxContent %d offset %d more %d index %d [%d:]", t.SearchNum, t.Total, t.CursorPos, t.MaxContent, t.offset, t.MaxContent-t.offset-1, t.CursorPos-t.MaxContent+1, t.CursorPos-t.MaxContent+1+t.MaxContent-t.offset-1)), tcell.StyleDefault)
	// t.printString(1, 2, fmt.Sprintf(" %d/%d", t.SearchNum, t.Total), tcell.StyleDefault.Foreground(tcell.ColorYellowGreen).Background(tcell.ColorDarkSlateGrey), false)
}

// 打印输出
func (t *TuiScreen) printString(x, y int, text string, style tcell.Style, isFouce bool) {
	if isFouce {
		t.CurContent = text
	}
	lx := 0
	lastX := x
	lastY := y
	position := 0

	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		rs := gr.Runes()

		if len(rs) == 1 {
			r := rs[0]
			if r < rune(' ') { // ignore control characters
				continue
			} else if r == '\n' {
				lastY++
				lx = 0
				continue
			} else if r == '\u000D' { // skip carriage return
				continue
			}
		}
		// TODO: 锁定范围
		// var xPos = w.left + w.lastX + lx
		// var yPos = w.top + w.lastY
		// if xPos < (w.left+w.width) && yPos < (w.top+w.height) {
		//      _screen.SetContent(xPos, yPos, rs[0], rs[1:], style)
		// }
		// lx += runewidth.StringWidth(string(rs))

		var xPos = lastX + lx
		var yPos = lastY
		// 定义第一个字符
		if isFouce && position == 0 {
			t.screen.SetCell(xPos, yPos, tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorDarkSlateGrey).Bold(true), '>')
			t.screen.SetCell(xPos+1, yPos, tcell.StyleDefault.Background(tcell.ColorDarkSlateGrey), ' ')
		} else if !isFouce && position == 0 {
			t.screen.SetCell(xPos, yPos, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSlateGrey), ' ')
			t.screen.SetCell(xPos+1, yPos, tcell.StyleDefault, ' ')
		}
		// 打印内容
		// 打印目标颜色 -> 逐字打印
		if len(t.Input) == 0 {
			t.screen.SetContent(xPos+2, yPos, rs[0], rs[1:], style)
		} else {
			// 有顺序匹配
			for ii, xx := range rs {
				// 包含部分加颜色
				if strings.Contains(string(t.Input), string(xx)) {
					if isFouce {
						t.screen.SetCell(xPos+2+ii, yPos, tcell.StyleDefault.Foreground(tcell.ColorDarkOliveGreen).Background(tcell.ColorDarkSlateGrey), xx)
					} else {
						t.screen.SetCell(xPos+2+ii, yPos, tcell.StyleDefault.Foreground(tcell.ColorDarkOliveGreen), xx)
					}
				} else {
					if isFouce {
						t.screen.SetCell(xPos+2+ii, yPos, tcell.StyleDefault.Background(tcell.ColorDarkSlateGrey), xx)
					} else {
						t.screen.SetCell(xPos+2+ii, yPos, tcell.StyleDefault, xx)
					}
				}
			}
		}

		lx += runewidth.StringWidth(string(rs))
		position++
	}
	lastX += lx
}

func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

func (t *TuiScreen) Clear() {
	t.screen.Sync()
	t.screen.Clear()
}

func (t *TuiScreen) Show() {
	t.screen.Show()
}

func (t *TuiScreen) Sync() {
	t.screen.Sync()
}

func (t *TuiScreen) SetCursor(x, y int, value string) {
	t.screen.SetCell(1, y, tcell.StyleDefault, []rune{'>', ' '}...)
	// t.screen.SetCursorStyle(tcell.CursorStyleDefault)
	t.screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
	t.screen.ShowCursor(x+3, 1)
}

func (t *TuiScreen) Draw(x1, y1, x2, y2 int, text string) {
	drawBox(t.screen, x1, y1, x2, y2, t.boxStyle, text)
}

func (t *TuiScreen) Event() {
	// fmt.Printf("\x1b[1000D")
	// fmt.Printf("\x1b[%dA", int(t.height/2))
	t.screen.Show()
	style := tcell.StyleDefault

	// Poll event
	ev := t.screen.PollEvent()

	// Process event
	switch ev := ev.(type) {
	case *tcell.EventResize:
		t.Sync()
		t.Clear()
		t.width, t.height = t.screen.Size()
		t.MaxContent = t.height - 4
		t.ShowResultFiles(true)
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			quit(t.screen)
		case tcell.KeyCtrlC:
			quit(t.screen)
		case tcell.KeyHome:
			t.CursorPos = 0
			t.ShowResultFiles(true)
		case tcell.KeyEnd:
			if len(t.SearchResult) == 0 {
				t.SearchResult = t.Files
			}
			t.CursorPos = len(t.SearchResult)
			t.ShowResultFiles(false)
		case tcell.KeyUp:
			if t.CursorPos > 0 {
				t.CursorPos = t.CursorPos - 1
				t.ShowResultFiles(true)
			}
		case tcell.KeyEnter:
			tmp := strings.Split(t.CurContent, " ")
			// 打印第二段
			if tmp[0] == "CMD" {
				err := CommandPty(strings.Join(tmp[1:], " "))
				if err != nil {
					panic(err)
				}
				t.Sync()
				t.Clear()
				t.ShowResultFiles(true)
			} else if tmp[0] == "HISTORY" {
				err := CommandPty(strings.Join(tmp[1:], " "))
				if err != nil {
					panic(err)
				}
				t.Sync()
				t.Clear()
				t.ShowResultFiles(true)
			} else if tmp[0] == "File" {
				err := CommandPty("vi " + strings.Join(tmp[1:], " "))
				if err != nil {
					panic(err)
				}
				t.Sync()
				t.Clear()
				t.ShowResultFiles(true)
			} else if tmp[0] == "Dir" {
				err := CommandPty("cd " + strings.Join(tmp[1:], " "))
				if err != nil {
					panic(err)
				}
				t.Sync()
				t.Clear()
				t.ShowResultFiles(true)
			} else {
				fmt.Fprintf(os.Stdout, "%s", strings.Join(tmp[1:], " "))
				os.Exit(0)
			}
		case tcell.KeyDown:
			if t.SearchNum-1 > t.CursorPos {
				t.CursorPos = t.CursorPos + 1
				t.ShowResultFiles(false)
			}
		case tcell.KeyCtrlW:
			t.Input = []rune{}
			t.ShowResultFiles(true)
		case tcell.KeyCtrlA:
			t.SetCursor(1, 1, "x")
			t.ShowResultFiles(true)
		case tcell.KeyCtrlL:
			t.Sync()
		case tcell.KeyCtrlK:
			// 向上移动
			if t.CursorPos > 0 {
				t.CursorPos = t.CursorPos - 1
				t.ShowResultFiles(true)
			}
		case tcell.KeyCtrlJ:
			if t.SearchNum-1 > t.CursorPos {
				t.CursorPos = t.CursorPos + 1
				t.ShowResultFiles(false)
			}
		case tcell.KeyBackspace2:
			// if len(t.Input) > 0 {
			//      // fmt.Println(t.Input)
			//      t.Input = t.Input[:len(t.Input)-1]
			//      // fmt.Println(t.Input)
			// }

			if len(t.Input) > 0 {
				t.Clear()
				t.Input = t.Input[:len(t.Input)-1]
				t.screen.SetCell(2, 1, tcell.StyleDefault, t.Input...)
				t.screen.SetCursorStyle(tcell.CursorStyleBlinkingUnderline)
				t.SetCursor(runewidth.StringWidth(string(t.Input)), 1, "x")
				// t.ShowFiles()
				t.FuzzySearch()
			}

		case tcell.KeyRune:
			// 查询重置搜索光标
			t.CursorPos = 0
			value := ev.Rune()
			switch value {
			case '0':
				t.Clear()
				t.Input = append(t.Input, value)
				t.SetCursor(1, 1, "x")
				t.screen.SetContent(2, 1, '0', nil, style)
				t.screen.SetCursorStyle(tcell.CursorStyleDefault)
			case '1':
				t.Clear()
				t.Input = append(t.Input, value)
				t.SetCursor(1, 1, "x")
				t.screen.SetContent(2, 1, '1', nil, style)
				t.screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
			case '2':
				t.Input = append(t.Input, value)
				t.screen.SetCell(2, 1, tcell.StyleDefault, t.Input...)
				t.screen.SetCursorStyle(tcell.CursorStyleSteadyBlock)
			case '3':
				t.Input = append(t.Input, value)
				t.screen.SetCell(2, 1, tcell.StyleDefault, t.Input...)
				t.screen.SetCursorStyle(tcell.CursorStyleBlinkingUnderline)
			case '4':
				t.Input = append(t.Input, value)
				t.screen.SetCell(2, 1, tcell.StyleDefault, t.Input...)
				t.screen.SetCursorStyle(tcell.CursorStyleSteadyUnderline)
			case '5':
				t.Input = append(t.Input, value)
				t.screen.SetCell(2, 1, tcell.StyleDefault, t.Input...)
				t.screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)
			case '6':
				t.Input = append(t.Input, value)
				t.screen.SetCell(2, 1, tcell.StyleDefault, t.Input...)
				t.screen.SetCursorStyle(tcell.CursorStyleSteadyBar)
			// case 'C':
			// 	t.Clear()
			// case 'c':
			//      t.Clear()
			default:
				t.Clear()
				t.Input = append(t.Input, value)
				t.screen.SetCell(2, 1, tcell.StyleDefault, t.Input...)
				t.screen.SetCursorStyle(tcell.CursorStyleSteadyBar)
				t.SetCursor(runewidth.StringWidth(string(t.Input)), 1, "x")
				// t.ShowFiles()
				t.FuzzySearch()
			}

			t.screen.Show()
		}
	case *tcell.EventMouse:
		x, y := ev.Position()
		button := ev.Buttons()
		// Only process button events, not wheel events
		button &= tcell.ButtonMask(0xff)

		if button != tcell.ButtonNone && t.X < 0 {
			t.X, t.Y = x, y

		}
		switch ev.Buttons() {
		case tcell.ButtonNone:
			if t.X >= 0 {
				label := fmt.Sprintf("%d,%d to %d,%d", t.X, t.Y, x, y)
				drawBox(t.screen, t.X, t.Y, x, y, t.boxStyle, label)
				// dishelloworld(s)
				t.X, t.Y = -1, -1

			}

		}

	}
}

func quit(s tcell.Screen) {
	s.Fini()
	os.Exit(0)

}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1

		}
		if row > y2 {
			break

		}

	}

}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func dishelloworld(s tcell.Screen) {
	w, h := s.Size()
	s.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorBlue.TrueColor()).Background(tcell.ColorWhite)
	emitStr(s, w/2-7, h/4, style, "Command Search")
	emitStr(s, w/2-7, h/4+1, tcell.StyleDefault, "Enter Search")
	s.Show()
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1

	}
	if x2 < x1 {
		x1, x2 = x2, x1

	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)

		}

	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)

	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)

	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)

	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)

}

func Run(unsearch []string, show bool) {
	var err error
	x := &TuiScreen{
		EnablePaste: true,
		EnableMouse: false,
		Unsearch:    unsearch,
		finish:      make(chan int),
		Files:       []string{},
	}

	x.init()

	go func() {
		x.GetCommand()
		err = x.Gopsutil()
		if err != nil {
			panic(err)
		}
		err = x.GetAllFiles(".", true)
		if err != nil {
			panic(err)
		}
		x.finish <- 1
	}()

	var s *spinner.Spinner
	if show {
		s = spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
		s.Start()
		s.Prefix = ""
		s.Suffix = " 扫描中..." // Start the spinner
		s.UpdateCharSet(spinner.CharSets[39])
		s.Reverse()
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			if show {
				s.Suffix = fmt.Sprintf("扫描中: %d MaxContent %d @%s", x.Total, x.MaxContent, time.Now().Format("2006-01-02 15:04:05"))
				s.Reverse()
				s.Restart()
			} else {
				x.ShowResultFiles(true)
			}

		case <-x.finish:
			goto Loop
		}

	}
Loop:
	if show {
		s.Stop()
	}

	ticker.Stop()

	// CommandPty("\x1b[10A")
	x.SetCursor(1, 1, "x")

	x.ShowAllFiles()

	for {
		x.Event()
	} // Run for some time to simulate work

}
