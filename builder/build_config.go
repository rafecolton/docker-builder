package builder

import (
	"path/filepath"
)

// BuildConfig contains fields that are passed to the builder to define
// parameters of the build.  File must be the Bobfile.  Top, if not provided,
// will default to "."
type BuildConfig struct {
	file string
	top  string
}

// NewBuildConfig returns a *BuildConfig after doing a little bit of setup.  If
// top is not provided (an empty string), it is defaulted to ".".  The top is
// sanitized by evaluating all symlinks so that it can be properly sanitized by
// the builder when the time comes
func NewBuildConfig(file, top string) (*BuildConfig, error) {
	//ret := &BuildConfig{}
	//ret.file = file
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

	return &BuildConfig{
		file: file,
		top:  resolved,
	}, nil
}

// File returns the Bobfile associated with the build config
func (b *BuildConfig) File() string {
	return b.file
}

// Top returns the repo directory after the setup has been done in NewBuildConfig
func (b *BuildConfig) Top() string {
	return b.top
}
