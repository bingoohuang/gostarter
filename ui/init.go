package ui

import (
	"html/template"
	"net/http"

	"github.com/bingoohuang/statiq/fs"
	"github.com/gin-gonic/gin"

	_ "go-starter/statiq"
)

var (
	BuiltTime string
	Sha1ver   string
)

func InitSha1verBuiltTime(sha1ver, builtTime string) {
	Sha1ver = sha1ver
	BuiltTime = builtTime
}

var StatiqFS *fs.StatiqFS

var homepageTpl *template.Template

func init() {
	StatiqFS, _ = fs.New()

	homepageTpl = loadTmpl("/homepage.html")

}

func loadTmpl(name string) *template.Template {
	res := string(StatiqFS.Files[name].Data)
	return template.Must(template.New(name).Funcs(fnMap).Parse(res))
}

func StaticHandler() http.Handler {
	return http.StripPrefix("/static", http.FileServer(StatiqFS))
}

func HomepageHandler(c *gin.Context) {
	args := struct {
		Sha1ver   string
		BuiltTime string
	}{
		Sha1ver:   Sha1ver,
		BuiltTime: BuiltTime,
	}
	JSONOrTpl(args, homepageTpl, c)
}

func JSONOrTpl(args interface{}, tpl *template.Template, c *gin.Context) {
	fmt := c.Query("format")
	if fmt == "json" {
		c.JSON(200, args)
	} else if err := tpl.Execute(c.Writer, args); err != nil {
		c.String(500, "%v", err)
	}
}
