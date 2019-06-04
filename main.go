package main

import (
	"go-starter/ui"

	"github.com/bingoohuang/gou"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	defer gou.Recover()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

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
