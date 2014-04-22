package parser

import (
	"fmt"
	"os/exec"
)

import (
	"github.com/rafecolton/bob/builderfile"
	"github.com/rafecolton/bob/parser/tag"
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
	Containers      map[string]builderfile.ContainerSection
}

/*
A CommandSequence is an intermediate data type in the parsing process. Once a
Builderfile is parsed into an InstructionSet, it is further parsed into a
CommandSequence, which is essential an array of strings where each string is a
command to be run.
*/
type CommandSequence struct {
	commands []*SubSequence
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
}

/*
A SubSequence is a logical grouping of commands such as a sequence of build,
tag, and push commands.  In addition, the subsequence metadata contains any
important metadata about the container build such as the name of the Dockerfile
and which files/dirs to exclude.
*/
type SubSequence struct {
	Metadata   *SubSequenceMetadata
	SubCommand []exec.Cmd
}

// turns InstructionSet structs into CommandSequence structs
func (parser *Parser) commandSequenceFromInstructionSet(is *InstructionSet) *CommandSequence {
	ret := &CommandSequence{
		commands: []*SubSequence{},
	}

	var containerCommands []exec.Cmd

	for _, v := range is.Containers {
		containerCommands = []exec.Cmd{}

		// ADD BUILD COMMANDS
		uuid, err := parser.NextUUID()
		if err != nil {
			return nil
		}

		name := fmt.Sprintf("%s/%s", v.Registry, v.Project)
		initialTag := fmt.Sprintf("%s:%s", name, uuid)
		buildArgs := []string{"docker", "build", "-t", initialTag}
		buildArgs = append(buildArgs, is.DockerBuildOpts...)
		buildArgs = append(buildArgs, ".")

		containerCommands = append(containerCommands, *&exec.Cmd{
			Path: "docker",
			Args: buildArgs,
		})

		// ADD TAG COMMANDS
		for _, t := range v.Tags {
			var tagObj tag.Tag
			tagArg := map[string]string{"tag": t}

			if len(t) > 4 && t[0:4] == "git:" {
				tagObj = tag.NewTag("git", tagArg)
			} else {
				tagObj = tag.NewTag("default", tagArg)
			}

			fullTag := fmt.Sprintf("%s:%s", name, tagObj.Tag())
			buildArgs = []string{"docker", "tag"}
			buildArgs = append(buildArgs, is.DockerTagOpts...)
			buildArgs = append(buildArgs, "<IMG>", fullTag)

			containerCommands = append(containerCommands, *&exec.Cmd{
				Path: "docker",
				Args: buildArgs,
			})
		}

		// ADD PUSH COMMANDS
		buildArgs = []string{"docker", "push", name}
		containerCommands = append(containerCommands, *&exec.Cmd{
			Path: "docker",
			Args: buildArgs,
		})

		ret.commands = append(ret.commands, &SubSequence{
			Metadata: &SubSequenceMetadata{
				Name:       v.Name,
				Dockerfile: v.Dockerfile,
				Included:   v.Included,
				Excluded:   v.Excluded,
			},
			SubCommand: containerCommands,
		})
	}

	return ret
}

// turns Builderfile structs into InstructionSet structs
func (parser *Parser) instructionSetFromBuilderfileStruct(file *builderfile.Builderfile) *InstructionSet {
	ret := &InstructionSet{
		DockerBuildOpts: file.Docker.BuildOpts,
		DockerTagOpts:   file.Docker.TagOpts,
		Containers:      *&map[string]builderfile.ContainerSection{},
	}

	if file.Containers == nil {
		file.Containers = map[string]builderfile.ContainerSection{}
	} else {
		globals, hasGlobals := file.Containers["global"]

		for k, v := range file.Containers {
			if k == "global" {
				continue
			}

			dockerfile := v.Dockerfile
			included := v.Included
			excluded := v.Excluded
			registry := v.Registry
			project := v.Project
			tags := v.Tags

			if hasGlobals {
				if dockerfile == "" {
					dockerfile = globals.Dockerfile
				}

				if registry == "" {
					registry = globals.Registry
				}

				if project == "" {
					project = globals.Project
				}

				if included == nil || len(included) == 0 {
					if globals.Included == nil {
						included = []string{}
					} else {
						included = globals.Included
					}
				}

				if excluded == nil || len(excluded) == 0 {
					if globals.Excluded == nil {
						excluded = []string{}
					} else {
						excluded = globals.Excluded
					}
				}

				if tags == nil || len(tags) == 0 {
					if globals.Tags == nil {
						tags = []string{}
					} else {
						tags = globals.Tags
					}
				}
			}

			containerSection := &builderfile.ContainerSection{
				Name:       k,
				Dockerfile: dockerfile,
				Included:   included,
				Excluded:   excluded,
				Registry:   registry,
				Project:    project,
				Tags:       tags,
			}

			ret.Containers[k] = *containerSection
		}
	}

	return ret
}
