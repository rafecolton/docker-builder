package log

import stdlog "log"

type nullLogger struct{}

func (nl *nullLogger) Print(v ...interface{})                 {}
func (nl *nullLogger) Printf(format string, v ...interface{}) {}
func (nl *nullLogger) Println(v ...interface{})               {}

func (nl *nullLogger) Fatal(v ...interface{}) {
	stdlog.Fatal(v...)
}
func (nl *nullLogger) Fatalln(v ...interface{}) {
	stdlog.Fatalln(v...)
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
func (nl *nullLogger) Write(p []byte) (int, error) {
	return len(p), nil
}

/*
NullLogger is an exported symbol for the nullLogger struct.
*/
type NullLogger struct{}

// Print prints.
func (nl *NullLogger) Print(v ...interface{}) {}

// Printf printfs.
func (nl *NullLogger) Printf(format string, v ...interface{}) {}

// Println printlns.
func (nl *NullLogger) Println(v ...interface{}) {}

// Fatal fatals.
func (nl *NullLogger) Fatal(v ...interface{}) {
	stdlog.Fatal(v...)
}

// Fatalln fatallns.
func (nl *NullLogger) Fatalln(v ...interface{}) {
	stdlog.Fatalln(v...)
}

// Fatalf fatalfs.
func (nl *NullLogger) Fatalf(format string, v ...interface{}) {
	stdlog.Fatalf(format, v...)
}

// Panicf panicfs.
func (nl *NullLogger) Panicf(format string, v ...interface{}) {
	stdlog.Panicf(format, v...)
}

// Panicln paniclns.
func (nl *NullLogger) Panicln(v ...interface{}) {
	stdlog.Panicln(v...)
}

// Write writes.
func (nl *NullLogger) Write(p []byte) (int, error) {
	return len(p), nil
}
