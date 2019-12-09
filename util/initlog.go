package util

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"

	"github.com/bingoohuang/gou/str"

	"github.com/spf13/pflag"

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
	loglevel := viper.GetString("loglevel")
	level := str.Decode(loglevel,
		"debug", logrus.DebugLevel,
		"info", logrus.InfoLevel,
		"warn", logrus.WarnLevel,
		"error", logrus.ErrorLevel,
		logrus.InfoLevel).(logrus.Level)
	logrus.SetLevel(level)

	var formatter logrus.Formatter

	if viper.GetString("logfmt") != "json" {
		// https://stackoverflow.com/a/48972299
		formatter = &TextFormatter{
			TextFormatter: prefixed.TextFormatter{
				DisableColors:   true,
				TimestampFormat: "2006-01-02 15:04:05",
				FullTimestamp:   true,
				ForceFormatting: true,
			},
			JoinLinesSeparator: `\n`,
		}
	}

	if !viper.GetBool("logrus") {
		logrus.SetFormatter(formatter)

		return os.Stdout
	}

	logdir := viper.GetString("logdir")
	if logdir != "" {
		if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
			logrus.Panicf("failed to create %s error %v\n", logdir, err)
		}

		return initLogger(level, logdir, filepath.Base(os.Args[0])+".log", formatter)
	}

	logrus.SetFormatter(formatter)

	return os.Stdout
}

// 参考链接： https://tech.mojotv.cn/2018/12/27/golang-logrus-tutorial
func initLogger(level logrus.Level, logDir, filename string, formatter logrus.Formatter) io.Writer {
	baseLogPath := path.Join(logDir, filename)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		logrus.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}

	logrus.SetLevel(level)

	//writerMap := lfshook.WriterMap{
	//	logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
	//	logrus.InfoLevel:  writer,
	//	logrus.WarnLevel:  writer,
	//	logrus.ErrorLevel: writer,
	//	logrus.FatalLevel: writer,
	//	logrus.PanicLevel: writer,
	//}

	writerMap := writer

	logrus.AddHook(lfshook.NewHook(writerMap, formatter))
	logrus.SetOutput(writer)

	return writer
}
