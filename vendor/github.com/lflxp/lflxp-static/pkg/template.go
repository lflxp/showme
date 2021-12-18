package pkg

import (
	"html/template"

	"github.com/gin-contrib/multitemplate"
)

var (
	HtmlTemp                      multitemplate.Render
	htmlTemplate, elementTemplate []byte
	err                           error
)

func init() {
	htmlTemplate, err = Static.ReadFile("static/video.html")
	if err != nil {
		panic(err)
	}
	elementTemplate, err = Static.ReadFile("static/main.html")
	if err != nil {
		panic(err)
	}

	HtmlTemp = multitemplate.New()
	registerTemplate(string(htmlTemplate), "index")
	registerTemplate(string(elementTemplate), "element")
}

func registerTemplate(temp, name string) {
	t, _ := template.New(name).Parse(temp)
	HtmlTemp.Add(name, t)
}
