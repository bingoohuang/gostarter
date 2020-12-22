package app

import (
	"github.com/bingoohuang/gostarter/internal/util"
	"net/http"

	"github.com/bingoohuang/gou/sy"

	"github.com/bingoohuang/gostarter/internal/ui"
	"github.com/bingoohuang/now"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// App wraps the application info.
type App struct {
	startupTime now.Now
	R           *gin.Engine
	UI          *ui.Context
}

// CreateApp creates the application.
func CreateApp() *App {
	util.InitFlags()

	app := &App{}
	app.startupTime = now.MakeNow()

	app.R = util.InitGin()
	app.UI = ui.CreateContext()

	return app
}

// Start starts the application.
func (a App) Start() {
	a.Route()
	a.run()
}

func (a App) run() {
	addr := viper.GetString("addr")

	// restart by self
	server := &http.Server{Addr: addr, Handler: a.R}

	if err := sy.UpdatePidFile("var/pid"); err != nil {
		logrus.Warnf("UpdatePidFile error %v", err)
	}

	logrus.Infof("gostarter started to run on addr %s", addr)

	if err := gracehttp.Serve(server); err != nil {
		panic(err)
	}
}
