/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	fzf "github.com/junegunn/fzf/src"
	"github.com/spf13/cobra"
)

var version string = "0.29.1"
var revision string = "cobra-dev"

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

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "fzf功能继承，实现当前文件快速搜索",
	Long:  `合并fzf功能，实现一个命令集成多种功能`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Args = os.Args[1:]
		fzf.Run(fzf.ParseOptions(), version, revision)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// =============================fzf========================
	searchCmd.Flags().BoolVarP(&extended, "extended", "x", false, "Extended-search mode (enabled by default; +x or --no-extended to disable)")
	searchCmd.Flags().BoolVarP(&exact, "exact", "e", false, "Enable Exact-match")
	searchCmd.Flags().StringVarP(&query, "query", "q", "", "Start the finder with the given query")
	searchCmd.Flags().StringVar(&algo, "algo", "v2", "Start the finder with the given query")
	searchCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter mode. Do not start interactive finder.")
	searchCmd.Flags().StringVarP(&nth, "nth", "n", "", "Comma-separated list of field index expressions")
	searchCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "AWK-style", "Field delimiter regex (default: AWK-style)")
	searchCmd.Flags().StringVar(&withnth, "with-nth", "", "Transform the presentation of each line using field index expressions")
	searchCmd.Flags().BoolVar(&literal, "literal", false, "Do not normalize latin script letters before matching")
	searchCmd.Flags().BoolVar(&nosort, "no-sort", false, "Do not sort the result")
	searchCmd.Flags().BoolVarP(&insensitive, "insensitive", "i", false, "Do not normalize latin script letters before matching")
	searchCmd.Flags().BoolVar(&phony, "phony", false, "Enable Exact-match")
	searchCmd.Flags().StringVar(&tiebreak, "tiebreak", "length", "omma-separated list of sort criteria to apply 	when the scores are tied [length|begin|end|index] 	(default: length)")
	searchCmd.Flags().BoolVar(&enabled, "enabled", false, "Do not perform search")
	searchCmd.Flags().BoolVar(&disabled, "disabled", false, "Do not perform search")
	searchCmd.Flags().IntVarP(&sort, "sort", "s", 0, "Do not sort the result")
	searchCmd.Flags().BoolVar(&tac, "tac", false, "Reverse the order of the input")
	searchCmd.Flags().IntVarP(&multi, "multi", "m", 0, "Enable multi-select with tab/shift-tab")
	searchCmd.Flags().BoolVar(&ansi, "ansi", false, "Enable processing of ANSI color codes")
	searchCmd.Flags().BoolVar(&mouse, "mouse", false, "Enable mouse")
	searchCmd.Flags().BoolVar(&mouse, "no-mouse", false, "Disabled mouse")
	searchCmd.Flags().BoolVar(&black, "black", false, "black")
	searchCmd.Flags().BoolVar(&bold, "bold", false, "bold")
	searchCmd.Flags().BoolVar(&cycle, "cycle", false, "cycle")
	searchCmd.Flags().BoolVar(&keepright, "keep-right", false, "keepright")
	searchCmd.Flags().BoolVar(&hscroll, "hscroll", false, "hscroll")
	searchCmd.Flags().BoolVar(&hscroll, "no-hscroll", false, "Disable horizontal scroll")
	searchCmd.Flags().IntVar(&hscrolloff, "hscroll-off", 0, "hscrolloff")
	searchCmd.Flags().IntVar(&scrolloff, "scroll-off", 0, "scrolloff")
	searchCmd.Flags().BoolVar(&fileword, "filepath-word", false, "Make word-wise movements respect path separator")
	searchCmd.Flags().StringVar(&jumplabels, "jump-labels", "", "Label characters for jump and jump-accept")
	searchCmd.Flags().BoolVarP(&select1, "select-1", "1", false, "Automatically select the only match")
	searchCmd.Flags().BoolVarP(&exit0, "exit-0", "0", false, "Exit immediately when there's no match")
	searchCmd.Flags().BoolVar(&readzero, "read0", false, "Read input delimited by ASCII NUL characters")
	searchCmd.Flags().BoolVar(&print0, "print0", false, "Print output delimited by ASCII NUL character")
	searchCmd.Flags().StringVar(&prompt, "expect", "", "Comma-separated list of keys to complete fzf")
	searchCmd.Flags().BoolVar(&printquery, "print-query", false, "Print query as the first line")
	searchCmd.Flags().StringVar(&pointer, "pointer", "", "Pointer to the current line (default: '>'")
	searchCmd.Flags().StringVar(&marker, "marker", "", "Multi-select marker (default: '>')")
	searchCmd.Flags().BoolVar(&sync, "sync", false, "Synchronous search for multi-staged filtering")
	searchCmd.Flags().StringVar(&preview, "preview", "", "Command to preview highlighted line ({})")
	searchCmd.Flags().StringVar(&preview, "preview-window", "right:50%", "Preview window layout (default: right:50%)")
	searchCmd.Flags().StringVar(&height, "height", "40%", "Display fzf window below the cursor with the give")
	searchCmd.Flags().IntVar(&minheight, "min-height", 0, "Minimum height when --height is given in percent default: 10")
	searchCmd.Flags().BoolVar(&unicode, "unicode", false, "unicode")
	searchCmd.Flags().IntVar(&tabstop, "tabstop", 8, "Number of spaces for a tab character (default: 8)")
	searchCmd.Flags().BoolVar(&clearonexit, "clear", false, "clearonexit")
	searchCmd.Flags().BoolVar(&versions, "versions", false, "versions")
	searchCmd.Flags().StringVar(&bind, "bind", "", "Custom key bindings. Refer to the man page.")
	searchCmd.Flags().StringVar(&layout, "layout", "default", "Choose layout: [default|reverse|reverse-list]")
	searchCmd.Flags().StringVar(&border, "border", "", "Draw border around the finder")
	searchCmd.Flags().StringVar(&border, "margin", "rounded", "Screen margin (TRBL | TB,RL | T,RL,B | T,R,B,L)")
	searchCmd.Flags().StringVar(&border, "padding", "", "Padding inside border (TRBL | TB,RL | T,RL,B | T,R,B,L)")
	searchCmd.Flags().StringVar(&border, "info", "", "Finder info style [default|inline|hidden]")
	searchCmd.Flags().StringVar(&border, "prompt", ">", "Input prompt (default: '> ')")
	searchCmd.Flags().StringVar(&border, "header", "", "String to print as header")
	searchCmd.Flags().IntVar(&headerlines, "header-lines", 0, "The first N lines of the input are treated as heade")
	searchCmd.Flags().BoolVar(&headerfirst, "header-first", false, "Print header before the prompt line")
	searchCmd.Flags().StringVar(&color, "color", "", "Base scheme (dark|light|16|bw) and/or custom colors")
	searchCmd.Flags().BoolVar(&versions, "no-bold", false, "Do not use bold text")
	searchCmd.Flags().StringVar(&color, "history", "", "History file")
	searchCmd.Flags().StringVar(&color, "history-size", "1000", "Maximum number of history entries (default: 1000)")
}
