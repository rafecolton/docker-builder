package main

import (
	"github.com/rafecolton/bob/config"
	"github.com/rafecolton/bob/log"
	"github.com/rafecolton/bob/parser"
	"github.com/rafecolton/bob/version"
)

import (
	"github.com/benmanns/goworker"
	"github.com/onsi/gocleanup"
	"github.com/wsxiaoys/terminal/color"
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
	logger = log.Initialize(false)

	if len(os.Args) > 1 && os.Args[1] == "-work" {
		flag.Parse()
		goworker.Register("DockerBuild", allTheThings)

		if err := goworker.Work(); err != nil {
			logger.Println(
				color.Sprintf("@{r!}Alas, something went wrong :'(@{|}\n----> %+v", err),
			)
		}
	} else {
		_ = allTheThings("")
	}

	gocleanup.Exit(0)
}
