package fzf

import (
	"os"

	"github.com/mattn/go-shellwords"
)

// ParseOptions parses Cobra command-line options
func ParseOptionsCobra() *Options {
	opts := defaultOptions()

	// Options from Env var
	words, _ := shellwords.Parse(os.Getenv("FZF_DEFAULT_OPTS"))
	if len(words) > 0 {
		parseOptions(opts, words)
	}

	// Options from command-line arguments
	// parseOptions(opts, os.Args[1:])

	// postProcessOptions(opts)
	return opts
}

// Step Two is Cobra Process
// Step Three
func PostProcessOptions(opts *Options) {
	postProcessOptions(opts)
}
