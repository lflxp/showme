package cmd

import (
	"fmt"
	"math"
	"os"
	"time"

	fzf "github.com/lflxp/fzf/src"
)

const (
	// Core
	coordinatorDelayMax  time.Duration = 100 * time.Millisecond
	coordinatorDelayStep time.Duration = 10 * time.Millisecond

	// Reader
	readerBufferSize       = 64 * 1024
	readerPollIntervalMin  = 10 * time.Millisecond
	readerPollIntervalStep = 5 * time.Millisecond
	readerPollIntervalMax  = 50 * time.Millisecond

	// Terminal
	initialDelay      = 20 * time.Millisecond
	initialDelayTac   = 100 * time.Millisecond
	spinnerDuration   = 100 * time.Millisecond
	previewCancelWait = 500 * time.Millisecond
	previewChunkDelay = 100 * time.Millisecond
	previewDelayed    = 500 * time.Millisecond
	maxPatternLength  = 300
	maxMulti          = math.MaxInt32

	// Matcher
	numPartitionsMultiplier = 8
	maxPartitions           = 32
	progressMinDuration     = 200 * time.Millisecond

	// Capacity of each chunk
	chunkSize int = 100

	// Pre-allocated memory slices to minimize GC
	slab16Size int = 100 * 1024 // 200KB * 32 = 12.8MB
	slab32Size int = 2048       // 8KB * 32 = 256KB

	// Do not cache results of low selectivity queries
	queryCacheMax int = chunkSize / 5

	// Not to cache mergers with large lists
	mergerCacheMax int = 100000

	// History
	defaultHistoryMax int = 1000

	// Jump labels
	defaultJumpLabels string = "asdfghjklqwertyuiopzxcvbnm1234567890ASDFGHJKLQWERTYUIOPZXCVBNM`~;:,<.>/?'\"!@#$%^&*()[{]}-_=+"

	exitCancel    = -1
	exitOk        = 0
	exitNoMatch   = 1
	exitError     = 2
	exitInterrupt = 130
)

var (
	extended bool   // -x --extended
	exact    bool   // -e --exact
	query    string // -q --query
	filter   string // -f --filter
	literal  bool   // --literal
	// TODO: FuzzyAlgo
	// TODO: Expect
	phony bool // --enabled --phony
	// TODO: Criteria
	// TODO: --bind parseKeymap
	// TODO: Theme
	// TODO: --toggle-sort
	// TODO: Delimiter -d --delimiter
	// TODO: Nth
	// TODO: WithNth
	sort int  // -s --sort 0
	tac  bool // --tac
	// TODO: Case
	multi int  // -m --multi
	ansi  bool // --ansi
	mouse bool // --no-mouse
	black bool // --black
	bold  bool // --bold
	// TODO: Layout
	cycle      bool // --cycle
	keepright  bool // --keep-right
	hscroll    bool // --hscroll
	hscrolloff int  // --hscroll-off
	scrolloff  int  // --scroll-off
	fileword   bool // --filepath-word
	// TODO: InfoStyle --info
	jumplabels string // --jump-labels
	select1    bool   // -1 --select-1
	exit0      bool   // -0 --exit-0
	readzero   bool   // --read0
	// TODO: --print0
	print0     bool
	printquery bool   // --print-query
	prompt     string // --prompt
	pointer    string // --pointer
	marker     string // --marker
	sync       bool   // --sync
	// TODO: History
	// TODO: Header
	headerlines int    // --no-header-lines
	headerfirst bool   // --no-header-first
	preview     string // --preview
	// TODO: --preview-window
	// TODO: Height
	height    string // --height=40%
	minheight int    // --min-height
	// TODO: margin
	// TODO: padding
	// TODO: BorderShape
	unicode     bool // --unicode
	tabstop     int  // --tabstop
	clearonexit bool // --clear
	versions    bool // version
)

func errorExit(msg string) {
	os.Stderr.WriteString(msg + "\n")
	os.Exit(exitError)
}

func parseFzfArgs(opts *fzf.Options) {
	// var historyMax int
	// if opts.History == nil {
	// 	historyMax = defaultHistoryMax
	// } else {
	// 	historyMax = opts.History.GetSize()
	// }
	// setHistory := func(path string) {
	// 	h, e := fzf.NewHistory(path, historyMax)
	// 	if e != nil {
	// 		errorExit(e.Error())
	// 	}
	// 	opts.History = h
	// }
	// setHistoryMax := func(max int) {
	// 	historyMax = max
	// 	if historyMax < 1 {
	// 		errorExit("history max must be a positive integer")
	// 	}
	// 	if opts.History != nil {
	// 		opts.History.SetSize(historyMax)
	// 	}
	// }
	validateJumpLabels := false
	// validatePointer := false
	// validateMarker := false
	opts.Extended = extended
	opts.Fuzzy = exact
	opts.Query = query
	opts.Filter = &filter
	opts.Normalize = literal
	// case "--algo":
	// opts.FuzzyAlgo = parseAlgo(nextString(allArgs, &i, "algorithm required (v1|v2)"))
	// case "--expect":
	// for k, v := range parseKeyChords(nextString(allArgs, &i, "key names required"), "key names required") {
	// opts.Expect[k] = v
	// }
	// case "--no-expect":
	// opts.Expect = make(map[tui.Event]string)
	opts.Phony = phony
	// case "--tiebreak":
	// opts.Criteria = parseTiebreak(nextString(allArgs, &i, "sort criterion required"))
	// case "--bind":
	// parseKeymap(opts.Keymap, nextString(allArgs, &i, "bind expression required"))
	// case "--color":
	// 	_, spec := optionalNextString(allArgs, &i)
	// 	if len(spec) == 0 {
	// 		opts.Theme = tui.EmptyTheme()
	// 	} else {
	// 		opts.Theme = parseTheme(opts.Theme, spec)
	// 	}
	// case "--toggle-sort":
	// 	parseToggleSort(opts.Keymap, nextString(allArgs, &i, "key name required"))
	// case "-d", "--delimiter":
	// 	opts.Delimiter = delimiterRegexp(nextString(allArgs, &i, "delimiter required"))
	// case "-n", "--nth":
	// 	opts.Nth = splitNth(nextString(allArgs, &i, "nth expression required"))
	// case "--with-nth":
	// 	opts.WithNth = splitNth(nextString(allArgs, &i, "nth expression required"))
	// case "-s", "--sort":
	// 	opts.Sort = optionalNumeric(allArgs, &i, 1)
	// case "+s", "--no-sort":
	opts.Sort = sort
	opts.Tac = tac
	// case "-i":
	// 	opts.Case = CaseIgnore
	// case "+i":
	// 	opts.Case = CaseRespect
	// case "-m", "--multi":
	// 	opts.Multi = optionalNumeric(allArgs, &i, maxMulti)
	// case "+m", "--no-multi":
	opts.Multi = multi
	opts.Ansi = ansi
	opts.Mouse = mouse
	// case "+c", "--no-color":
	// 	opts.Theme = tui.NoColorTheme()
	// case "+2", "--no-256":
	// 	opts.Theme = tui.Default16
	opts.Black = black
	opts.Bold = bold
	// case "--layout":
	// 	opts.Layout = parseLayout(
	// 		nextString(allArgs, &i, "layout required (default / reverse / reverse-list)"))
	// case "--reverse":
	// 	opts.Layout = layoutReverse
	// case "--no-reverse":
	// 	opts.Layout = layoutDefault
	opts.Cycle = cycle
	opts.KeepRight = keepright
	opts.Hscroll = hscroll
	opts.HscrollOff = hscrolloff
	opts.ScrollOff = scrolloff
	opts.FileWord = fileword
	// case "--info":
	// 	opts.InfoStyle = parseInfoStyle(
	// 		nextString(allArgs, &i, "info style required"))
	// case "--no-info":
	// 	opts.InfoStyle = infoHidden
	// case "--inline-info":
	// 	opts.InfoStyle = infoInline
	// case "--no-inline-info":
	// 	opts.InfoStyle = infoDefault
	opts.JumpLabels = jumplabels
	// if len(jumplabels) > 0 {
	// 	validateJumpLabels = true
	// }
	opts.Select1 = select1
	opts.Exit0 = exit0
	opts.ReadZero = readzero
	if print0 {
		opts.Printer = func(str string) { fmt.Print(str, "\x00") }
		opts.PrintSep = "\x00"

	} else {

		opts.Printer = func(str string) { fmt.Println(str) }
		opts.PrintSep = "\n"
	}
	opts.PrintQuery = printquery
	opts.Prompt = prompt
	opts.Pointer = pointer
	// validatePointer = true
	opts.Marker = marker
	// validateMarker = true
	opts.Sync = sync
	// case "--no-history":
	// 	opts.History = nil
	// case "--history":
	// 	setHistory(nextString(allArgs, &i, "history file path required"))
	// case "--history-size":
	// 	setHistoryMax(nextInt(allArgs, &i, "history max size required"))
	// case "--no-header":
	// 	opts.Header = []string{}
	opts.HeaderLines = headerlines
	// case "--header":
	// 	opts.Header = strLines(nextString(allArgs, &i, "header string required"))
	// case "--header-lines":
	// 	opts.HeaderLines = atoi(
	// 		nextString(allArgs, &i, "number of header lines required"))
	opts.HeaderFirst = headerfirst
	// opts.Preview.command = preview
	// case "--preview-window":
	// 	parsePreviewWindow(&opts.Preview,
	// 		nextString(allArgs, &i, "preview window layout required: [up|down|left|right][,SIZE[%]][,border-BORDER_OPT][,wrap][,cycle][,hidden][,+SCROLL[OFFSETS][/DENOM]][,~HEADER_LINES][,default]"))
	// case "--height":
	// 	opts.Height = parseHeight(nextString(allArgs, &i, "height required: HEIGHT[%]"))
	opts.MinHeight = minheight
	// case "--no-height":
	// 	opts.Height = sizeSpec{}
	// case "--no-margin":
	// 	opts.Margin = defaultMargin()
	// case "--no-padding":
	// 	opts.Padding = defaultMargin()
	// case "--no-border":
	// 	opts.BorderShape = tui.BorderNone
	// case "--border":
	// 	hasArg, arg := optionalNextString(allArgs, &i)
	// // 	opts.BorderShape = parseBorder(arg, !hasArg)
	// case "--no-unicode":
	opts.Unicode = unicode
	// case "--margin":
	// 	opts.Margin = parseMargin(
	// 		"margin",
	// 		nextString(allArgs, &i, "margin required (TRBL / TB,RL / T,RL,B / T,R,B,L)"))
	// case "--padding":
	// 	opts.Padding = parseMargin(
	// 		"padding",
	// 		nextString(allArgs, &i, "padding required (TRBL / TB,RL / T,RL,B / T,R,B,L)"))
	opts.Tabstop = tabstop
	opts.ClearOnExit = clearonexit
	opts.Version = versions
	// default:
	// 	if match, value := optString(arg, "--algo="); match {
	// 		opts.FuzzyAlgo = parseAlgo(value)
	// 	} else if match, value := optString(arg, "-q", "--query="); match {
	// 		opts.Query = value
	// 	} else if match, value := optString(arg, "-f", "--filter="); match {
	// 		opts.Filter = &value
	// 	} else if match, value := optString(arg, "-d", "--delimiter="); match {
	// 		opts.Delimiter = delimiterRegexp(value)
	// 	} else if match, value := optString(arg, "--border="); match {
	// 		opts.BorderShape = parseBorder(value, false)
	// 	} else if match, value := optString(arg, "--prompt="); match {
	// 		opts.Prompt = value
	// 	} else if match, value := optString(arg, "--pointer="); match {
	// 		opts.Pointer = value
	// 		validatePointer = true
	// 	} else if match, value := optString(arg, "--marker="); match {
	// 		opts.Marker = value
	// 		validateMarker = true
	// 	} else if match, value := optString(arg, "-n", "--nth="); match {
	// 		opts.Nth = splitNth(value)
	// 	} else if match, value := optString(arg, "--with-nth="); match {
	// 		opts.WithNth = splitNth(value)
	// 	} else if match, _ := optString(arg, "-s", "--sort="); match {
	// 		opts.Sort = 1 // Don't care
	// 	} else if match, value := optString(arg, "-m", "--multi="); match {
	// 		opts.Multi = atoi(value)
	// 	} else if match, value := optString(arg, "--height="); match {
	// 		opts.Height = parseHeight(value)
	// 	} else if match, value := optString(arg, "--min-height="); match {
	// 		opts.MinHeight = atoi(value)
	// 	} else if match, value := optString(arg, "--layout="); match {
	// 		opts.Layout = parseLayout(value)
	// 	} else if match, value := optString(arg, "--info="); match {
	// 		opts.InfoStyle = parseInfoStyle(value)
	// 	} else if match, value := optString(arg, "--toggle-sort="); match {
	// 		parseToggleSort(opts.Keymap, value)
	// 	} else if match, value := optString(arg, "--expect="); match {
	// 		for k, v := range parseKeyChords(value, "key names required") {
	// 			opts.Expect[k] = v
	// 		}
	// 	} else if match, value := optString(arg, "--tiebreak="); match {
	// 		opts.Criteria = parseTiebreak(value)
	// 	} else if match, value := optString(arg, "--color="); match {
	// 		opts.Theme = parseTheme(opts.Theme, value)
	// 	} else if match, value := optString(arg, "--bind="); match {
	// 		parseKeymap(opts.Keymap, value)
	// 	} else if match, value := optString(arg, "--history="); match {
	// 		setHistory(value)
	// 	} else if match, value := optString(arg, "--history-size="); match {
	// 		setHistoryMax(atoi(value))
	// 	} else if match, value := optString(arg, "--header="); match {
	// 		opts.Header = strLines(value)
	// 	} else if match, value := optString(arg, "--header-lines="); match {
	// 		opts.HeaderLines = atoi(value)
	// 	} else if match, value := optString(arg, "--preview="); match {
	// 		opts.Preview.command = value
	// 	} else if match, value := optString(arg, "--preview-window="); match {
	// 		parsePreviewWindow(&opts.Preview, value)
	// 	} else if match, value := optString(arg, "--margin="); match {
	// 		opts.Margin = parseMargin("margin", value)
	// 	} else if match, value := optString(arg, "--padding="); match {
	// 		opts.Padding = parseMargin("padding", value)
	// 	} else if match, value := optString(arg, "--tabstop="); match {
	// 		opts.Tabstop = atoi(value)
	// 	} else if match, value := optString(arg, "--hscroll-off="); match {
	// 		opts.HscrollOff = atoi(value)
	// 	} else if match, value := optString(arg, "--scroll-off="); match {
	// 		opts.ScrollOff = atoi(value)
	// 	} else if match, value := optString(arg, "--jump-labels="); match {
	// 		opts.JumpLabels = value
	// 		validateJumpLabels = true
	// 	} else {
	// 		errorExit("unknown option: " + arg)
	// 	}
	// }

	if opts.HeaderLines < 0 {
		errorExit("header lines must be a non-negative integer")
	}

	if opts.HscrollOff < 0 {
		errorExit("hscroll offset must be a non-negative integer")
	}

	if opts.ScrollOff < 0 {
		errorExit("scroll offset must be a non-negative integer")
	}

	if opts.Tabstop < 1 {
		errorExit("tab stop must be a positive integer")
	}

	if len(opts.JumpLabels) == 0 {
		errorExit("empty jump labels")
	}

	if validateJumpLabels {
		for _, r := range opts.JumpLabels {
			if r < 32 || r > 126 {
				errorExit("non-ascii jump labels are not allowed")
			}
		}
	}

	// if validatePointer {
	// 	if err := validateSign(opts.Pointer, "pointer"); err != nil {
	// 		errorExit(err.Error())
	// 	}
	// }

	// if validateMarker {
	// 	if err := validateSign(opts.Marker, "marker"); err != nil {
	// 		errorExit(err.Error())
	// 	}
	// }
}

// func validateSign(sign string, signOptName string) error {
// 	if sign == "" {
// 		return fmt.Errorf("%v cannot be empty", signOptName)
// 	}
// 	for _, r := range sign {
// 		if !unicode.IsGraphic(r) {
// 			return fmt.Errorf("invalid character in %v", signOptName)
// 		}
// 	}
// 	if runewidth.StringWidth(sign) > 2 {
// 		return fmt.Errorf("%v display width should be up to 2", signOptName)
// 	}
// 	return nil
// }
