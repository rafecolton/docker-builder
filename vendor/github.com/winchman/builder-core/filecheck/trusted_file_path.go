package filecheck

import (
	"path/filepath"
)

// TrustedFilePathState represents the state of a TrustedFilePath
type TrustedFilePathState int

const (
	// Unchecked means the TrustedFilePath has not been checked for sanitization
	Unchecked TrustedFilePathState = iota

	// Errored means there was an error in creation
	Errored

	// OK means the path has been check and is sanitary
	OK

	// NotOK means the path has been check and is NOT sanitary
	NotOK
)

// TrustedFilePath contains the fields of a path to a file including "top", the
// top level directory containing the file, and "file", the relative path from
// top to the file. Top, if not provided will default to "."
type TrustedFilePath struct {
	file  string
	top   string
	State TrustedFilePathState
	Error error
}

// NewTrustedFilePathOptions are options for creating a new trusted file path -
// used for disambiguating the args
type NewTrustedFilePathOptions struct {
	File string
	Top  string
}

// NewTrustedFilePath returns struct representation of the path to a file.  If
// top is not provided (an empty string), it is defaulted to ".".  The top is
// sanitized by evaluating all symlinks so that it can be properly sanitized by
// the builder when the time comes.  This treats "top" as trusted (i.e.
// relative paths and symlinks can be evaluated safely).
func NewTrustedFilePath(opts NewTrustedFilePathOptions) (*TrustedFilePath, error) {
	top := opts.Top
	file := opts.File
	if top == "" {
		top = "."
	}
	abs, err := filepath.Abs(top)
	if err != nil {
		return &TrustedFilePath{
			State: Errored,
			Error: err,
		}, err
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return &TrustedFilePath{
			State: Errored,
			Error: err,
		}, err
	}

	return &TrustedFilePath{
		file:  file,
		top:   resolved,
		State: Unchecked,
	}, nil
}

// File returns the file component of the trusted file path
func (b *TrustedFilePath) File() string {
	return b.file
}

// Top returns the top level directory containing the file.  The path to the
// top level directory as provided by Top() is considered to be trusted.
func (b *TrustedFilePath) Top() string {
	return b.top
}

// FullPath returns the full path to the trusted file (i.e. top + "/" + file)
func (b *TrustedFilePath) FullPath() string {
	return b.top + "/" + b.file
}
