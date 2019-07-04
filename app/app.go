package app

import (
	"go-starter/ui"
	"go-starter/util"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bingoohuang/now"
	"github.com/facebookgo/grace/gracehttp"
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

	// restart by self
	server := &http.Server{Addr: addr, Handler: a.R}

	updatePidFile()

	logrus.Infof("go-starter started to run on addr %s", addr)
	if err := gracehttp.Serve(server); err != nil {
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

// kill -USR2 {pid} 会执行重启
func updatePidFile() {
	pidFile := "var/pid"
	envPidFile := os.Getenv("PID_FILE")
	if envPidFile != "" {
		pidFile = envPidFile
	}

	bytes, err := ioutil.ReadFile(pidFile)
	if err != nil {
		logrus.Errorf("read pid file error %s", err.Error())
		return
	}

	oldPid, err := strconv.Atoi(strings.TrimSpace(string(bytes)))
	if err != nil {
		logrus.Errorf("trans pid file error %s", err.Error())
		return
	}

	logrus.Infof("old pid is %d, new pid is %d", os.Getpid(), oldPid)

	if os.Getpid() != oldPid {
		_ = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	}
}
