package util

import (
	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
)

func InitLog() io.Writer {
	logdir := viper.GetString("logdir")
	if logdir != "" {
		if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
			logrus.Panicf("failed to create %s error %v\n", logdir, err)
			return os.Stdout
		}

		loglevel := viper.GetString("loglevel")
		return gou.InitLogger(loglevel, logdir, filepath.Base(os.Args[0])+".log")
	} else {
		logrus.SetLevel(logrus.DebugLevel)
		return os.Stdout
	}
}
