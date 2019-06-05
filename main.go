package main

import (
	"go-starter/ui"

	"github.com/bingoohuang/now"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var startupTime = now.MakeNow()

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/panic", func(c *gin.Context) { panic("panic " + now.MakeNow().String()) })
	r.GET("/health", func(c *gin.Context) { c.String(200, "OK "+now.MakeNow().String()) })
	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	r.GET("/stats", func(c *gin.Context) {
		c.JSON(200, struct {
			StartupTime string
			Count       int
		}{
			StartupTime: startupTime.String(),
			Count:       100,
		})
	})

	if viper.GetBool("ui") {
		g := r.Group("/", auth)
		g.GET("/", ui.HomepageHandler)

		g.GET("/static/*filename", gin.WrapH(ui.StaticHandler()))
	}

	addr := viper.GetString("addr")
	logrus.Infof("go-starter started to run on addr %s\n", addr)
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
