package main

import (
	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/parser"

	"github.com/codegangsta/cli"
	"github.com/modcloth/queued-command-runner"
)

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
