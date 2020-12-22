package util

import (
	"fmt"

	"github.com/bingoohuang/gou/cnf"

	_ "net/http/pprof" // nolint
	"os"

	"github.com/bingoohuang/gou/htt"
	"github.com/bingoohuang/gou/lo"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// InitFlags initializes the plfags.
func InitFlags() {
	help := pflag.BoolP("help", "h", false, "help")
	ipo := pflag.BoolP("init", "i", false, "init to create template config file and ctl.sh")
	pflag.StringP("addr", "a", ":30057", "http address to listen and serve")

	pflag.BoolP("ui", "u", false, "enable simple ui")

	pprofAddr := htt.PprofAddrPflag()

	// Add more pflags can be set from command line
	// ...

	pflag.Parse()

	cnf.CheckUnknownPFlags()

	if *help {
		fmt.Printf("Built on %s from sha1 %s\n", Compile, Version)
		pflag.PrintDefaults()
		os.Exit(0)
	}

	Ipo(*ipo)
	htt.StartPprof(*pprofAddr)

	// 从当前位置读取config.toml配置文件
	viper.SetConfigName("cnf")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	lo.Err(viper.ReadInConfig())

	// Watching and re-reading config files
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viper.SetEnvPrefix("GOSTARTER")
	viper.AutomaticEnv()

	// 设置一些配置默认值
	// viper.SetDefault("InfluxAddr", "http://127.0.0.1:8086")
	// viper.SetDefault("CheckIntervalSeconds", 60)

	_ = viper.BindPFlags(pflag.CommandLine)
}
