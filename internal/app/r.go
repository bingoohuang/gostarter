package app

import (
	"net/http"

	"github.com/bingoohuang/gostarter/demo"
	"github.com/bingoohuang/gostarter/internal/model"
	"github.com/bingoohuang/gostarter/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Route defines the routing of http.
func (a App) Route() {
	r := a.R

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.Rsp{Status: http.StatusOK, Message: "started time:" + a.startupTime.P})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.Rsp{Status: http.StatusOK, Message: "pong, built time:" + util.Compile})
	})
	r.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			StartupTime string
			Count       int
		}{
			StartupTime: a.startupTime.String(),
			Count:       100,
		})
	})

	// curl "http://127.0.0.1:30057/DemoWrapBindJSON" -X POST  \
	// -d '{"name":"bingoohuang","age":100}' -H "Content-Type: application/json"
	r.POST("/DemoWrapBindJSON", demo.WrapBindJSON())
	// curl "http://127.0.0.1:30057/DemoWrapBindJSONRouter" -X POST  \
	// -d '{"name":"bingoohuang","age":100}' -H "Content-Type: application/json"
	r.POST("/DemoWrapBindJSONRouter", demo.WrapBindJSONRouter())

	if viper.GetBool("ui") {
		g := r.Group("/", util.Auth)
		g.GET("/", a.UI.HomepageHandler)
		g.GET("/static/*filename", gin.WrapH(a.UI.StaticHandler()))
	}
}
