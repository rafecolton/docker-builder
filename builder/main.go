package main

import (
	. "github.com/rafecolton/builder/config"
	linter "github.com/rafecolton/builder/linter"
	versiontrick "github.com/rafecolton/builder/version"
	//. "github.com/wsxiaoys/terminal/color"
)

var (
	opts Options
)

func main() {

	runtime := Configure()

	// check for version/rev/branch options
	version := versiontrick.Init(runtime)
	infoAndExitCheck(runtime.Options, version)

	_ = linter.Lint(opts.Lintfile)
	//Log.Println(data)

	// parse file catch error and suggest -l / --lint
	// run the build with parsed data

}

func infoAndExitCheck(opts Options, version *versiontrick.VersionTrick) {
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
