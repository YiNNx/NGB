package util

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"ngb/config"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = getLogger()
	Logger.Info("logger started")
}

func getLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.SetFormatter(formatter())
	logger.SetLevel(logrus.DebugLevel)

	baseLogPath := path.Join(config.C.Log.Path, config.C.Log.File)
	writer, err := rotatelogs.New(
		baseLogPath+"-%Y-%m-%d",
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logger.Fatal(err)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})

	logger.AddHook(lfHook)
	return logger
}

func formatter() *nested.Formatter {
	fmtter := &nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(" [%v:%v]", filepath.Base(fullPath), line)
		},
	}
	fmtter.NoColors = false
	return fmtter
}

//var stdFormatter  &prefixed.TextFormatter{
//	FullTimestamp:   true,
//	TimestampFormat: "2006-01-02.15:04:05.000000",
//	ForceFormatting: true,
//	ForceColors:     true,
//	DisableColors:   false,
//}
