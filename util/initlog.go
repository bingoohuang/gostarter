package util

import (
	"io"
	"os"
	"path/filepath"

	"github.com/bingoohuang/gou/str"

	"github.com/spf13/pflag"

	"github.com/bingoohuang/gou/lo"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// DeclareLogPFlags declares the log pflags.
func DeclareLogPFlags() {
	pflag.StringP("loglevel", "", "info", "debug/info/warn/error")
	pflag.StringP("logdir", "", "var/logs", "log dir")
	pflag.StringP("logfmt", "", "text/json", "default text")
	pflag.BoolP("logrus", "", true, "enable logrus")
}

// SetupLog setup log parameters.
func SetupLog() io.Writer {
	if !viper.GetBool("logrus") {
		logrus.SetLevel(logrus.DebugLevel)
		return os.Stdout
	}

	logfmt := viper.GetString("logfmt")

	if logfmt != "json" {
		// https://stackoverflow.com/a/48972299
		logrus.SetFormatter(&prefixed.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		})
	}

	loglevel := viper.GetString("loglevel")
	logrus.SetLevel(str.Decode(loglevel,
		"debug", logrus.DebugLevel,
		"info", logrus.InfoLevel,
		"warn", logrus.WarnLevel,
		"error", logrus.ErrorLevel,
		logrus.InfoLevel).(logrus.Level))

	logdir := viper.GetString("logdir")
	if logdir != "" {
		if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
			logrus.Panicf("failed to create %s error %v\n", logdir, err)
		}

		w := lo.InitLogger(loglevel, logdir, filepath.Base(os.Args[0])+".log")
		logrus.SetOutput(w)

		return w
	}

	return os.Stdout
}
