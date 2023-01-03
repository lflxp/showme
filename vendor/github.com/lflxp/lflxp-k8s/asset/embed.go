package asset

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"path"

	"github.com/lflxp/lflxp-k8s/core/middlewares/jwt/framework"
	"github.com/lflxp/lflxp-k8s/core/middlewares/jwt/services"

	log "github.com/go-eden/slf4go"

	"github.com/gin-gonic/gin"
)

//go:embed script
var cleanscriptDir embed.FS

//go:embed docs
var docs embed.FS

//go:embed d2admin
var d2admin embed.FS

//go:embed node_modules
var nodeModules embed.FS

func RegisterAsset(router *gin.Engine) {
	router.Any("/d2admin/*any", func(c *gin.Context) {
		// staticServer := wrapHandler(http.FS(dashboard))
		staticServer := http.FileServer(http.FS(d2admin))
		// TODO: 遇到404 就返回前端根路径下index.html的资源 路径不变
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/node_modules/*any", func(c *gin.Context) {
		// staticServer := wrapHandler(http.FS(dashboard))
		staticServer := http.FileServer(http.FS(nodeModules))
		// TODO: 遇到404 就返回前端根路径下index.html的资源 路径不变
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	router.Any("/script/*any", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(cleanscriptDir))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	needUserGroup := router.Group("")
	needUserGroup.Use(framework.NewGinJwtMiddlewares(services.AllUserAuthorizator).MiddlewareFunc())
	needUserGroup.Any("/docs/*any", func(c *gin.Context) {
		staticServer := wrapHandlerDocs(http.FS(docs))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Request.URL.Path = "/base/favicon.ico"
		router.HandleContext(c)
	})
}

// 自定义FileServer
type NotFoundRedirectRespWr struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
}

func (w *NotFoundRedirectRespWr) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *NotFoundRedirectRespWr) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}
	return len(p), nil // Lie that we successfully written it
}

func wrapHandlerDocs(h http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hserver := http.FileServer(h)
		nfrw := &NotFoundRedirectRespWr{ResponseWriter: w}
		hserver.ServeHTTP(nfrw, r)
		if nfrw.status == 404 {
			log.Debugf("Redirecting %s to index.html.", r.RequestURI)
			// hardCode
			path := "/docs/index.html"
			f, err := h.Open(path)
			if err != nil {
				msg, code := toHTTPError(err)
				http.Error(w, msg, code)
				return
			}
			defer f.Close()

			d, err := f.Stat()
			if err != nil {
				msg, code := toHTTPError(err)
				http.Error(w, msg, code)
				return
			}

			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.Header().Add("Accept-Ranges", "bytes")
			w.Header().Add("Content-Length", string(d.Size()))
			http.ServeContent(w, r, path, d.ModTime(), f)
		}
	}
}

func wrapHandler(h http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hserver := http.FileServer(h)
		nfrw := &NotFoundRedirectRespWr{ResponseWriter: w}
		hserver.ServeHTTP(nfrw, r)
		if nfrw.status == 404 {
			log.Debugf("Redirecting %s to index.html.", r.RequestURI)
			// hardCode
			path := "/dashboard/index.html"
			f, err := h.Open(path)
			if err != nil {
				msg, code := toHTTPError(err)
				http.Error(w, msg, code)
				return
			}
			defer f.Close()

			d, err := f.Stat()
			if err != nil {
				msg, code := toHTTPError(err)
				http.Error(w, msg, code)
				return
			}

			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.Header().Add("Accept-Ranges", "bytes")
			w.Header().Add("Content-Length", string(d.Size()))
			http.ServeContent(w, r, path, d.ModTime(), f)
		}
	}
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found", http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

// AssetHandler returns an http.Handler that will serve files from
// the Assets embed.FS. When locating a file, it will strip the given
// prefix from the request and prepend the root to the filesystem.
func AssetHandler(prefix string, assets embed.FS, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)

		// If we can't find the asset, fs can handle the error
		file, err := assets.Open(assetPath)
		if err != nil {
			return nil, err
		}

		// Otherwise assume this is a legitimate request routed correctly
		return file, err
	})

	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}
