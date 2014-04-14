package main

import (
	"github.com/rafecolton/builder/builderfile"
)

/*
Builder is responsible for taking a Builderfile struct and knowing what to do
with it to build docker containers.
*/
type Builder struct {
	file *builderfile.Builderfile
}

/*
NewBuilder returns an instance of a Builder struct.  The function exists in
case we want to initialize our Builders with something.
*/
func NewBuilder() *Builder {
	return &Builder{}
}

/*
Build performs the actual building work of the Builder.
*/
func (builder *Builder) Build(file *builderfile.Builderfile) error {
	builder.file = file

	return nil
}