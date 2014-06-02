package main

import (
	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/config"
	"github.com/modcloth/docker-builder/log"
	"github.com/modcloth/docker-builder/parser"
	"github.com/modcloth/docker-builder/version"
)

import (
	"fmt"

	"github.com/modcloth/queued-command-runner"
	"github.com/onsi/gocleanup"
	"github.com/wsxiaoys/terminal/color"
)

var runtime *config.Runtime
var ver *version.Version
var par *parser.Parser
var logger log.Logger

func main() {
	runtime = config.NewRuntime()
	ver = version.NewVersion()

	// if user requests version/branch/rev
	if runtime.Version {
		runtime.Println(ver.Version)
	} else if runtime.VersionFull {
		runtime.Println(ver.VersionFull)
	} else if runtime.Branch {
		runtime.Println(ver.Branch)
	} else if runtime.Rev {
		runtime.Println(ver.Rev)
	} else if runtime.Lintfile != "" {
		// lint
		par, _ = parser.NewParser(runtime.Lintfile, runtime)
		par.AssertLint()
	} else {
		if runtime.Builderfile == "" {
			runtime.Builderfile = "Bobfile"
		}
		// otherwise, build
		par, err := parser.NewParser(runtime.Builderfile, runtime)
		if err != nil {
			runtime.Println(
				color.Sprintf("@{r!}Alas, could not generate parser@{|}\n----> %q", err),
			)
			gocleanup.Exit(73)
		}

		commandSequence, err := par.Parse()
		if err != nil {
			runtime.Println(color.Sprintf("@{r!}Alas, could not parse@{|}\n----> %q", err))
			gocleanup.Exit(23)
		}

		bob, err := builder.NewBuilder(runtime, true)
		if err != nil {
			runtime.Println(
				color.Sprintf(
					"@{r!}Alas, I am unable to complete my assigned build because of...@{|}\n----> %q",
					err,
				),
			)
			gocleanup.Exit(61)
		}

		bob.Builderfile = runtime.Builderfile

		if err = bob.Build(commandSequence); err != nil {
			runtime.Println(
				color.Sprintf(
					"@{r!}Alas, I am unable to complete my assigned build because of...@{|}\n----> %q",
					err,
				),
			)
			gocleanup.Exit(29)
		}
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
