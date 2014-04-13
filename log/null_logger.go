package log

import stdlog "log"

type nullLogger struct{}

func (me *nullLogger) Print(v ...interface{})                 {}
func (me *nullLogger) Printf(format string, v ...interface{}) {}
func (me *nullLogger) Println(v ...interface{})               {}

func (me *nullLogger) Fatal(v ...interface{}) {
	stdlog.Fatal(v...)
}
func (me *nullLogger) Fatalf(format string, v ...interface{}) {
	stdlog.Fatalf(format, v...)
}
func (me *nullLogger) Panicf(format string, v ...interface{}) {
	stdlog.Panicf(format, v...)
}
func (me *nullLogger) Panicln(v ...interface{}) {
	stdlog.Panicln(v...)
}
