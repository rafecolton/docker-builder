package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/sylphon/build-runner/builder"
	"github.com/sylphon/build-runner/builderfile"
	"github.com/sylphon/build-runner/parser"
)

type Stream int

const (
	stdin Stream = iota
	Stdout
	Stderr
)

type LogMsg interface {
	BuildID() string
	Level() int // type may change
	Msg() string
	Stream() Stream
}

type StatusMsg interface {
	BuildID() int
	Status() int // type may change
	Msg() string
	Error() error // should be checked for non-nil
}

var bf = &builderfile.Builderfile{
	Version: 1,
	ContainerArr: []*builderfile.ContainerSection{
		&builderfile.ContainerSection{
			Name:       "app",
			Dockerfile: "Dockerfile",
			Registry:   "quay.io/rafecolton",
			Project:    "docker-builder",
			Tags:       []string{"latest", "git:sha", "git:tag", "git:branch"},
			SkipPush:   true,
		},
	},
}

func main() {
	if err := RunBuild(bf, os.Getenv("GOPATH")+"/src/github.com/rafecolton/docker-builder"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
unit_config must include a unique ID that gets returned when reporting status
*/
func RunBuild(fileZero *builderfile.Builderfile, contextDir string, channels ...chan interface{}) error {
	var fileOne *builderfile.Builderfile
	var err error
	var logger = logrus.New()
	//logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.DebugLevel

	if fileZero == nil {
		return errors.New("builderfile may not be nil")
	}

	if fileZero.Version == 0 {
		fileOne, err = builderfile.Convert0to1(fileZero)
	} else {
		fileOne = fileZero
	}

	if err = fileOne.HandleDeprecatedStanzas(); err != nil {
		return err
	}

	p := parser.NewParser("", logger)

	instructionSet := p.InstructionSetFromBuilderfileStruct(fileOne)
	commandSequence := p.CommandSequenceFromInstructionSet(instructionSet)

	bob, err := builder.NewBuilder(logger, true)
	if err != nil {
		return err
	}

	if buildErr := bob.BuildCommandSequence(commandSequence); buildErr != nil {
		return buildErr
	}

	/* TODO:
	- struct => instruction set => command sequence
	*/

	/*
	  docker client
	  example config:
	  ---
	  docker:
	    build_opts:
	    - --force-rm
	    - --no-cache

	*/
	/*
	  TODO:
	  - parse unit config
	  - validate presence of contextDir
	  - do teh build
	  - report logs and status
	*/

	return nil
}

func log() {
}
