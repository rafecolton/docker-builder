package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

var un string
var pwd string

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

func setUnAndPwd(c *cli.Context) {
	// lowest priority

	// ENV
	un = config.Username
	pwd = config.Password

	// command line
	cliUn := c.String("username")
	cliPwd := c.String("password")

	// if username passed on command line, use cl one instead
	if cliUn != "" {
		un = cliUn
	}

	// if password passed on command line, used cl one instead
	if cliPwd != "" {
		pwd = cliPwd
	}

	// highest priority
}
