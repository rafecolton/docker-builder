package main

import (
	"github.com/modcloth/docker-builder/analyzer"
	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/parser"
	"github.com/modcloth/docker-builder/version"
)

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/modcloth/queued-command-runner"
	"github.com/onsi/gocleanup"
	"github.com/wsxiaoys/terminal/color"
)

var ver = version.NewVersion()
var par *parser.Parser

var log = logrus.New()
var logger = struct {
	*logrus.Logger
	io.Writer
}{log, log.Out}

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
		logger.Level = logrus.Debug
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

func lint(c *cli.Context) {
	par, _ = parser.NewParser(c.Args().First(), logger)
	par.AssertLint()
}

func build(c *cli.Context) {
	builder.SkipPush = c.Bool("skip-push")
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

func initialize(c *cli.Context) {
	dir := c.Args().First()
	if dir == "" {
		dir = "."
	}

	file, err := analyzer.ParseAnalysisFromDir(dir)
	if err != nil {
		exitErr(1, "@{r!}Unable to initialize Bobfile@{|}\n----> %q", err)
	}

	bobfilePath := filepath.Join(dir, "Bobfile")

	//no error when stating, file already exists, rename with timestamp
	if _, err := os.Stat(bobfilePath); err == nil {
		bobfilePath = fmt.Sprintf("%s.%d", bobfilePath, int32(time.Now().Unix()))
	}

	outfile, err := os.Create(bobfilePath)
	if err != nil {
		exitErr(86, "@{r!}Unable to create output file %q@{|}\n----> %q", bobfilePath, err)
	}
	defer outfile.Close()

	// TODO: figure out why this isn't getting written by the toml encoder
	dockerSectionHeader := []byte("[docker]\n\n")
	if _, err := outfile.Write(dockerSectionHeader); err != nil {
		exitErr(127, "@{r!}Unable to write to output file %q@{|}\n----> %q", bobfilePath, err)
	}

	encoder := toml.NewEncoder(outfile)
	if err = encoder.Encode(file); err != nil {
		exitErr(123, "@{r!}Unable to write to output file %q@{|}\n----> %q", bobfilePath, err)
	}

	vimFtComment := []byte("\n\n# vim:ft=toml")
	if _, err := outfile.Write(vimFtComment); err != nil {
		exitErr(127, "@{r!}Unable to write to output file %q@{|}\n----> %q", bobfilePath, err)
	}

	logger.Printf("Successfully created %q\n", bobfilePath)
}

func exitErr(exitCode int, fmtString string, args ...interface{}) {
	logger.Println(color.Sprintf(fmtString, args...))
	gocleanup.Exit(exitCode)
}
