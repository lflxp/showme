package str2html

import (
	"fmt"
	"html/template"
)

// string转html
type Str2html struct {
}

func (s *Str2html) Transfer() interface{} {
	return func(data string) template.HTML {
		if err := s.check(data); err != nil {
			return template.HTML(err.Error())
		}
		return template.HTML(data)
	}
}

func (s *Str2html) check(data string) error {
	// log.Debug("str2html raw", data)
	if data == "" {
		return fmt.Errorf("data is empty")
	}
	return nil
}
