package app

import (
	"go-starter/demo"
	"go-starter/model"
	"go-starter/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (a App) Route() {
	r := a.R

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, model.Rsp{Status: 200, Message: "started time:" + a.startupTime.P})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, model.Rsp{Status: 200, Message: "pong, built time:" + util.Compile})
	})
	r.GET("/stats", func(c *gin.Context) {
		c.JSON(200, struct {
			StartupTime string
			Count       int
		}{
			StartupTime: a.startupTime.String(),
			Count:       100,
		})
	})

	// curl "http://127.0.0.1:30057/DemoWrapBindJSON" -X POST -d '{"name":"bingoohuang","age":100}' -H "Content-Type: application/json"
	r.POST("/DemoWrapBindJSON", demo.WrapBindJSON())
	// curl "http://127.0.0.1:30057/DemoWrapBindJSONRouter" -X POST -d '{"name":"bingoohuang","age":100}' -H "Content-Type: application/json"
	r.POST("/DemoWrapBindJSONRouter", demo.WrapBindJSONRouter())

	if viper.GetBool("ui") {
		g := r.Group("/", util.Auth)
		g.GET("/", a.UI.HomepageHandler)
		g.GET("/static/*filename", gin.WrapH(a.UI.StaticHandler()))
	}
}
