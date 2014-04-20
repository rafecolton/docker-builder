package parser

import "os/exec"

import "github.com/rafecolton/bob/builderfile"

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
A CommandSequence is an intermediate data type in the parsing process. Once a
Builderfile is parsed into an InstructionSet, it is further parsed into a
CommandSequence, which is essential an array of strings where each string is a
command to be run.
*/
type CommandSequence struct {
	commands []exec.Cmd
}
