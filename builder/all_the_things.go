package main

import (
	builder "github.com/rafecolton/bob"
	"github.com/rafecolton/bob/config"
	"github.com/rafecolton/bob/parser"
	"github.com/rafecolton/bob/version"
)

import (
	"github.com/onsi/gocleanup"
	"github.com/wsxiaoys/terminal/color"
)

import (
	"fmt"
	"os"
)

var allTheThings = func(queue string, args ...interface{}) (fake error) {
	if queue == "docker-build" {
		first := args[0].(map[string]interface{})
		pwd := first["pwd"].(string)
		build := fmt.Sprintf("%s/%s", pwd, first["build"].(string))

		os.Setenv("PWD", pwd)

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
