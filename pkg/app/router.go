package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bingoohuang/gostarter/pkg/controllers"

	"github.com/bingoohuang/gostarter/pkg/model/cook"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/vibrantbyte/go-antpath/antpath"

	"github.com/bingoohuang/gg/pkg/ginx/adapt"
	"github.com/bingoohuang/gg/pkg/ginx/anyfn"
	"github.com/bingoohuang/golog/pkg/ginlogrus"
	"github.com/bingoohuang/gostarter/pkg/conf"

	"github.com/bingoohuang/gg/pkg/logx"
	"github.com/bingoohuang/gostarter/pkg/ging"
	"github.com/bingoohuang/gostarter/pkg/model"
	"github.com/bingoohuang/gostarter/pkg/ui"
	"github.com/gin-gonic/gin"
)

var startupTime = time.Now().Format(`2006-01-02 15:03:04`)

// StartWeb start the web
func StartWeb() {
	r := gin.New()
	r.Use(ginlogrus.Logger(nil, true), gin.Recovery())
	r.Use(sessions.Sessions("gostarter", cookie.NewStore([]byte("PAS6roKzT0ES3BpuChp7hacQAmR9soka"))), loginFilter)

	r.GET("/health", func(c *gin.Context) {
		ging.JSON(c, model.Rsp{Status: http.StatusOK, Message: "started time:" + startupTime})
	})
	r.GET("/ping", func(c *gin.Context) {
		ging.JSON(c, model.Rsp{Status: http.StatusOK, Message: "pong"})
	})
	r.GET("/stats", func(c *gin.Context) {
		ging.JSON(c, struct {
			StartupTime string
			Count       int
		}{
			StartupTime: startupTime,
			Count:       100,
		})
	})

	ar := adapt.Adapt(r)
	af := anyfn.NewAdapter()
	registerWrappers(af)
	ar.RegisterAdapter(af)

	controllers.Register(ar, af)

	g := r.Group("/", ui.Auth)
	g.GET("/", ui.Handler(""))

	addr := fmt.Sprintf(":%d", conf.Conf.Port)
	log.Printf("start to listen on %s", addr)
	err := r.Run(addr)
	logx.Fatalf(err, "run failed: %v", err)
}

func loginFilter(c *gin.Context) {
	if user, _ := cook.GetLogin(c); user != nil {
		c.Set(model.LoginUser, user)
		return
	}

	if pathMatchersAny(c.Request.URL.Path, conf.Conf.IgnorePaths) {
		log.Printf("D! 忽略授权 %s ", c.Request.URL.Path)
		return
	}

	log.Printf("W! 未授权访问 %s ", c.Request.URL.Path)
	c.JSON(http.StatusUnauthorized, nil)
	c.Abort()
}

var antMatcher = antpath.New()

func pathMatchersAny(uri string, paths []string) bool {
	for _, p := range paths {
		if antMatcher.Match(p, uri) {
			return true
		}
	}

	return false
}

func GetLogin(c *gin.Context) *model.Login {
	if v, _ := c.Get(model.LoginUser); v != nil {
		return v.(*model.Login)
	}
	return nil
}
