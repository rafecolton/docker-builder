package linter

import (
	//"fmt"
	"github.com/rafecolton/builder/config"
)

type Linter struct {
	*config.Runtime
}

// like init
func New(runtime *config.Runtime) *Linter {
	return &Linter{
		runtime,
	}
}

func (me *Linter) Lint() string {
	return "parsing \"" + me.Options.Lintfile + "\""
}
