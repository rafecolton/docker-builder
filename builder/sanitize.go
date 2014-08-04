package builder

import (
	"path/filepath"
	"regexp"
)

const (
	// DotDotSanitizeErrorMessage is the error message used in errors that occur
	// because a provided Bobfile path contains ".."
	DotDotSanitizeErrorMessage = "bobfile path must not contain .."

	// InvalidPathSanitizeErrorMessage is the error message used in errors that
	// occur because a provided Bobfile path is invalid
	InvalidPathSanitizeErrorMessage = "bobfile path is invalid"

	// SymlinkSanitizeErrorMessage is the error message used in errors that
	// occur because a provided Bobfile path contains symlinks
	SymlinkSanitizeErrorMessage = "bobfile path must not contain symlinks"
)

var dotDotRegex = regexp.MustCompile("\\.\\.")

// SanitizeBuilderfilePath checks for disallowed entries in the provided
// Bobfile path and returns either a sanitized version of the path or an error
func SanitizeBuilderfilePath(config *BuildConfig) (string, Error) {
	var file = config.File()
	var top = config.Top()

	if dotDotRegex.MatchString(file) {
		return "", &SanitizeError{Message: DotDotSanitizeErrorMessage}
	}

	abs, err := filepath.Abs(top + "/" + file)
	if err != nil {
		return "", &SanitizeError{Message: InvalidPathSanitizeErrorMessage}
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", &SanitizeError{Message: InvalidPathSanitizeErrorMessage}
	}

	if abs != resolved {
		return "", &SanitizeError{Message: SymlinkSanitizeErrorMessage}
	}

	clean := filepath.Clean(abs)

	return clean, nil
}
