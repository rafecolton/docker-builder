package main

import (
	builder "github.com/modcloth/bob"
	"github.com/modcloth/bob/config"
	"github.com/modcloth/bob/log"
	"github.com/modcloth/bob/parser"
	"github.com/modcloth/bob/version"
)

import (
	"github.com/benmanns/goworker"
	"github.com/onsi/gocleanup"
	"github.com/wsxiaoys/terminal/color"
)

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var runWorker = func() {
	logger = log.Initialize(false)
	flag.Parse()
	goworker.Register("DockerBuild", workerFunc)

	if err := goworker.Work(); err != nil {
		logger.Println(
			color.Sprintf("@{r!}Alas, something went wrong :'(@{|}\n----> %+v", err),
		)
	}
}

var workerFunc = func(queue string, args ...interface{}) (fake error) {
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

		bob, err := builder.NewBuilder(logger, true)
		if err != nil {
			logger.Println(
				color.Sprintf(
					"@{r!}Alas, I am unable to complete my assigned build because of...@{|}\n----> %+v",
					err,
				),
			)
			gocleanup.Exit(61)
		}
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

	return errors.New("invalid attempt to use as a goworker")
}

var allTheThings = func() {
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

		bob, err := builder.NewBuilder(runtime, true)
		if err != nil {
			runtime.Println(
				color.Sprintf(
					"@{r!}Alas, I am unable to complete my assigned build because of...@{|}\n----> %+v",
					err,
				),
			)
			gocleanup.Exit(61)
		}

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
}