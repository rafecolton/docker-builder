package parser

import (
	"errors"
)

import (
	"github.com/onsi/gocleanup"
	"github.com/wsxiaoys/terminal/color"
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

// helper functions
func (parser *Parser) printLintSuccessMessage() {
	parser.Println(color.Sprintf("@{g!}Hooray@{|}, %s is a valid Builderfile!", parser.filename))
}

func (parser *Parser) printLintFailMessage(err error) {
	parser.Println(
		color.Sprintf(
			"@{r!}Alas@{|}, %s is not a valid Builderfile\n----> %+v", parser.filename, err,
		),
	)
}
