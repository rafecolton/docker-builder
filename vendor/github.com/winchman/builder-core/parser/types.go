package parser

import (
	"github.com/winchman/builder-core/communication"
	"github.com/winchman/builder-core/unit-config"
)

/*
Parser is a struct that contains a Builderfile and knows how to parse it both
as raw text and to convert toml to a Builderfile struct.  It also knows how to
tell if the Builderfile is valid (openable) or nat.
*/
type Parser struct {
	contextDir string
	reporter   *comm.Reporter
}

// NewParserOptions encapsulates all of the options necessary when creating a new parser
type NewParserOptions struct {
	ContextDir string
	Log        comm.LogChan
	Event      comm.EventChan
}

/*
NewParser returns an initialized Parser.  Not currently necessary, as no
default values are assigned to a new Parser, but useful to have in case we need
to change this.
*/
func NewParser(opts NewParserOptions) *Parser {
	return &Parser{
		contextDir: opts.ContextDir,
		reporter:   comm.NewReporter(opts.Log, opts.Event),
	}
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
	Containers      []unitconfig.ContainerSection
}

/*
A CommandSequence is an intermediate data type in the parsing process. Once a
Builderfile is parsed into an InstructionSet, it is further parsed into a
CommandSequence, which is essential an array of strings where each string is a
command to be run.
*/
type CommandSequence struct {
	Commands []*SubSequence
}

/*
SubSequenceMetadata contains any important metadata about the container build
such as the name of the Dockerfile and which files/dirs to exclude.
*/
type SubSequenceMetadata struct {
	Name       string
	Dockerfile string
	UUID       string
	SkipPush   bool
}

/*
A SubSequence is a logical grouping of commands such as a sequence of build,
tag, and push commands.  In addition, the subsequence metadata contains any
important metadata about the container build such as the name of the Dockerfile
and which files/dirs to exclude.
*/
type SubSequence struct {
	Metadata   *SubSequenceMetadata
	SubCommand []DockerCmd
}
