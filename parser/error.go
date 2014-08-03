package parser

// Error is an interface for any error types returned by the parser
// package / during the parsing process
type Error interface {
	// Error returns the error message to satisfy the error interface
	Error() string

	// ExitCode returns the code that should be used when exiting as a result
	// of this error
	ExitCode() int
}

// OSPathError is used for any instance of os.PathError that is encountered
// during parsing
type OSPathError struct {
	error
}

// ExitCode returns the exit code parser errors related to os.PathError.  It is
// the same value for all OSPathError instances
func (err *OSPathError) ExitCode() int {
	return 17
}

// TOMLParseError is used for errors related to parsing a .toml file
type TOMLParseError struct {
	error
}

// ExitCode returns the exit code for toml parsing errors.  It is the same
// value for all toml parsing errors.
func (err *TOMLParseError) ExitCode() int {
	return 5
}

// BuilderfileConvertError is used for errors encountered while converting
// Bobfile versions
type BuilderfileConvertError struct {
	error
}

// ExitCode returns the exit code Bobfile conversion errors.  It is the same
// value for all Bobfile conversion errors.
func (err *BuilderfileConvertError) ExitCode() int {
	return 7
}
