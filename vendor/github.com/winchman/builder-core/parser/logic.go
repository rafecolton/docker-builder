package parser

import (
	"github.com/fsouza/go-dockerclient"
	gouuid "github.com/nu7hatch/gouuid"

	"github.com/winchman/builder-core/unit-config"
)

// Parse - does the parsing!
func (parser *Parser) Parse(file *unitconfig.UnitConfig) *CommandSequence {
	return parser.step2(parser.step1(file))
}

// step1 (formerly InstructionSetFromBuilderfileStruct) turns a UnitConfig
// struct into an InstructionSet struct - one of the intermediate steps to
// building, will eventually be made private
func (parser *Parser) step1(file *unitconfig.UnitConfig) *InstructionSet {
	ret := &InstructionSet{
		DockerBuildOpts: file.Docker.BuildOpts,
		DockerTagOpts:   file.Docker.TagOpts,
		Containers:      []unitconfig.ContainerSection{},
	}

	if file.ContainerArr == nil {
		file.ContainerArr = []*unitconfig.ContainerSection{}
	}

	if file.ContainerGlobals == nil {
		file.ContainerGlobals = &unitconfig.ContainerSection{}
	}
	globals := file.ContainerGlobals

	for _, container := range file.ContainerArr {
		container = mergeGlobals(container, globals)
		ret.Containers = append(ret.Containers, *container)
	}

	return ret
}

func mergeGlobals(container, globals *unitconfig.ContainerSection) *unitconfig.ContainerSection {
	if container.Dockerfile == "" {
		container.Dockerfile = globals.Dockerfile
	}
	if container.Registry == "" {
		container.Registry = globals.Registry
	}
	if container.Project == "" {
		container.Project = globals.Project
	}
	if container.CfgUn == "" {
		container.CfgUn = globals.CfgUn
	}
	if container.CfgPass == "" {
		container.CfgPass = globals.CfgPass
	}
	if container.CfgEmail == "" {
		container.CfgEmail = globals.CfgEmail
	}

	if container.Tags == nil {
		container.Tags = []string{}
	}
	if len(container.Tags) == 0 && globals.Tags != nil {
		container.Tags = globals.Tags
	}

	container.SkipPush = container.SkipPush || globals.SkipPush

	return container
}

// step2 (formerly CommandSequenceFromInstructionSet) turns an InstructionSet struct into a
// CommandSequence struct - one of the intermediate steps to building, will
// eventually be made private
func (parser *Parser) step2(is *InstructionSet) *CommandSequence {
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
		uuid, err := gouuid.NewV4()
		if err != nil {
			return nil
		}

		name := v.Registry + "/" + v.Project
		initialTag := name + ":" + uuid.String()

		// get docker registry credentials
		un := v.CfgUn
		pass := v.CfgPass
		email := v.CfgEmail

		buildOpts := docker.BuildImageOptions{
			Name:           initialTag,
			RmTmpContainer: true,
			ContextDir:     parser.contextDir,
			Auth: docker.AuthConfiguration{
				Username: un,
				Password: pass,
				Email:    email,
			},
			AuthConfigs: docker.AuthConfigurations{
				Configs: map[string]docker.AuthConfiguration{
					v.Registry: docker.AuthConfiguration{
						Password:      pass,
						Username:      un,
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
		for _, provided := range v.Tags {
			var tagValue = NewTag(provided).Evaluate(parser.contextDir)

			tagCmd := &TagCmd{
				Repo: name,
				Tag:  tagValue,
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
					Tag:       tagValue,
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
				UUID:       uuid.String(),
			},
			SubCommand: containerCommands,
		})
	}

	return ret
}
