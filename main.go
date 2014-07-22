package main

import (
	"fmt"
	"os"

	"github.com/modcloth/docker-builder/conf"
	"github.com/modcloth/docker-builder/parser"
	"github.com/modcloth/docker-builder/server"
	"github.com/modcloth/docker-builder/version"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/kelseyhightower/envconfig"
	"github.com/modcloth/kamino"
	"github.com/onsi/gocleanup"
)

var ver = version.NewVersion()
var par *parser.Parser

//Logger is the logger for the docker-builder main
var Logger *logrus.Logger

func init() {
	// parse env config
	if err := envconfig.Process("docker_builder", &conf.Config); err != nil {
		Logger.WithField("err", err).Fatal("envconfig error")
	}

	// set default config port
	if conf.Config.Port == 0 {
		conf.Config.Port = 5000
	}

	// set logger defaults
	Logger = logrus.New()
	Logger.Formatter = &logrus.TextFormatter{ForceColors: true}
}

func main() {
	app := cli.NewApp()
	app.Name = "docker-builder"
	app.Usage = "docker-builder (a.k.a. \"Bob\") builds Docker images from a friendly config file"
	app.Version = fmt.Sprintf("%s %s", ver.Version, app.Compiled)
	app.Flags = []cli.Flag{
		cli.BoolFlag{"branch", "print branch and exit"},
		cli.BoolFlag{"rev", "print revision and exit"},
		cli.BoolFlag{"version-short", "print long version and exit"},
		cli.BoolFlag{"quiet, q", "produce no output, only exit codes"},
		cli.StringFlag{"log-level, l", conf.Config.LogLevel, "log level (options: debug/d, info/i, warn/w, error/e, fatal/f, panic/p)"},
		cli.StringFlag{"log-format, f", conf.Config.LogFormat, "log output format (options: text/t, json/j)"},
		cli.StringFlag{"dockercfg-un", conf.Config.CfgUn, "Docker registry username"},
		cli.StringFlag{"dockercfg-pass", conf.Config.CfgPass, "Docker registry password"},
		cli.StringFlag{"dockercfg-email", conf.Config.CfgEmail, "Docker registry email"},
	}
	app.Action = func(c *cli.Context) {
		ver = version.NewVersion()
		if c.GlobalBool("branch") {
			fmt.Println(ver.Branch)
		} else if c.GlobalBool("rev") {
			fmt.Println(ver.Rev)
		} else if c.GlobalBool("version-short") {
			fmt.Println(ver.Version)
		} else {
			cli.ShowAppHelp(c)
		}
	}
	app.Before = func(c *cli.Context) error {
		logLevel := c.String("log-level")
		logFormat := c.String("log-format")

		setLogger(logLevel, logFormat)
		kamino.Logger = Logger

		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:        "init",
			Usage:       "init [dir] - initialize the given directory (default '.') with a Bobfile",
			Description: "Make educated guesses to fill out a Bobfile given a directory with a Dockerfile",
			Action:      initialize,
		},
		{
			Name:        "build",
			Usage:       "build [file] - build Docker images from the provided Bobfile",
			Description: "Build Docker images from the provided Bobfile.",
			Action:      build,
			Flags: []cli.Flag{
				cli.BoolFlag{"skip-push", "override Bobfile behavior and do not push any images (useful for testing)"},
			},
		},
		{
			Name:        "lint",
			Usage:       "lint [file] - validates whether or not your Bobfile is parsable",
			Description: "Validate whether or not your Bobfile is parsable.",
			Action:      lint,
		},
		{
			Name:        "serve",
			Usage:       "serve <options> - start a small HTTP web server for receiving build requests",
			Description: server.Description,
			Action:      func(c *cli.Context) { server.Logger(Logger); server.Serve(c) },
			Flags: []cli.Flag{
				cli.IntFlag{"port, p", conf.Config.Port, "port on which to serve"},
				cli.StringFlag{"api-token, t", "", "GitHub API token"},
				cli.BoolFlag{"skip-push", "override Bobfile behavior and do not push any images (useful for testing)"},
				cli.StringFlag{"username", "", "username for basic auth"},
				cli.StringFlag{"password", "", "password for basic auth"},
				cli.StringFlag{"travis-token", "", "Travis API token for webhooks"},
				cli.StringFlag{"github-secret", "", "GitHub secret for webhooks"},
				cli.BoolFlag{"no-travis", "do not include route for Travis CI webhook"},
				cli.BoolFlag{"no-github", "do not include route for GitHub webhook"},
			},
		},
	}

	app.Run(os.Args)
	gocleanup.Exit(0)
}
