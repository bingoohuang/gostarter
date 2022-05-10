package app

import (
	"net/http"
	"time"

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

	g := r.Group("/", ui.Auth)
	g.GET("/", ui.Handler(""))

	err := r.Run(":8080")
	logx.Fatalf(err, "run failed: %v", err)
}
