package main

import (
	"github.com/rafecolton/bob/config"
	"github.com/rafecolton/bob/parser"
	"github.com/rafecolton/bob/version"
	"github.com/wsxiaoys/terminal/color"
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
	par = parser.NewParser()
	opts := runtime.Options

	// if user requests version/branch/rev
	if opts.Version {
		runtime.Println(ver.Version)
		os.Exit(0)
	} else if opts.VersionFull {
		runtime.Println(ver.VersionFull)
		os.Exit(0)
	} else if opts.Branch {
		runtime.Println(ver.Branch)
		os.Exit(0)
	} else if opts.Rev {
		runtime.Println(ver.Rev)
		os.Exit(0)
	}

	//does linting!
	if opts.Lintfile != "" {
		par.Builderfile = opts.Lintfile

		_, err := par.Parse()
		if err != nil {
			runtime.Println(color.Sprintf("@{r!}Alas@{|}, %s is not a valid Builderfile!\n----> %+v", opts.Lintfile, err))
			os.Exit(5)
		}

		runtime.Printf(color.Sprintf("@{g!}Hooray@{|}, %s is a valid Builderfile!", opts.Lintfile))
		os.Exit(0)
	}

	if opts.Builderfile != "" {
		par.Builderfile = opts.Builderfile
		file, err := par.Parse()
		if err != nil {
			runtime.Println(color.Sprintf("@{r!}Alas@{|}, cannot build %s, it is not a valid Builderfile!\n----> %+v", opts.Builderfile, err))
			os.Exit(7)
		}

		bill := NewBuilder()

		_ = bill.Build(file)
	}

	//otherwise, nothing to do!
	config.Usage()
}
