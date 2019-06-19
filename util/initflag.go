package util

import (
	"fmt"
	// pprof debug
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"

	"github.com/bingoohuang/gou"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitFlags() {
	help := pflag.BoolP("help", "h", false, "help")
	ipo := pflag.BoolP("init", "i", false, "init to create template config file and ctl.sh")
	pflag.StringP("addr", "a", ":30057", "http address to listen and serve")
	pflag.StringP("loglevel", "l", "info", "debug/info/warn/error")
	pflag.StringP("logdir", "d", "./var", "log dir")
	pflag.BoolP("logrus", "o", true, "enable logrus")
	pflag.BoolP("ui", "u", false, "enable simple ui")
	pprofAddr := gou.PprofAddrPflag()

	// Add more pflags can be set from command line
	// ...

	pflag.Parse()

	args := pflag.Args()
	if len(args) > 0 {
		fmt.Printf("Unknown args %s\n", strings.Join(args, " "))
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if *help {
		fmt.Printf("Built on %s from sha1 %s\n", Compile, Version)
		pflag.PrintDefaults()
		os.Exit(0)
	}

	Ipo(*ipo)
	gou.StartPprof(*pprofAddr)

	// 从当前位置读取config.toml配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	gou.LogErr(viper.ReadInConfig())

	// Watching and re-reading config files
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viper.SetEnvPrefix("GO_STARTER")
	viper.AutomaticEnv()

	// 设置一些配置默认值
	// viper.SetDefault("InfluxAddr", "http://127.0.0.1:8086")
	// viper.SetDefault("CheckIntervalSeconds", 60)

	_ = viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("logrus") {
		logdir := viper.GetString("logdir")
		if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
			logrus.Panicf("failed to create %s error %v\n", logdir, err)
		}

		loglevel := viper.GetString("loglevel")
		gou.InitLogger(loglevel, logdir, filepath.Base(os.Args[0])+".log")
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
