package pages

import (
	"embed"
	"html/template"
	"net/http"

	t1 "github.com/lflxp/lflxp-music/core/middlewares/template"

	"github.com/gin-gonic/gin"
)

//go:embed static
var Static embed.FS

//go:embed views
var Templates embed.FS

func RegisterTemplate(router *gin.Engine) {
	if router == nil {
		panic("router nil")
	}
	router.StaticFS("/adminfs", http.FS(Static))
	// 基于embed注册templates模板
	t := template.Must(template.New("").Funcs(t1.FuncMap).ParseFS(Templates, "views/*/*"))
	router.SetHTMLTemplate(t)
}
