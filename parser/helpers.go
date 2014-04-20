package parser

import (
	"github.com/wsxiaoys/terminal/color"
	"os"
)

// helper functions
func (parser *Parser) printLintSuccessMessage() {
	parser.Println(color.Sprintf("@{g!}Hooray@{|}, %s is a valid Builderfile!", parser.filename))
}

func (parser *Parser) printLintFailMessage(err error) {
	parser.Println(color.Sprintf("@{r!}Alas@{|}, %s is not a valid Builderfile\n----> %+v", parser.filename, err))
}

/*
IsOpenable examines the Builderfile provided to the Parser and returns a bool
indicating whether or not the file exists and openable.
*/
func (parser *Parser) IsOpenable() bool {

	file, err := os.Open(parser.filename)
	defer file.Close()

	if err != nil {
		return false
	}

	return true
}
