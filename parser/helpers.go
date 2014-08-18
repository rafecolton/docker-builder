package parser

import (
	"os/exec"

	"github.com/rafecolton/docker-builder/builderfile"
	"github.com/rafecolton/docker-builder/conf"
	"github.com/rafecolton/docker-builder/parser/tag"
)

// turns InstructionSet structs into CommandSequence structs
func (parser *Parser) commandSequenceFromInstructionSet(is *InstructionSet) *CommandSequence {
	ret := &CommandSequence{
		Commands: []*SubSequence{},
	}

	var containerCommands []DockerCmd
	var tagCommands []DockerCmd
	var pushCommands []DockerCmd

	for _, v := range is.Containers {
		containerCommands = []DockerCmd{}
		tagCommands = []DockerCmd{}
		pushCommands = []DockerCmd{}

		// ADD BUILD COMMANDS
		uuid, err := parser.NextUUID()
		if err != nil {
			return nil
		}

		name := v.Registry + "/" + v.Project
		initialTag := name + ":" + uuid
		buildArgs := []string{"docker", "build", "-t", initialTag}
		buildArgs = append(buildArgs, is.DockerBuildOpts...)
		buildArgs = append(buildArgs, ".")

		containerCommands = append(containerCommands, &BuildCmd{
			Cmd: &exec.Cmd{
				Path: "docker",
				Args: buildArgs,
			},
		})

		// get docker registry credentials
		un := v.CfgUn
		pass := v.CfgPass
		email := v.CfgEmail
		if un == "" {
			un = conf.Config.CfgUn
		}
		if pass == "" {
			pass = conf.Config.CfgPass
		}
		if email == "" {
			email = conf.Config.CfgEmail
		}

		// ADD TAG COMMANDS
		for _, t := range v.Tags {
			var tagObj tag.Tag
			tagArg := map[string]string{
				"tag": t,
				"top": parser.top,
			}

			if len(t) > 4 && t[0:4] == "git:" {
				tagObj = tag.NewTag("git", tagArg)
			} else {
				tagObj = tag.NewTag("default", tagArg)
			}

			tagCmd := &TagCmd{
				Repo: name,
				Tag:  tagObj.Tag(),
			}
			for _, opt := range is.DockerTagOpts {
				if opt == "-f" || opt == "--force" {
					tagCmd.Force = true
				}
			}

			tagCommands = append(tagCommands, tagCmd)

			// ADD CORRESPONDING PUSH COMMAND
			if !v.SkipPush {
				pushCmd := &PushCmd{
					Image:     name,
					Tag:       tagObj.Tag(),
					AuthUn:    un,
					AuthPwd:   pass,
					AuthEmail: email,
					Registry:  v.Registry,
				}
				pushCommands = append(pushCommands, pushCmd)
			}
		}

		containerCommands = append(containerCommands, tagCommands...)
		containerCommands = append(containerCommands, pushCommands...)

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

func mergeGlobals(container, globals *builderfile.ContainerSection) *builderfile.ContainerSection {

	if container.Tags == nil {
		container.Tags = []string{}
	}

	if container.Dockerfile == "" {
		container.Dockerfile = globals.Dockerfile
	}

	if container.Registry == "" {
		container.Registry = globals.Registry
	}

	if container.Project == "" {
		container.Project = globals.Project
	}

	if len(container.Tags) == 0 && globals.Tags != nil {
		container.Tags = globals.Tags
	}

	container.SkipPush = container.SkipPush || globals.SkipPush

	if container.CfgUn == "" {
		container.CfgUn = globals.CfgUn
	}

	if container.CfgPass == "" {
		container.CfgPass = globals.CfgPass
	}

	if container.CfgEmail == "" {
		container.CfgEmail = globals.CfgEmail
	}

	return container
}

// turns Builderfile structs into InstructionSet structs
func (parser *Parser) instructionSetFromBuilderfileStruct(file *builderfile.Builderfile) *InstructionSet {
	ret := &InstructionSet{
		DockerBuildOpts: file.Docker.BuildOpts,
		DockerTagOpts:   file.Docker.TagOpts,
		Containers:      []builderfile.ContainerSection{},
	}

	if file.ContainerArr == nil {
		file.ContainerArr = []*builderfile.ContainerSection{}
	}

	if file.ContainerGlobals == nil {
		file.ContainerGlobals = &builderfile.ContainerSection{}
	}
	globals := file.ContainerGlobals

	for _, container := range file.ContainerArr {
		container = mergeGlobals(container, globals)
		ret.Containers = append(ret.Containers, *container)
	}

	return ret
}
