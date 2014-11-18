package parser

import (
	"github.com/fsouza/go-dockerclient"

	"github.com/sylphon/build-runner/builderfile"
	"github.com/sylphon/build-runner/parser/tag"
)

// CommandSequenceFromInstructionSet turns an InstructionSet struct into a
// CommandSequence struct - one of the intermediate steps to building, will
// eventually be made private
func (parser *Parser) CommandSequenceFromInstructionSet(is *InstructionSet) *CommandSequence {
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

		// ADD BUILD COMMAND
		uuid, err := parser.NextUUID()
		if err != nil {
			return nil
		}

		name := v.Registry + "/" + v.Project
		initialTag := name + ":" + uuid

		// get docker registry credentials
		un := v.CfgUn
		pass := v.CfgPass
		email := v.CfgEmail

		buildOpts := docker.BuildImageOptions{
			Name:           initialTag,
			RmTmpContainer: true,
			ContextDir:     parser.top,
			Auth: docker.AuthConfiguration{
				Username: un,
				Password: pass,
				Email:    email,
			},
			AuthConfigs: docker.AuthConfigurations{
				Configs: map[string]docker.AuthConfiguration{
					v.Registry: docker.AuthConfiguration{
						Username:      un,
						Password:      pass,
						Email:         email,
						ServerAddress: v.Registry,
					},
				},
			},
		}

		for _, opt := range is.DockerBuildOpts {
			switch opt {
			case "--force-rm":
				buildOpts.ForceRmTmpContainer = true
			case "--no-cache":
				buildOpts.NoCache = true
			case "-q", "--quiet":
				buildOpts.SuppressOutput = true
			case "--no-rm":
				// Is "--no-rm" this the best way to handle this since default is true?
				// Maybe so, just document it somewhere (TODO)
				buildOpts.RmTmpContainer = false
			}
		}

		containerCommands = append(containerCommands, &BuildCmd{
			buildOpts:     buildOpts,
			origBuildOpts: is.DockerBuildOpts,
		})

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

// InstructionSetFromBuilderfileStruct turns a UnitConfig struct into an
// InstructionSet struct - one of the intermediate steps to building, will
// eventually be made private
func (parser *Parser) InstructionSetFromBuilderfileStruct(file *builderfile.UnitConfig) *InstructionSet {
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
