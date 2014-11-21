package main

import (
	"github.com/sylphon/builder-core/unit-config"

	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

func lint(c *cli.Context) {
	_, err := unitconfig.ReadFromFile("./"+c.Args().First(), unitconfig.TOML)
	if err != nil {
		Logger.Error(err)
		gocleanup.Exit(1)
	}
	gocleanup.Exit(0)
}
