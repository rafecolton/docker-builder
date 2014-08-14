package main

import (
	"github.com/modcloth/docker-builder/parser"

	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

func lint(c *cli.Context) {
	p := parser.NewParser(c.Args().First(), Logger)
	if _, err := p.Parse(); err != nil {
		p.Error(err.Error())
		gocleanup.Exit(err.ExitCode())
	}
	gocleanup.Exit(0)
}
