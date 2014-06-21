package main

import (
	"github.com/modcloth/docker-builder/analyzer"
	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/parser"
	"github.com/modcloth/docker-builder/version"
)

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/modcloth/queued-command-runner"
	"github.com/onsi/gocleanup"
)

var ver = version.NewVersion()
var par *parser.Parser
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
		exitErr(73, "unable to generate parser", err)
	}

	commandSequence, err := par.Parse()
	if err != nil {
		exitErr(23, "unable to parse", err)
	}

	bob, err := builder.NewBuilder(logger, true)
	if err != nil {
		exitErr(61, "unable to build", err)
	}

	bob.Builderfile = builderfile

	if err = bob.Build(commandSequence); err != nil {
		exitErr(29, "unable to build", err)
	}

	if builder.WaitForPush {
	WaitForPush:
		for {
			select {
			case <-runner.Done:
				break WaitForPush
			case err := <-runner.Errors:
				exitErr(1, "error when running push command", map[string]interface{}{"command": err.CommandStr, "error": err})
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
		exitErr(1, "unable to create Bobfile", err)
	}

	bobfilePath := filepath.Join(dir, "Bobfile")

	//no error when stating, file already exists, rename with timestamp
	if _, err := os.Stat(bobfilePath); err == nil {
		bobfilePath = fmt.Sprintf("%s.%d", bobfilePath, int32(time.Now().Unix()))
	}

	outfile, err := os.Create(bobfilePath)
	if err != nil {
		exitErr(86, "unable to create output file", map[string]interface{}{"output_file": bobfilePath, "error": err})
	}
	defer outfile.Close()

	// TODO: figure out why this isn't getting written by the toml encoder
	dockerSectionHeader := []byte("[docker]\n\n")
	if _, err := outfile.Write(dockerSectionHeader); err != nil {
		exitErr(127, "unable to write to output file", map[string]interface{}{"output_file": bobfilePath, "error": err})
	}

	encoder := toml.NewEncoder(outfile)
	if err = encoder.Encode(file); err != nil {
		exitErr(123, "unable to write to output file", map[string]interface{}{"output_file": bobfilePath, "error": err})
	}

	vimFtComment := []byte("\n\n# vim:ft=toml")
	if _, err := outfile.Write(vimFtComment); err != nil {
		exitErr(127, "unable to write to output file", map[string]interface{}{"output_file": bobfilePath, "error": err})
	}

	logger.WithFields(logrus.Fields{"output_file": bobfilePath}).Info("successfully initialized")
}

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

	logger.WithFields(fields).Error(message)
	gocleanup.Exit(exitCode)
}
