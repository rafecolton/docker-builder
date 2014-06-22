package main

import (
	"github.com/modcloth/docker-builder/parser"

	"github.com/codegangsta/cli"
)

func lint(c *cli.Context) {
	par, _ = parser.NewParser(c.Args().First(), Logger)
	par.AssertLint()
}
