package util

import (
	"io"
	"os"
	"path/filepath"
	"regexp"

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
	pflag.StringP("logfmt", "", "", "log format(text/json), default text")
	pflag.BoolP("logrus", "", true, "enable logrus")
}

// TextFormatter extends the prefixed.TextFormatter with line joining.
type TextFormatter struct {
	prefixed.TextFormatter
	JoinLinesSeparator string
}

var reNewLines = regexp.MustCompile(`\r?\n`) // nolint

// Format formats the log output.
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.JoinLinesSeparator != "" {
		entry.Message = reNewLines.ReplaceAllString(entry.Message, f.JoinLinesSeparator)
	}

	return f.TextFormatter.Format(entry)
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
		logrus.SetFormatter(&TextFormatter{
			TextFormatter: prefixed.TextFormatter{
				DisableColors:   true,
				TimestampFormat: "2006-01-02 15:04:05",
				FullTimestamp:   true,
				ForceFormatting: true,
			},
			JoinLinesSeparator: `\n`,
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
