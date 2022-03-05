package cmd

var (
	extended bool   // -x --extended
	exact    bool   // -e --exact
	query    string // -q --query
	filter   string // -f --filter
	literal  bool   // --literal
	algo     string
	// TODO: FuzzyAlgo
	// TODO: Expect
	phony bool // --enabled --phony
	// TODO: Criteria
	// TODO: --bind parseKeymap
	bind string
	// TODO: Theme
	// TODO: --toggle-sort
	// TODO: Delimiter -d --delimiter
	delimiter string
	// TODO: Nth
	// TODO: WithNth
	withnth string
	nth     string
	sort    int // -s --sort 0
	nosort  bool
	tac     bool // --tac
	// TODO: Case
	multi       int  // -m --multi
	ansi        bool // --ansi
	mouse       bool // --no-mouse
	black       bool // --black
	bold        bool // --bold
	insensitive bool
	enabled     bool
	disabled    bool
	tiebreak    string
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
	layout      string
	border      string
	color       string
)
