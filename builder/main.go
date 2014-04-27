package main

import (
	builder "github.com/rafecolton/bob"
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
	"fmt"
	"os"
)

func init() {
}

var runtime *config.Runtime
var ver *version.Version
var par *parser.Parser
var runAsWorker = flag.Bool("work", false, "Run as a Goworker")

func main() {
	allTheThings := func(queue string, args ...interface{}) (errFake error) {
		if queue == "docker-build" {
			first := args[0].(map[string]interface{})
			pwd := first["pwd"].(string)
			build := fmt.Sprintf("%s/%s", pwd, first["build"].(string))

			os.Setenv("PWD", pwd)

			logger := log.Initialize(false)

			par, err := parser.NewParser(build, logger)
			if err != nil {
				logger.Println(
					color.Sprintf("@{r!}Alas, could not generate parser@{|}\n----> %+v", err),
				)
				gocleanup.Exit(73)
			}

			commandSequence, err := par.Parse()
			if err != nil {
				runtime.Println(color.Sprintf("@{r!}Alas, could not parse@{|}\n----> %+v", err))
				gocleanup.Exit(23)
			}

			bob := builder.NewBuilder(logger, true)
			bob.Builderfile = build

			if err = bob.Build(commandSequence); err != nil {
				logger.Println(
					color.Sprintf(
						"@{r!}Alas, I am unable to complete my assigned build because of...@{|}\n----> %+v",
						err,
					),
				)
				gocleanup.Exit(29)
			}

			return
		}

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
				runtime.Builderfile = "bob.toml"
			}
			// otherwise, build
			par, err := parser.NewParser(runtime.Builderfile, runtime)
			if err != nil {
				runtime.Println(
					color.Sprintf("@{r!}Alas, could not generate parser@{|}\n----> %+v", err),
				)
				gocleanup.Exit(73)
			}

			commandSequence, err := par.Parse()
			if err != nil {
				runtime.Println(color.Sprintf("@{r!}Alas, could not parse@{|}\n----> %+v", err))
				gocleanup.Exit(23)
			}

			bob := builder.NewBuilder(runtime, true)
			bob.Builderfile = runtime.Builderfile

			if err = bob.Build(commandSequence); err != nil {
				runtime.Println(
					color.Sprintf(
						"@{r!}Alas, I am unable to complete my assigned build because of...@{|}\n----> %+v",
						err,
					),
				)
				gocleanup.Exit(29)
			}
		}

		gocleanup.Exit(0)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "-work" {
		flag.Parse()
		goworker.Register("DockerBuild", allTheThings)
		if err := goworker.Work(); err != nil {
			fmt.Println(
				color.Sprintf("@{r!}Alas, something went wrong :'(@{|}\n----> %+v", err),
			)
		}
		gocleanup.Exit(0)
	}

	runtime = config.NewRuntime()
	ver = version.NewVersion()

	allTheThings("")
}
