package parser

import (
	"errors"
)

import (
	"github.com/onsi/gocleanup"
)

/*
Lint parses a builderfile and returns either nil if the file was parsed
successfully or an error indicating that parsing failed and the file is
invalid.
*/
func (parser *Parser) Lint() error {
	//TODO: this should probably call Parse instead of the unexported method rawToStruct
	_, err := parser.rawToStruct()

	if err != nil {
		return err
	}

	return nil
}

/*
AssertLint is like Lint except that instead of returning an nil/error to
indicate success/failure, it exits nonzero if linting fails.
*/
func (parser *Parser) AssertLint() {
	if !parser.IsOpenable() {
		parser.printLintFailMessage(errors.New("unable to open file"))
		gocleanup.Exit(17)
	}

	err := parser.Lint()
	if err != nil {
		parser.printLintFailMessage(err)
		gocleanup.Exit(5)
	} else {
		parser.printLintSuccessMessage()
		gocleanup.Exit(0)
	}
}
