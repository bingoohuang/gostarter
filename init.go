package main

import (
	"bytes"
	"fmt"
	"go-starter/util"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	_ "go-starter/statiq"

	"github.com/bingoohuang/gou"
	"github.com/bingoohuang/statiq/fs"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
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
		fmt.Printf("Built on %s from sha1 %s\n", util.Compile, util.Version)
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if *ipo {
		if err := ipoInit(); err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}

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

func ipoInit() error {
	sfs, err := fs.New()
	if err != nil {
		return err
	}

	if err = initCtl(sfs, "/ctl.tpl.sh", "./ctl"); err != nil {
		return err
	}

	if err = initConfigFile(sfs, "/config.tpl.toml", "./config.toml"); err != nil {
		return err
	}

	return nil
}

func initCtl(sfs *fs.StatiqFS, ctlTplName, ctlFilename string) error {
	if _, err := os.Stat(ctlFilename); err == nil {
		fmt.Println(ctlFilename + " already exists, ignored!")
		return nil
	} else if os.IsNotExist(err) {
		// continue
	} else {
		return err
	}

	ctl := string(sfs.Files[ctlTplName].Data)
	tpl, err := template.New(ctlTplName).Parse(ctl)
	if err != nil {
		return err
	}

	binArgs := argsExcludeInit()

	var content bytes.Buffer
	m := map[string]string{"BinName": os.Args[0], "BinArgs": strings.Join(binArgs, " ")}
	if err := tpl.Execute(&content, m); err != nil {
		return err
	}

	// 0755->即用户具有读/写/执行权限，组用户和其它用户具有读写权限；
	if err = ioutil.WriteFile(ctlFilename, content.Bytes(), 0755); err != nil {
		return err
	}

	fmt.Println(ctlFilename + " created!")
	return nil
}

func initConfigFile(sfs *fs.StatiqFS, configTplFileName, configFileName string) error {
	if _, err := os.Stat(configFileName); err == nil {
		fmt.Printf("%s already exists, ignored!\n", configFileName)
		return nil
	} else if os.IsNotExist(err) {
		// continue
	} else {
		return err
	}

	conf := sfs.Files[configTplFileName].Data
	// 0644->即用户具有读写权限，组用户和其它用户具有只读权限；
	if err := ioutil.WriteFile(configFileName, conf, 0644); err != nil {
		return err
	}

	fmt.Println(configFileName + " created!")

	return nil
}

func argsExcludeInit() []string {
	binArgs := make([]string, 0, len(os.Args)-2)
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if strings.Index(arg, "-i") == 0 || strings.Index(arg, "--init") == 0 {
			continue
		}

		if strings.Index(arg, "-") != 0 {
			arg = strconv.Quote(arg)
		}

		binArgs = append(binArgs, arg)
	}

	return binArgs
}
