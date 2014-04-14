package config

import (
	"fmt"
	"os"
)

import (
	flags "github.com/jessevdk/go-flags"
	builderlogger "github.com/rafecolton/builder/log"
)

var (
	parser *flags.Parser
	opts   Options
)

func Usage() {
	parser.WriteHelp(os.Stderr)

}

func New() *Runtime {

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

	me := &Runtime{
		Quiet:   opts.Quiet,
		Options: opts,
		Log:     logger.Log,
	}

	return me
}

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

type Runtime struct {
	Quiet bool
	builderlogger.Log
	Options
}

func (me *Runtime) Print(v ...interface{}) {
	me.Log.Print(v...)
}

func (me *Runtime) Println(v ...interface{}) {
	me.Log.Println(v...)
}

func (me *Runtime) Printf(format string, v ...interface{}) {
	me.Log.Printf(format, v...)
}
