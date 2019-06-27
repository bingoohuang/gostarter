package app

import (
	"go-starter/ui"
	"go-starter/util"

	"github.com/bingoohuang/now"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type App struct {
	startupTime now.Now
	R           *gin.Engine
	UI          *ui.Context
}

func (a App) Start() {
	a.Route()
	a.Run()
}

func (a App) Run() {
	addr := viper.GetString("addr")
	logrus.Infof("go-starter started to run on addr %s\n", addr)
	if err := a.R.Run(addr); err != nil {
		panic(err)
	}
}

func CreateApp() *App {
	util.InitFlags()

	app := &App{}
	app.startupTime = now.MakeNow()

	app.R = util.InitGin(util.InitLog())
	app.UI = ui.CreateContext()

	return app
}
