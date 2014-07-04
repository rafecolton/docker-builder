package main

import (
	"github.com/modcloth/docker-builder/server"

	"github.com/codegangsta/cli"
)

func serve(context *cli.Context) {
	server.Logger(Logger)
	server.Serve(context)
}
