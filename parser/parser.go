package parser

import (
	"errors"
	"os"
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
An InstructionSet is an intermediate datatype - once a Builderfile is parsed
and the TOML is validated, the parser parses the data into an InstructionSet.
The primary purpose of this step is to merge any global container options into
the sections for the individual containers.
*/
type InstructionSet struct {
	DockerBuildOpts []string
	DockerTagOpts   []string
	Containers      map[string]builderfile.ContainerSection
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
	//TODO: this should probably call Parse instead of the unexported method rawToStruct
	_, err := parser.rawToStruct()

	if err != nil {
		parser.printLintFailMessage(err)
		return err
	}

	parser.printLintSuccessMessage()
	return nil
}

/*
AssertLint is like Lint except that instead of returning an nil/error to
indicate success/failure, it exits nonzero if linting fails.
*/
func (parser *Parser) AssertLint() {
	if !parser.IsOpenable() {
		parser.printLintFailMessage(errors.New("unable to open file"))
		os.Exit(17)
	}

	err := parser.Lint()
	if err != nil {
		parser.printLintFailMessage(err)
		os.Exit(5)
	}
}
