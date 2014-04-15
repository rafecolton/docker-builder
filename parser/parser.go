package parser

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/rafecolton/bob/builderfile"
	"io/ioutil"
	"os"
)

/*
Parser is a struct that contains a Builderfile and knows how to parse it both
as raw text and to convert toml to a Builderfile struct.  It also knows how to
tell if the Builderfile is valid (openable) or nat.
*/
type Parser struct {
	Builderfile string
}

/*
NewParser returns an initialized Parser.  Not currently necessary, as no
default values are assigned to a new Parser, but useful to have in case we need
to change this.
*/
func NewParser() *Parser {
	return &Parser{}
}

var (
	err error
)

/*
Parse uses ParseRaw to retrieve the data from the Builderfile and to parse the
toml into a Builderfile struct instance.
*/
func (parser *Parser) Parse() (*builderfile.Builderfile, error) {
	bf := &builderfile.Builderfile{}

	raw, err := parser.ParseRaw()
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(raw, &bf); err != nil {
		return nil, err
	}

	return bf, nil
}

/*
ParseRaw opens the Builderfile, extracts the data, and returns it as a string.
*/
func (parser *Parser) ParseRaw() (string, error) {

	if !parser.IsOpenable() {
		return "", fmt.Errorf("%s is not openable", parser.Builderfile)
	}

	bytes, err := ioutil.ReadFile(parser.Builderfile)

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

	file, err := os.Open(parser.Builderfile)
	defer file.Close()

	if err != nil {
		return false
	}

	return true
}
