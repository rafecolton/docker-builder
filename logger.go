package main

import (
	"github.com/Sirupsen/logrus"
)

func setLogger(level, format string) {
	switch level {
	case "debug", "d":
		Logger.Level = logrus.DebugLevel
	case "info", "i":
		Logger.Level = logrus.InfoLevel
	case "warn", "w":
		Logger.Level = logrus.WarnLevel
	case "error", "e":
		Logger.Level = logrus.ErrorLevel
	case "fatal", "f":
		Logger.Level = logrus.FatalLevel
	case "panic", "p":
		Logger.Level = logrus.PanicLevel
	default:
		Logger.Level = logrus.InfoLevel
	}

	switch format {
	case "text", "t":
		Logger.Formatter = new(logrus.TextFormatter)
	case "json", "j":
		Logger.Formatter = new(logrus.JSONFormatter)
	case "force-color", "fc":
		Logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	default:
		Logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	}

}
