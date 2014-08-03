package builder

import (
	"path/filepath"
	"regexp"
)

const (
	DotDotSanitizeErrorMessage      = "bobfile path must not contain .."
	InvalidPathSanitizeErrorMessage = "bobfile path is invalid"
	SymlinkSanitizeErrorMessage     = "bobfile path must not contain symlinks"
)

var dotDotRegex = regexp.MustCompile("\\.\\.")

func SanitizeBuilderfilePath(file string) (string, BuilderError) {
	if dotDotRegex.MatchString(file) {
		return "", &SanitizeError{message: DotDotSanitizeErrorMessage}
	}

	abs, err := filepath.Abs("./" + file)
	if err != nil {
		return "", &SanitizeError{message: InvalidPathSanitizeErrorMessage}
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", &SanitizeError{message: InvalidPathSanitizeErrorMessage}
	}

	if abs != resolved {
		return "", &SanitizeError{message: SymlinkSanitizeErrorMessage}
	}

	clean := filepath.Clean(abs)

	return clean, nil
}
