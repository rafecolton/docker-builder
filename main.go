package main

import (
	"fmt"
	"os"

	"github.com/rafecolton/docker-builder/conf"
	"github.com/rafecolton/docker-builder/parser"
	"github.com/rafecolton/docker-builder/server"
	"github.com/rafecolton/docker-builder/version"

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
	app.Version = ver.Version + " " + app.Compiled.String()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "branch",
			Usage: "print branch and exit",
		},
		cli.BoolFlag{
			Name:  "rev",
			Usage: "print revision and exit",
		},
		cli.BoolFlag{
			Name:  "version-short",
			Usage: "print long version and exit",
		},
		cli.BoolFlag{Name: "quiet, q",
			Usage: "produce no output, only exit codes",
		},
		cli.StringFlag{
			Name:  "log-level, l",
			Value: conf.Config.LogLevel,
			Usage: "log level (options: debug/d, info/i, warn/w, error/e, fatal/f, panic/p)",
		},
		cli.StringFlag{
			Name:  "log-format, f",
			Value: conf.Config.LogFormat,
			Usage: "log output format (options: text/t, json/j)",
		},
		cli.StringFlag{
			Name:  "dockercfg-un",
			Value: conf.Config.CfgUn,
			Usage: "Docker registry username",
		},
		cli.StringFlag{
			Name:  "dockercfg-pass",
			Value: conf.Config.CfgPass,
			Usage: "Docker registry password",
		},
		cli.StringFlag{
			Name:  "dockercfg-email",
			Value: conf.Config.CfgEmail,
			Usage: "Docker registry email",
		},
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
				cli.BoolFlag{
					Name:  "skip-push",
					Usage: "override Bobfile behavior and do not push any images (useful for testing)",
				},
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
				cli.IntFlag{
					Name:  "port, p",
					Value: conf.Config.Port,
					Usage: "port on which to serve",
				},
				cli.StringFlag{
					Name:  "api-token, t",
					Value: "",
					Usage: "GitHub API token",
				},
				cli.BoolFlag{
					Name:  "skip-push",
					Usage: "override Bobfile behavior and do not push any images (useful for testing)",
				},
				cli.StringFlag{
					Name:  "username",
					Value: "",
					Usage: "username for basic auth",
				},
				cli.StringFlag{
					Name:  "password",
					Value: "",
					Usage: "password for basic auth",
				},
				cli.StringFlag{
					Name:  "travis-token",
					Value: "",
					Usage: "Travis API token for webhooks",
				},
				cli.StringFlag{
					Name:  "github-secret",
					Value: "",
					Usage: "GitHub secret for webhooks",
				},
				cli.BoolFlag{
					Name:  "no-travis",
					Usage: "do not include route for Travis CI webhook",
				},
				cli.BoolFlag{
					Name:  "no-github",
					Usage: "do not include route for GitHub webhook",
				},
			},
		},
	}

	app.Run(os.Args)
	gocleanup.Exit(0)
}
