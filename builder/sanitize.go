package builder

import (
	"path/filepath"
	"regexp"
)

const (
	// DotDotSanitizeErrorMessage is the error message used in errors that occur
	// because a provided Bobfile path contains ".."
	DotDotSanitizeErrorMessage = "file path must not contain .."

	// InvalidPathSanitizeErrorMessage is the error message used in errors that
	// occur because a provided Bobfile path is invalid
	InvalidPathSanitizeErrorMessage = "file path is invalid"

	// SymlinkSanitizeErrorMessage is the error message used in errors that
	// occur because a provided Bobfile path contains symlinks
	SymlinkSanitizeErrorMessage = "file path must not contain symlinks"
)

var dotDotRegex = regexp.MustCompile(`\.\.`)

// SanitizeTrustedFilePath checks for disallowed entries in the provided
// file path and returns either a sanitized version of the path or an error
func SanitizeTrustedFilePath(trustedFilePath *TrustedFilePath) (string, Error) {
	var file = trustedFilePath.File()
	var top = trustedFilePath.Top()

	if dotDotRegex.MatchString(file) {
		return "", &SanitizeError{
			Message:  DotDotSanitizeErrorMessage,
			Filename: file,
		}
	}

	abs, err := filepath.Abs(top + "/" + file)
	if err != nil {
		return "", &SanitizeError{
			Message:  InvalidPathSanitizeErrorMessage,
			error:    err,
			Filename: file,
		}
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", &SanitizeError{
			Message:  InvalidPathSanitizeErrorMessage,
			error:    err,
			Filename: file,
		}
	}

	if abs != resolved {
		return "", &SanitizeError{
			Message:  SymlinkSanitizeErrorMessage,
			Filename: file,
		}
	}

	clean := filepath.Clean(abs)

	return clean, nil
}
