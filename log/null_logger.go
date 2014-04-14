package log

import stdlog "log"

type nullLogger struct{}

func (nl *nullLogger) Print(v ...interface{})                 {}
func (nl *nullLogger) Printf(format string, v ...interface{}) {}
func (nl *nullLogger) Println(v ...interface{})               {}

func (nl *nullLogger) Fatal(v ...interface{}) {
	stdlog.Fatal(v...)
}
func (nl *nullLogger) Fatalf(format string, v ...interface{}) {
	stdlog.Fatalf(format, v...)
}
func (nl *nullLogger) Panicf(format string, v ...interface{}) {
	stdlog.Panicf(format, v...)
}
func (nl *nullLogger) Panicln(v ...interface{}) {
	stdlog.Panicln(v...)
}
