package parser

import (
	"github.com/modcloth/docker-builder/builderfile"
	"io/ioutil"
)

import (
	"github.com/BurntSushi/toml"
)

// Step 1 of Parse
func (parser *Parser) getRaw() (string, Error) {
	bytes, err := ioutil.ReadFile(parser.filename)
	if err != nil {
		return "", &OSPathError{error: err}
	}

	raw := string(bytes)

	return raw, nil
}

// Step 2 of Parse
func (parser *Parser) rawToStruct() (*builderfile.Builderfile, Error) {
	raw, err := parser.getRaw()
	if err != nil {
		return nil, err
	}

	file := &builderfile.Builderfile{}
	if _, err := toml.Decode(raw, &file); err != nil {
		return nil, &TOMLParseError{error: err}
	}

	file.Clean()

	return file, nil
}

// Step 2.5 of Parse - handle Builderfile version

func (parser *Parser) convertBuilderfileVersion() (*builderfile.Builderfile, Error) {
	//var fileZero *builderfile.Builderfile

	fileZero, err := parser.rawToStruct()
	if err != nil {
		return nil, err
	}

	// check version, do conversion
	if fileZero.Version == 1 {
		return fileZero, nil
	}

	builderfile.Logger(parser.Logger)
	fileOne, goErr := builderfile.Convert0to1(fileZero)
	if goErr != nil {
		return nil, &BuilderfileConvertError{error: goErr}
	}

	return fileOne, nil
}

// Step 3 of Parse
func (parser *Parser) structToInstructionSet() (*InstructionSet, Error) {
	file, err := parser.convertBuilderfileVersion()
	if err != nil {
		return nil, err
	}

	return parser.instructionSetFromBuilderfileStruct(file), nil
}

// Step 4 of Parse()
func (parser *Parser) instructionSetToCommandSequence() (*CommandSequence, Error) {
	is, err := parser.structToInstructionSet()
	if err != nil {
		return nil, err
	}

	return parser.commandSequenceFromInstructionSet(is), nil
}

// wrapper function for the final step
func (parser *Parser) finalStep() (interface{}, Error) {
	return parser.instructionSetToCommandSequence()
}

/*
Parse further parses the Builderfile struct into an InstructionSet struct,
merging the global container options into the individual container sections.
*/
func (parser *Parser) Parse() (*CommandSequence, Error) {
	ret, err := parser.finalStep()
	if err != nil {
		return nil, err
	}

	return ret.(*CommandSequence), nil
}
