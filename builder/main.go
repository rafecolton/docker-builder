package main

import (
	"github.com/rafecolton/builder/config"
	"github.com/rafecolton/builder/parser"
	"github.com/rafecolton/builder/version"
)

//import (
//"os"
//"path"
//)

var runtime *config.Runtime
var ver *version.Version
var par *parser.Parser

func main() {

	runtime = config.New()
	ver = version.New()
	par = parser.New()
	opts := runtime.Options

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

	par.Builderfile = "./spec/fixtures/Builderfile"
	bf, _ := par.Parse()

	runtime.Printf("%+v\n", bf)
	runtime.Println(bf.Docker.BuildOpts)

}
