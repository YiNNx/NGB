package util

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"ngb/config"
	"path"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = getLogger()
	Logger.Info("logger started")
}

func getLogger() *logrus.Logger {
	logger := logrus.New()

	logger.Formatter = new(logrus.JSONFormatter)
	logger.SetReportCaller(true)
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
