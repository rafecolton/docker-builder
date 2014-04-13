package main

import (
	"fmt"
	config "github.com/rafecolton/builder/config"
	//lint "github.com/rafecolton/builder/linter"
	"github.com/rafecolton/builder/version"
	//. "github.com/wsxiaoys/terminal/color"
)

func main() {

	runtime := config.New()

	// check for version/rev/branch options
	ver := version.New()
	infoCheck(runtime.Options, ver)

	//lint := linter.New(runtime)
	//_ = linter.Lint()

	//runtime.Println(data)
	//Log.Println(data)

	// parse file catch error and suggest -l / --lint
	// run the build with parsed data
}

func infoCheck(opts config.Options, ver *version.Version) {
	if opts.Version {
		fmt.Println(ver.Version)
	}

	if opts.VersionFull {
		fmt.Println(ver.VersionFull)
	}

	if opts.Branch {
		fmt.Println(ver.Branch)
	}

	if opts.Rev {
		fmt.Println(ver.Rev)
	}
}
