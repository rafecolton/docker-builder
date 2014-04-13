package main

import (
	config "github.com/rafecolton/builder/config"
	"github.com/rafecolton/builder/version"
)

var runtime config.Runtime
var ver version.Version

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
		runtime.Println(ver.Version)
	}

	if opts.VersionFull {
		runtime.Println(ver.VersionFull)
	}

	if opts.Branch {
		runtime.Println(ver.Branch)
	}

	if opts.Rev {
		runtime.Println(ver.Rev)
	}
}
