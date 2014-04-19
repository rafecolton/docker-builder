package parser

import (
	"fmt"
	"github.com/rafecolton/bob/builderfile"
	"io/ioutil"
)

import (
	"github.com/BurntSushi/toml"
)

// Step 1 of Parse
func (parser *Parser) getRaw() (string, error) {

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

// Step 2 of Parse
func (parser *Parser) rawToStruct() (*builderfile.Builderfile, error) {
	file := &builderfile.Builderfile{}

	raw, err := parser.getRaw()
	if err != nil {
		parser.printLintFailMessage(err)
		return nil, err
	}

	if _, err := toml.Decode(raw, &file); err != nil {
		parser.printLintFailMessage(err)
		return nil, err
	}

	return file, nil
}

/*
Parse further parses the Builderfile struct into an InstructionSet struct,
merging the global container options into the individual container sections.
*/
func (parser *Parser) Parse() (*InstructionSet, error) {
	// TODO:
	return nil, nil
}
