package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/sylphon/build-runner/builder"
	"github.com/sylphon/build-runner/builderfile"
	"github.com/sylphon/build-runner/conf"
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

var example = &builderfile.Builderfile{
	Version: 1,
	ContainerArr: []*builderfile.ContainerSection{
		&builderfile.ContainerSection{
			Name:       "app",
			Dockerfile: "Dockerfile",
			Registry:   "quay.io/rafecolton",
			Project:    "build-runner-test",
			Tags:       []string{"latest", "git:sha", "git:tag", "git:branch"},
			SkipPush:   true,
		},
	},
}

func main() {
	if err := RunBuild(example, os.Getenv("GOPATH")+"/src/github.com/rafecolton/docker-builder"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RunBuild(fileZero *builderfile.Builderfile, contextDir string, channels ...chan interface{}) error {
	var fileOne *builderfile.Builderfile
	var err error
	var logger = logrus.New()
	logger.Level = logrus.DebugLevel

	if err := envconfig.Process("build_runner", &conf.Config); err != nil {
		logger.WithField("err", err).Fatal("envconfig error")
	}

	// set default config port
	if conf.Config.Port == 0 {
		conf.Config.Port = 5000
	}

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

	return nil
}
