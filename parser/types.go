package parser

import (
	"github.com/modcloth/docker-builder/builderfile"
)

/*
An InstructionSet is an intermediate datatype - once a Builderfile is parsed
and the TOML is validated, the parser parses the data into an InstructionSet.
The primary purpose of this step is to merge any global container options into
the sections for the individual containers.
*/
type InstructionSet struct {
	DockerBuildOpts []string
	DockerTagOpts   []string
	Containers      []builderfile.ContainerSection
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
	Included   []string
	Excluded   []string
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
