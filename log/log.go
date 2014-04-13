package log

import (
	color "github.com/wsxiaoys/terminal/color"
	stdlog "log"
	"os"
)

type Log interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

type BuilderLogger struct {
	Log Log
}

func Initialize(quiet bool) *BuilderLogger {
	l := &BuilderLogger{}

	if quiet {
		l.Log = &nullLogger{}
	} else {
		l.Log = stdlog.New(os.Stderr, color.Sprint("@{g!}[builder] "), stdlog.LstdFlags)
	}
	return l
}
