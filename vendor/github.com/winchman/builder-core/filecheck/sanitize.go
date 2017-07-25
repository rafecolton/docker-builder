package filecheck

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

const (
	// DotDotSanitizeErrorMessage is the error message used in errors that occur
	// because a provided Bobfile path contains ".."
	dotDotSanitizeErrorMessage = "file path must not contain .."

	// InvalidPathSanitizeErrorMessage is the error message used in errors that
	// occur because a provided Bobfile path is invalid
	invalidPathSanitizeErrorMessage = "file path is invalid"

	// SymlinkSanitizeErrorMessage is the error message used in errors that
	// occur because a provided Bobfile path contains symlinks
	symlinkSanitizeErrorMessage = "file path must not contain symlinks"

	// DoesNotExistSanitizeErrorMessage is the error message used in cases
	// where the error results in the requested file not existing
	doesNotExistSanitizeErrorMessage = "file requested does not exist"
)

var dotDotRegex = regexp.MustCompile(`\.\.`)

// Sanitize checks for disallowed entries in the provided file path and sets
// the State and Error values of the trustedFilePath
func (trustedFilePath *TrustedFilePath) Sanitize() {
	var file = trustedFilePath.File()
	var top = trustedFilePath.Top()

	if dotDotRegex.MatchString(file) {
		trustedFilePath.State = NotOK
		trustedFilePath.Error = errors.New(dotDotSanitizeErrorMessage)
		return
	}

	abs, err := filepath.Abs(top + "/" + file)
	if err != nil {
		trustedFilePath.State = NotOK
		trustedFilePath.Error = errors.New(invalidPathSanitizeErrorMessage)
		return
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		trustedFilePath.State = NotOK
		trustedFilePath.Error = errors.New(invalidPathSanitizeErrorMessage)
		if os.IsNotExist(err) {
			trustedFilePath.Error = errors.New(doesNotExistSanitizeErrorMessage)
		}
		return
	}

	if abs != resolved {
		trustedFilePath.State = NotOK
		trustedFilePath.Error = errors.New(symlinkSanitizeErrorMessage)
		return
	}

	clean := filepath.Clean(abs)
	trustedFilePath.top = filepath.Dir(clean)
	trustedFilePath.file = filepath.Base(clean)
	trustedFilePath.State = OK
}
