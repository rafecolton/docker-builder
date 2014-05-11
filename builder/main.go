package main

import (
	builder "github.com/modcloth/bob"
	"github.com/modcloth/bob/config"
	"github.com/modcloth/bob/log"
	"github.com/modcloth/bob/parser"
	"github.com/modcloth/bob/version"
)

import (
	"github.com/modcloth/queued-command-runner"
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
var logger log.Logger
var runAsWorker = flag.Bool("work", false, "Run as a Goworker")

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-work" {
		runWorker()
	} else {
		allTheThings()
	}

	if builder.WaitForPush {
	WaitForPush:
		for {
			select {
			case <-runner.Done:
				break WaitForPush
			case err := <-runner.Errors:
				fmt.Println(
					color.Sprintf("@{r!}Uh oh, something went wrong while running %q@{|}\n----> %q", err.CommandStr, err),
				)
				gocleanup.Exit(1)
			}
		}
	}

	gocleanup.Exit(0)
}
