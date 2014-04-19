package main

import (
	//builder "github.com/rafecolton/bob"
	"github.com/rafecolton/bob/config"
	"github.com/rafecolton/bob/parser"
	"github.com/rafecolton/bob/version"
)

import (
	"fmt"
	//"github.com/wsxiaoys/terminal/color"
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
		par, _ = parser.NewParser(runtime.Lintfile, runtime)
		par.AssertLint()

		os.Exit(0)
	}

	// does building
	if runtime.Builderfile != "" {
		//par = parser.NewParser(runtime.Builderfile, runtime)

		//instructions, err := par.Parse()
		//if err != nil {
		////TODO: print something here
		//os.Exit(23)
		//}

		//bob := builder.NewBuilder()

		//_ = bob.Build(instructions)

		os.Exit(0)
	}

	fmt.Printf("runtime: %+v\n", runtime)
	par, _ := parser.NewParser(runtime.Builderfile, runtime)

	_, _ = par.Parse()

	//otherwise, nothing to do!
	//config.Usage()
}
