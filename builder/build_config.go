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

func (b *BuildConfig) File() string {
	return b.file
}

func (b *BuildConfig) Top() string {
	return b.top
}
