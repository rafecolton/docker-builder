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
//"fmt"
//"os"
)

var runtime *config.Runtime
var ver *version.Version
var par *parser.Parser

func main() {

	runtime = config.NewRuntime()
	ver = version.NewVersion()

	// if user requests version/branch/rev
	if runtime.Version {
		runtime.Println(ver.Version)
		gocleanup.Exit(0)
	} else if runtime.VersionFull {
		runtime.Println(ver.VersionFull)
		gocleanup.Exit(0)
	} else if runtime.Branch {
		runtime.Println(ver.Branch)
		gocleanup.Exit(0)
	} else if runtime.Rev {
		runtime.Println(ver.Rev)
		gocleanup.Exit(0)
	}

	// does linting
	if runtime.Lintfile != "" {
		par, _ = parser.NewParser(runtime.Lintfile, runtime)
		par.AssertLint()

		gocleanup.Exit(0)
	}

	// does building
	if runtime.Builderfile != "" {
		par, err := parser.NewParser(runtime.Builderfile, runtime)
		if err != nil {
			runtime.Println(color.Sprintf("@{r!}Alas@{|}, could not generate parser\n----> %+v", err))
			gocleanup.Exit(73)
		}

		commandSequence, err := par.Parse()
		if err != nil {
			runtime.Println(color.Sprintf("@{r!}Alas@{|}, could not parse\n----> %+v", err))
			gocleanup.Exit(23)
		}

		bob := builder.NewBuilder(runtime, true)

		if err = bob.Build(commandSequence); err != nil {
			runtime.Println(err)
			gocleanup.Exit(29)
		}

		err = bob.Build(commandSequence)
		if err != nil {
			runtime.Println(color.Sprintf("@{r!}Alas@{|}, I am unable to complete my assigned build\n----> %+v", err))
			gocleanup.Exit(41)
		}

		gocleanup.Exit(0)
	}

	//otherwise, nothing to do!
	config.Usage()
	gocleanup.Exit(2)
}
