package main

import (
	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/log"
	builderlogger "github.com/modcloth/docker-builder/log"
	"github.com/modcloth/docker-builder/parser"
	"github.com/modcloth/docker-builder/version"
)

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/modcloth/queued-command-runner"
	"github.com/onsi/gocleanup"
	"github.com/wsxiaoys/terminal/color"
)

var ver = version.NewVersion()
var par *parser.Parser
var logger log.Logger

func main() {
	app := cli.NewApp()
	app.Name = "builder"
	app.Usage = "builder (a.k.a. \"Bob\") builds Docker images from a friendly config file"
	app.Version = fmt.Sprintf("%s %s", ver.Version, app.Compiled)
	app.Flags = []cli.Flag{
		cli.BoolFlag{"branch", "print branch and exit"},
		cli.BoolFlag{"rev", "print revision and exit"},
		cli.BoolFlag{"version-short", "print long version and exit"},
		cli.BoolFlag{"quiet, q", "produce no output, only exit codes"},
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
		logger = builderlogger.Initialize(c.GlobalBool("quiet"))
		return nil
	}
	app.Commands = []cli.Command{
		//{
		//Name:      "init",
		//ShortName: "i",
		//Usage:     "",
		//Action:    initialize,
		//},
		{
			Name:        "build",
			ShortName:   "b",
			Usage:       "build <file> - build Docker images from the provided Bobfile",
			Description: "Build Docker images from the provided Bobfile.",
			Action:      build,
		},
		{
			Name:        "lint",
			ShortName:   "l",
			Usage:       "lint <file> - validates whether or not your Bobfile is parsable",
			Description: "Validate whether or not your Bobfile is parsable.",
			Action:      lint,
		},
	}

	app.Run(os.Args)
	gocleanup.Exit(0)
}

func lint(c *cli.Context) {
	par, _ = parser.NewParser(c.Args().First(), logger)
	par.AssertLint()
}

func build(c *cli.Context) {
	builderfile := c.Args().First()
	if builderfile == "" {
		builderfile = "Bobfile"
	}

	par, err := parser.NewParser(builderfile, logger)
	if err != nil {
		exitErr(73, "@{r!}Alas, could not generate parser@{|}\n----> %q", err)
	}

	commandSequence, err := par.Parse()
	if err != nil {
		exitErr(23, "@{r!}Alas, could not parse@{|}\n----> %q", err)
	}

	bob, err := builder.NewBuilder(logger, true)
	if err != nil {
		exitErr(61, "@{r!}Alas, I am unable to complete my assigned build because of...@{|}\n----> %q", err)
	}

	bob.Builderfile = builderfile

	if err = bob.Build(commandSequence); err != nil {
		exitErr(29, "@{r!}Alas, I am unable to complete my assigned build because of...@{|}\n----> %q", err)
	}

	if builder.WaitForPush {
	WaitForPush:
		for {
			select {
			case <-runner.Done:
				break WaitForPush
			case err := <-runner.Errors:
				exitErr(1, "@{r!}Uh oh, something went wrong while running %q@{|}\n----> %q", err.CommandStr, err)
			}
		}
	}
}

//func initialize(c *cli.Context) {
//}

func exitErr(exitCode int, fmtString string, args ...interface{}) {
	logger.Println(color.Sprintf(fmtString, args...))
	gocleanup.Exit(exitCode)
}
