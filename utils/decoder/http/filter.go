package http

import "regexp"

// url: /tet & method: get & body: balabal & Content-Type: application/json
type Filter struct {
	filterStr string
	filters   map[string]*regexp.Regexp
}

func (f Filter) String() string {
	return f.filterStr
}

func (f *Filter) IsEmpty() bool {
	return len(f.filters) == 0
}

func (f *Filter) compile() {
	if len(f.filterStr) == 0 {
		return
	}
	pattern := regexp.MustCompile("\\s+&\\s+")
	filters := pattern.Split(f.filterStr, -1)
	subPattern := regexp.MustCompile("\\s*:\\s*")
	for _, filter := range filters {
		result := subPattern.Split(filter, 2)
		f.filters[result[0]] = regexp.MustCompile("(?im:" + result[1] + ")")
	}
}

func NewFilter(filterStr string) *Filter {
	f := &Filter{filterStr, map[string]*regexp.Regexp{}}
	f.compile()
	return f
}
