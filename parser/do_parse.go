package parser

import (
	"fmt"
	"github.com/rafecolton/bob/builderfile"
	"github.com/rafecolton/bob/parser/uuid"
	"io/ioutil"
)

import (
	"github.com/BurntSushi/toml"
)

// Step 1 of Parse
func (parser *Parser) getRaw() (string, error) {
	_ = uuid.NewUUIDGenerator(true)

	if !parser.IsOpenable() {
		return "", fmt.Errorf("%s is not openable", parser.filename)
	}

	bytes, err := ioutil.ReadFile(parser.filename)

	if err != nil {
		return "", err
	}

	raw := string(bytes)

	return raw, nil
}

// Step 2 of Parse
func (parser *Parser) rawToStruct() (*builderfile.Builderfile, error) {
	raw, err := parser.getRaw()
	if err != nil {
		return nil, err
	}

	file := &builderfile.Builderfile{}
	if _, err := toml.Decode(raw, &file); err != nil {
		return nil, err
	}

	if file.Docker.BuildOpts == nil {
		file.Docker.BuildOpts = []string{}
	}

	if file.Docker.TagOpts == nil {
		file.Docker.TagOpts = []string{}
	}

	return file, nil
}

// Step 3 of Parse
func (parser *Parser) structToInstructionSet() (*InstructionSet, error) {
	file, err := parser.rawToStruct()
	if err != nil {
		return nil, err
	}

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

	return ret, nil
}

// Step 4 of Parse()
func (parser *Parser) instructionSetToCommandSequence() (*CommandSequence, error) {
	//latestID, err := parser.dclient.LatestImage()
	_, err := parser.dclient.LatestImage()

	if err != nil {
		return nil, err
	}

	//TODO: fill this in

	return nil, nil
}

func (parser *Parser) finalStep() (interface{}, error) {
	return parser.instructionSetToCommandSequence()
}

/*
Parse further parses the Builderfile struct into an InstructionSet struct,
merging the global container options into the individual container sections.
*/
func (parser *Parser) Parse() (*CommandSequence, error) {
	ret, err := parser.finalStep()
	if err != nil {
		return nil, err
	}

	return ret.(*CommandSequence), nil
}
