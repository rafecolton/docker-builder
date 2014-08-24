package builder

import (
	"path/filepath"
)

// TrustedFilePath contains the fields of a path to a file including "top", the
// top level directory containing the file, and "file", the relative path from
// top to the file. Top, if not provided will default to "."
type TrustedFilePath struct {
	file string
	top  string
}

// NewTrustedFilePath returns struct representation of the path to a file.  If
// top is not provided (an empty string), it is defaulted to ".".  The top is
// sanitized by evaluating all symlinks so that it can be properly sanitized by
// the builder when the time comes.  This treats "top" as trusted (i.e.
// relative paths and symlinks can be evaluated safely).
func NewTrustedFilePath(file, top string) (*TrustedFilePath, error) {
	if top == "" {
		top = "."
	}
	abs, err := filepath.Abs(top)
	if err != nil {
		return nil, err
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return nil, err
	}

	return &TrustedFilePath{
		file: file,
		top:  resolved,
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

// Dir returns the dirname of the full trusted file path
func (b *TrustedFilePath) Dir() string {
	return filepath.Dir(b.top + "/" + b.file)
}

// Dir returns the basename of the full trusted file path
func (b *TrustedFilePath) Base() string {
	return filepath.Base(b.top + "/" + b.file)
}
