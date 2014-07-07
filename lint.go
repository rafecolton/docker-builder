package main

import (
	"github.com/modcloth/docker-builder/parser"

	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

func lint(c *cli.Context) {
	par, _ = parser.NewParser(c.Args().First(), Logger)
	exitCode := par.AssertLint()
	gocleanup.Exit(exitCode)
}
