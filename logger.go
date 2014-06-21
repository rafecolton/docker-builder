package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/modcloth/queued-command-runner"
)

var logger *logrus.Logger

func init() {
	//set logger defaults
	logger = logrus.New()
	logger.Formatter = &logrus.TextFormatter{ForceColors: true}

	//set logger defaults from env
	setLogger(os.Getenv("DOCKER_BUILDER_LOG_LEVEL"), os.Getenv("DOCKER_BUILDER_LOG_FORMAT"))
	runner.Logger = logger
}

func setLogger(level, format string) {
	switch level {
	case "debug", "d":
		logger.Level = logrus.Debug
	case "info", "i":
		logger.Level = logrus.Info
	case "warn", "w":
		logger.Level = logrus.Warn
	case "error", "e":
		logger.Level = logrus.Error
	case "fatal", "f":
		logger.Level = logrus.Fatal
	case "panic", "p":
		logger.Level = logrus.Panic
	default:
		logger.Level = logrus.Info
	}

	switch format {
	case "text", "t":
		logger.Formatter = new(logrus.TextFormatter)
	case "json", "j":
		logger.Formatter = new(logrus.JSONFormatter)
	case "force-color", "fc":
		logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	default:
		logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	}

}
