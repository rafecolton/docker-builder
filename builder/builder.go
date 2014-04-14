package main

import (
	"github.com/rafecolton/builder/builderfile"
)

type Builder struct {
	file *builderfile.Builderfile
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (me *Builder) Build(file *builderfile.Builderfile) error {
	me.file = file

	return nil
}
