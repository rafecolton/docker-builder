package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/onsi/gocleanup"
)

func exitErr(exitCode int, message string, args interface{}) {
	var fields logrus.Fields

	switch args.(type) {
	case error:
		fields = logrus.Fields{
			"error": args,
		}
	case map[string]interface{}:
		fields = args.(map[string]interface{})
	}

	Logger.WithFields(fields).Error(message)
	gocleanup.Exit(exitCode)
}
