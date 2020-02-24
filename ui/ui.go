package ui

import (
	"html/template"
	"net/http"

	"github.com/bingoohuang/gostarter/util"

	"github.com/bingoohuang/statiq/fs"
	"github.com/gin-gonic/gin"

	// import static resources
	_ "github.com/bingoohuang/gostarter/statiq"
)

// Context defines the context of UI.
type Context struct {
	sfs         *fs.StatiqFS
	homepageTpl *template.Template
	fnMap       template.FuncMap
}

// CreateContext creates the Context.
func CreateContext() *Context {
	c := &Context{}
	c.sfs, _ = fs.New()
	c.homepageTpl = c.loadTmpl("/homepage.html")
	c.fnMap = template.FuncMap{
		"showData": showData,
		"showTime": showTime,
	}

	return c
}

func (c Context) loadTmpl(name string) *template.Template {
	res := string(c.sfs.Files[name].Data)
	return template.Must(template.New(name).Funcs(c.fnMap).Parse(res))
}

// StaticHandler handles the static resources requests.
func (c Context) StaticHandler() http.Handler {
	return http.StripPrefix("/static", http.FileServer(c.sfs))
}

// HomepageHandler handles the homepage request.
func (c Context) HomepageHandler(g *gin.Context) {
	args := struct {
		Sha1ver   string
		BuiltTime string
	}{
		Sha1ver:   util.Version,
		BuiltTime: util.Compile,
	}
	c.JSONOrTpl(args, c.homepageTpl, g)
}

// JSONOrTpl handles the JSON or HTML requests.
func (c Context) JSONOrTpl(args interface{}, tpl *template.Template, g *gin.Context) {
	fmt := g.Query("format")
	if fmt == "json" {
		g.JSON(http.StatusOK, args)
	} else if err := tpl.Execute(g.Writer, args); err != nil {
		g.String(http.StatusInternalServerError, "%v", err)
	}
}
