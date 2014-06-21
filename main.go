package main

import (
	"github.com/modcloth/docker-builder/parser"
	"github.com/modcloth/docker-builder/version"
)

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

var ver = version.NewVersion()
var par *parser.Parser

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
		cli.StringFlag{"log-level, l", "info", "log level (options: debug/d, info/i, warn/w, error/e, fatal/f, panic/p)"},
		cli.StringFlag{"log-format, f", "text", "log output format (options: text/t, json/j)"},
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
		setLogger(c.String("log-level"), c.String("log-format"))
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:        "init",
			ShortName:   "i",
			Usage:       "init [dir] - initialize the given directory (default '.') with a Bobfile",
			Description: "Make educated guesses to fill out a Bobfile given a directory with a Dockerfile",
			Action:      initialize,
		},
		{
			Name:        "build",
			ShortName:   "b",
			Usage:       "build [file] - build Docker images from the provided Bobfile",
			Description: "Build Docker images from the provided Bobfile.",
			Action:      build,
			Flags: []cli.Flag{
				cli.BoolFlag{"skip-push", "override Bobfile behavior and do not push any images (useful for testing)"},
			},
		},
		{
			Name:        "lint",
			ShortName:   "l",
			Usage:       "lint [file] - validates whether or not your Bobfile is parsable",
			Description: "Validate whether or not your Bobfile is parsable.",
			Action:      lint,
		},
	}

	app.Run(os.Args)
	gocleanup.Exit(0)
}
