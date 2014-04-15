package main

import (
	builder "github.com/rafecolton/bob"
	"github.com/rafecolton/bob/config"
	"github.com/rafecolton/bob/parser"
	"github.com/rafecolton/bob/version"
)

import (
	"os"
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
		os.Exit(0)
	} else if runtime.VersionFull {
		runtime.Println(ver.VersionFull)
		os.Exit(0)
	} else if runtime.Branch {
		runtime.Println(ver.Branch)
		os.Exit(0)
	} else if runtime.Rev {
		runtime.Println(ver.Rev)
		os.Exit(0)
	}

	// does linting
	if runtime.Lintfile != "" {
		par = parser.NewParser(runtime.Lintfile, runtime)
		par.AssertLint()

		os.Exit(0)
	}

	// does building
	if runtime.Builderfile != "" {
		par = parser.NewParser(runtime.Builderfile, runtime)

		//ignoring the error because we elect to blow up in the parser instead
		file, _ := par.Parse(true)

		bob := builder.NewBuilder()

		_ = bob.Build(file)
		os.Exit(0)
	}

	//otherwise, nothing to do!
	config.Usage()
}
