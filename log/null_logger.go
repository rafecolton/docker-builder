package log

import stdlog "log"

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
