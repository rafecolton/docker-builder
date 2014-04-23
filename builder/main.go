package main

import (
	builder "github.com/rafecolton/bob"
	"github.com/rafecolton/bob/config"
	"github.com/rafecolton/bob/parser"
	"github.com/rafecolton/bob/version"
)

import (
	"github.com/onsi/gocleanup"
	//"github.com/wsxiaoys/terminal/color"
)

import (
	"fmt"
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
		//par = parser.NewParser(runtime.Builderfile, runtime)

		//instructions, err := par.Parse()
		//if err != nil {
		////TODO: print something here
		//gocleanup.Exit(23)
		//}

		gocleanup.Exit(0)
	}

	bob := builder.NewBuilder(runtime, true)
	fmt.Println(bob.LatestImageTaggedWithUUID("foo"))

	//_ = bob.Build(instructions)

	//otherwise, nothing to do!
	//config.Usage()
}
