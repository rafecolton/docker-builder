package bob

import (
//"github.com/rafecolton/bob/parser"
)

/*
Builder is responsible for taking a Builderfile struct and knowing what to do
with it to build docker containers.
*/
type Builder interface {
	Build(placeHolderArg string) error
}

/*
NewBuilder returns an instance of a Builder struct.  The function exists in
case we want to initialize our Builders with something.
*/
func NewBuilder(shouldBeRegular bool) Builder {
	if !shouldBeRegular {
		return &nullBob{}
	}
	return &regularBob{}
}

/*
Build is currently a placeholder function but will eventually have a fixed
output and be used for testing
*/
func (nullbob *nullBob) Build(placeHolderArg string) error {
	return nil
}

type nullBob struct{}

/*
Build is currently a placeholder function but will eventually be used to do the
actual work of building.
*/
func (bob *regularBob) Build(placeHolderArg string) error {
	return nil
}

type regularBob struct{}
