package main

import (
	"github.com/Sirupsen/logrus"
)

func setLogger(level, format string) {
	switch level {
	case "debug", "d":
		Logger.Level = logrus.Debug
	case "info", "i":
		Logger.Level = logrus.Info
	case "warn", "w":
		Logger.Level = logrus.Warn
	case "error", "e":
		Logger.Level = logrus.Error
	case "fatal", "f":
		Logger.Level = logrus.Fatal
	case "panic", "p":
		Logger.Level = logrus.Panic
	default:
		Logger.Level = logrus.Info
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
