package addcolumns

import (
	"encoding/json"
	"strings"
)

type Columns struct {
	Field     string //state
	Title     string
	Checkbox  bool
	Rowspan   int
	Colspan   int
	Align     string
	Valign    string
	Sortable  bool
	editable  bool
	Formatter string
}

func NewColumns(filed string, checkbox bool) *Columns {
	return &Columns{
		Field:    filed,
		Title:    strings.ToUpper(filed),
		Checkbox: checkbox,
		Align:    "center",
		Valign:   "middle",
		Sortable: true,
	}
}

func MutilColumms(data ...string) []Columns {
	tmp := []Columns{Columns{Checkbox: true}}
	for _, x := range data {
		tmp = append(tmp, *NewColumns(x, false))
	}
	return tmp
}

func DirectJson(data ...string) (string, error) {
	tmp := MutilColumms(data...)
	jsons, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}
	return strings.ToLower(string(jsons)), nil
}
