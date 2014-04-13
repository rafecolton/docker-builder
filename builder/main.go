package main

import (
	config "github.com/rafecolton/builder/config"
	lint "github.com/rafecolton/builder/linter"
	versiontrick "github.com/rafecolton/builder/version"
	//. "github.com/wsxiaoys/terminal/color"
)

func main() {

	runtime := config.New()

	// check for version/rev/branch options
	version := versiontrick.Init(runtime)
	infoAndExitCheck(runtime.Options, version)

	linter := lint.New(runtime)
	_ = linter.Lint()

	//runtime.Println(data)
	//Log.Println(data)

	// parse file catch error and suggest -l / --lint
	// run the build with parsed data
}

func infoAndExitCheck(opts config.Options, version *versiontrick.VersionTrick) {
	if opts.Version {
		version.VersionAndExit()
	}

	if opts.VersionFull {
		version.VersionFullAndExit()
	}

	if opts.Branch {
		version.BranchAndExit()
	}

	if opts.Rev {
		version.RevAndExit()
	}
}
