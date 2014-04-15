package log

import (
	color "github.com/wsxiaoys/terminal/color"
	stdlog "log"
	"os"
)

/*
Log is the interface for all general logging methods.
*/
type Log interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

/*
BuilderLogger is the implementation of the Log interface for this project.
*/
type BuilderLogger struct {
	Log Log
}

/*
Initialize returns a BuilderLogger that either contains a null logger (that
prints nothing) or a standard logger (from the log package) with
project-specific output.
*/
func Initialize(quiet bool) *BuilderLogger {
	l := &BuilderLogger{}

	if quiet {
		l.Log = &nullLogger{}
	} else {
		l.Log = stdlog.New(os.Stderr, color.Sprint("@{g!}[builder] "), stdlog.LstdFlags)
	}
	return l
}
