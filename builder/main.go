package main

import (
	"fmt"
	//. "github.com/rafecolton/builder"
	flags "github.com/jessevdk/go-flags"
	linter "github.com/rafecolton/builder/linter"
	versiontrick "github.com/rafecolton/builder/version"
	"os"
)

type Options struct {
	// Inform and Exit
	Version     bool `short:"v" description:"Print version and exit"`
	VersionFull bool `long:"version" description:"Print long version and exit"`
	Branch      bool `short:"b" long:"branch" description:"Print branch and exit"`
	Rev         bool `short:"r" long:"rev" description:"Print revision and exit"`

	// Runtime Options
	Quiet bool `short:"q" long:"quiet" description:"Produce no output, only exit codes"`

	// Features
	Lintfile    string `short:"l" long:"lint" descrpition:"Lint the provided file. Compatible with -q/--quiet"`
	Builderfile string `short:"f" long:"builderfile" descrpition:"The configuration file for Builder"`
}

var opts Options

func main() {

	// parse args, make sure we exit 0 with -h/--help
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		arg1 := os.Args[1]
		if arg1 == "-h" || arg1 == "--help" {
			os.Exit(0)
		} else {
			fmt.Println("Unable to parse args")
			os.Exit(3)
		}
	}

	// check for version/rev/branch options
	version := versiontrick.Init()
	InfoAndExitCheck(opts, version)

	data := linter.Lint(opts.Lintfile)
	fmt.Println(data)

	// parse file catch error and suggest -l / --lint
	// run the build with parsed data

}

func InfoAndExitCheck(opts Options, version *versiontrick.VersionTrick) {
	if opts.Version {
		version.VersionAndExit()
	} else if opts.VersionFull {
		version.VersionFullAndExit()
	} else if opts.Branch {
		version.BranchAndExit()
	} else if opts.Rev {
		version.RevAndExit()
	}

}
