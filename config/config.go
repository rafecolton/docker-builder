package config

import (
	"fmt"
	"os"
)

import (
	flags "github.com/jessevdk/go-flags"
	builderlogger "github.com/rafecolton/bob/log"
)

var (
	parser *flags.Parser
	opts   Options
)

/*
Usage is like running the builder with -h/--help - it simply prints the usage
message to stderr.
*/
func Usage() {
	parser.WriteHelp(os.Stderr)

}

/*
Runtime is a struct of convenience, used for keeping track of our conf options
(i.e. passed on the command line or specified otherwise) as well as other
useful, global-ish things.
*/
type Runtime struct {
	Quiet bool
	builderlogger.Log
	Options
}

/*
NewRuntime returns a new Runtime struct instance that contains all of the
global-ish things specific to this invokation of builder.
*/
func NewRuntime() *Runtime {
	parser = flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		arg1 := os.Args[1]
		if arg1 == "-h" || arg1 == "--help" {
			os.Exit(0)
		} else {
			fmt.Println("Unable to parse args")
			os.Exit(3)
		}
	}

	logger := builderlogger.Initialize(opts.Quiet)

	runtime := &Runtime{
		Quiet:   opts.Quiet,
		Options: opts,
		Log:     logger.Log,
	}

	return runtime
}

/*
Options are our command line options, set using the
https://github.com/jessevdk/go-flags library.
*/
type Options struct {
	// Inform and Exit
	Version     bool `short:"v" description:"Print version and exit"`
	VersionFull bool `long:"version" description:"Print long version and exit"`
	Branch      bool `long:"branch" description:"Print branch and exit"`
	Rev         bool `long:"rev" description:"Print revision and exit"`

	// Runtime Options
	Quiet bool `short:"q" long:"quiet" description:"Produce no output, only exit codes" default:"false"`

	// Features
	Lintfile    string `short:"l" long:"lint" description:"Lint the provided file. Compatible with -q/--quiet"`
	Builderfile string `short:"f" long:"builderfile" description:"The configuration file for Builder"`
}

/*
Print passes through calls to Print to logger owned by the Runtime object.
Used primarily as a convenience.
*/
func (config *Runtime) Print(v ...interface{}) {
	config.Log.Print(v...)
}

/*
Println passes through calls to Println to logger owned by the Runtime object.
Used primarily as a convenience.
*/
func (config *Runtime) Println(v ...interface{}) {
	config.Log.Println(v...)
}

/*
Printf passes through calls to Printf to logger owned by the Runtime object.
Used primarily as a convenience.
*/
func (config *Runtime) Printf(format string, v ...interface{}) {
	config.Log.Printf(format, v...)
}
