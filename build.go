package main

import (
	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/parser"

	"github.com/codegangsta/cli"
)

func build(c *cli.Context) {
	builder.SkipPush = c.Bool("skip-push")
	builderfile := c.Args().First()
	if builderfile == "" {
		builderfile = "Bobfile"
	}

	par, err := parser.NewParser(builderfile, Logger)
	if err != nil {
		exitErr(73, "unable to generate parser", err)
	}

	commandSequence, err := par.Parse()
	if err != nil {
		exitErr(23, "unable to parse", err)
	}

	bob, err := builder.NewBuilder(Logger, true)
	if err != nil {
		exitErr(61, "unable to build", err)
	}

	bob.Builderfile = builderfile

	if err = bob.Build(commandSequence); err != nil {
		exitErr(29, "unable to build", err)
	}
}
