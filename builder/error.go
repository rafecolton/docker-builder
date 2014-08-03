package builder

// Error is an interface for any error types returned by the builder
// package / during the building process
type Error interface {
	// Error returns the error message to satisfy the error interface
	Error() string

	// ExitCode returns the code that should be used when exiting as a result
	// of this error
	ExitCode() int
}

// SanitizeError is used for errors related to sanitizing a given Bobfile path
type SanitizeError struct {
	message string
}

// Error returns the error message for a SanitizeError.  It is expected to be
// set at the time that the struct instance is created
func (err *SanitizeError) Error() string {
	return err.message
}

// ExitCode returns the exit code for errors related to sanitizing the Bobfile
// path.  It is the same value for all sanitize errors.
func (err *SanitizeError) ExitCode() int {
	return 67
}

// ParserRelatedError is used for errors encounted during the building process
// that are related to parsing the Bobfile
type ParserRelatedError struct {
	message  string
	exitCode int
}

// Error returns the error message for a ParserRelatedError.  It is expected to
// be set at the time that the struct instance is created
func (err *ParserRelatedError) Error() string {
	return err.message
}

// ExitCode returns the exit code for errors related to parsing the Bobfile
// path.  It is expected to be set during the time that struct instance is
// created
func (err *ParserRelatedError) ExitCode() int {
	return err.exitCode
}

// BuildRelatedError is used for build-related errors produced by the builder package
// that are encountered during the build process
type BuildRelatedError struct {
	message  string
	exitCode int
}

// Error returns the error message for a build-related error.  It is expected
// to be set at the time that the struct instance is created
func (err *BuildRelatedError) Error() string {
	return err.message
}

// ExitCode returns the exit code for errors related to the build process.  It
// is expected to be set during the time that struct instance is created
func (err *BuildRelatedError) ExitCode() int {
	return err.exitCode
}
