package main

import (
	"github.com/modcloth/bob/config"
	"github.com/modcloth/bob/log"
	"github.com/modcloth/bob/parser"
	"github.com/modcloth/bob/version"
)

import (
	"github.com/onsi/gocleanup"
)

import (
	"flag"
	"os"
)

var runtime *config.Runtime
var ver *version.Version
var par *parser.Parser
var logger log.Logger
var runAsWorker = flag.Bool("work", false, "Run as a Goworker")

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-work" {
		runWorker()
	} else {
		allTheThings()
	}

	gocleanup.Exit(0)
}
