package parser

import (
	"fmt"
	"os"
	"os/exec"
)

import (
	"github.com/rafecolton/bob/builderfile"
	"github.com/rafecolton/bob/parser/tag"
)

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

// turns InstructionSet structs into CommandSequence structs
func (parser *Parser) commandSequenceFromInstructionSet(is *InstructionSet) *CommandSequence {
	ret := &CommandSequence{
		Commands: []*SubSequence{},
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
		if !v.SkipPush {
			buildArgs = []string{"docker", "push", name}
			containerCommands = append(containerCommands, *&exec.Cmd{
				Path: "docker",
				Args: buildArgs,
			})
		}

		ret.Commands = append(ret.Commands, &SubSequence{
			Metadata: &SubSequenceMetadata{
				Name:       v.Name,
				Dockerfile: v.Dockerfile,
				Included:   v.Included,
				Excluded:   v.Excluded,
				UUID:       uuid,
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
			skipPush := v.SkipPush

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
				SkipPush:   skipPush,
			}

			ret.Containers[k] = *containerSection
		}
	}

	return ret
}
