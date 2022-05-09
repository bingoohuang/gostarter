package ui

import (
	"fmt"
	"github.com/bingoohuang/gg/pkg/v"
	"github.com/bingoohuang/gostarter/pkg/ging"
	"io/fs"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	WebFS       fs.FS
	WebTemplate *template.Template
)

var WebFuncMap = template.FuncMap{
	"showData": func(t interface{}) string {
		return fmt.Sprintf("%+v", t)
	},
	"showTime": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02 15:04:05")
	},
}

func Handler(contextPath string) func(*gin.Context) {
	return func(g *gin.Context) {
		switch g.Request.URL.Path {
		case "/":
			IndexHandler(g)
		default:
			webFS := http.FileServer(http.FS(WebFS))
			if contextPath != "" {
				webFS = http.StripPrefix(contextPath, webFS)
			}
			webFS.ServeHTTP(g.Writer, g.Request)
		}
	}
}

// IndexHandler handles the homepage request.
func IndexHandler(g *gin.Context) {
	args := struct {
		Version string
	}{
		Version: v.Version(),
	}

	JSONOrTpl(args, "index.html", g)
}

// JSONOrTpl handles the JSON or HTML requests.
func JSONOrTpl(data interface{}, templateFile string, g *gin.Context) {
	if g.GetHeader("Accept") == "application/json" || g.Query("format") == "json" {
		ging.JSON(g, data)
		return
	}

	if err := WebTemplate.ExecuteTemplate(g.Writer, templateFile, data); err != nil {
		g.String(http.StatusInternalServerError, "%v", err)
	}
}
