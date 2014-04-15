package parser

import (
	"fmt"
	"io/ioutil"
	"os"
)

import (
	"github.com/BurntSushi/toml"
	"github.com/wsxiaoys/terminal/color"
)

import (
	"github.com/rafecolton/bob/builderfile"
	"github.com/rafecolton/bob/log"
)

/*
Parser is a struct that contains a Builderfile and knows how to parse it both
as raw text and to convert toml to a Builderfile struct.  It also knows how to
tell if the Builderfile is valid (openable) or nat.
*/
type Parser struct {
	filename string
	log.Log
}

/*
NewParser returns an initialized Parser.  Not currently necessary, as no
default values are assigned to a new Parser, but useful to have in case we need
to change this.
*/
func NewParser(filename string, logger log.Log) *Parser {
	return &Parser{
		Log:      logger,
		filename: filename,
	}
}

var (
	err error
)

/*
Lint parses a builderfile and returns either nil if the file was parsed
successfully or an error indicating that parsing failed and the file is
invalid.
*/
func (parser *Parser) Lint() error {
	_, err := parser.Parse(false)

	if err != nil {
		parser.lintFailAction(false)
		return err
	}

	parser.Println(color.Sprintf("@{g!}Hooray@{|}, %s is a valid Builderfile!", parser.filename))
	return nil
}

func (parser *Parser) lintFailAction(assert bool) {
	parser.Println(color.Sprintf("@{r!}Alas@{|}, %s is not a valid Builderfile!\n----> %+v", parser.filename, err))
	if assert {
		os.Exit(5)
	}

}

/*
AssertLint is like Lint except that instead of returning an nil/error to
indicate success/failure, it exits nonzero if linting fails.
*/
func (parser *Parser) AssertLint() {
	err := parser.Lint()
	if err != nil {
		parser.lintFailAction(true)
	}
}

/*
Parse uses ParseRaw to retrieve the data from the Builderfile and to parse the
toml into a Builderfile struct instance.
*/
func (parser *Parser) Parse(assert bool) (*builderfile.Builderfile, error) {
	file := &builderfile.Builderfile{}

	raw, err := parser.ParseRaw()
	if err != nil {
		parser.lintFailAction(assert)
		return nil, err
	}

	if _, err := toml.Decode(raw, &file); err != nil {
		parser.lintFailAction(assert)
		return nil, err
	}

	return file, nil
}

/*
ParseRaw opens the Builderfile, extracts the data, and returns it as a string.
*/
func (parser *Parser) ParseRaw() (string, error) {

	if !parser.IsOpenable() {
		return "", fmt.Errorf("%s is not openable", parser.filename)
	}

	bytes, err := ioutil.ReadFile(parser.filename)

	if err != nil {
		return "", err
	}

	raw := string(bytes)

	return raw, nil
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
