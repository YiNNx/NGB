package util

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger := logrus.New()

	Logger.Formatter = new(logrus.JSONFormatter)
	Logger.SetReportCaller(true)
	Logger.SetLevel(logrus.DebugLevel)

	filePath := path.Join("log", "log-")
	writer, err := rotatelogs.New(
		filePath+".%Y-%m-%d-%H-%M",
		rotatelogs.WithLinkName(filePath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		Logger.Fatal(err)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})

	Logger.AddHook(lfHook)

	Logger.Info("logger started successfully")
}
