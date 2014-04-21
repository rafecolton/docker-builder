package bob

import (
	"github.com/rafecolton/bob/parser"
)

/*
Builder is responsible for taking a Builderfile struct and knowing what to do
with it to build docker containers.
*/
type Builder interface {
	Build(commands *parser.CommandSequence) error
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
func (nullbob *nullBob) Build(commands *parser.CommandSequence) error {
	return nil
}

type nullBob struct{}

/*
Build is currently a placeholder function but will eventually be used to do the
actual work of building.
*/
func (bob *regularBob) Build(commands *parser.CommandSequence) error {
	/*
		  TODO:
		  - inject setup and teardown commands
		  - integrate with gocleanup
		  - take docker stuff out of parser and put here
		  - setup/teardown process:
			1. create tmp dir in work dir
			2. if include is empty, start with all, otherwise start with include
				2a. remove excludes
			3. copy results into tmpdir
			4. copy dockerfile into tmpdir as 'Dockerfile'
			5. build
			6. tag
			7. push
	*/
	return nil
}

type regularBob struct{}
