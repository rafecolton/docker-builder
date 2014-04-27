package main

import (
	"github.com/rafecolton/bob/config"
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
	"fmt"
	"os"
)

var runtime *config.Runtime
var ver *version.Version
var par *parser.Parser
var runAsWorker = flag.Bool("work", false, "Run as a Goworker")

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-work" {
		flag.Parse()
		goworker.Register("DockerBuild", allTheThings)

		if err := goworker.Work(); err != nil {
			fmt.Println(
				color.Sprintf("@{r!}Alas, something went wrong :'(@{|}\n----> %+v", err),
			)
		}
	} else {
		runtime = config.NewRuntime()
		ver = version.NewVersion()

		_ = allTheThings("")
	}

	gocleanup.Exit(0)
}
