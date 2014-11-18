package buildrunner

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/sylphon/build-runner/builder"
	"github.com/sylphon/build-runner/builderfile"
	"github.com/sylphon/build-runner/conf"
	"github.com/sylphon/build-runner/parser"
)

// Stream corresponds to a file stream (stdout/stderr)
type Stream int

const (
	stdin Stream = iota

	// Stdout indicates the LogMsg's message should be printed, if at all, to stdout
	Stdout

	// Stderr indicates the LogMsg's message should be printed, if at all, to stderr
	Stderr
)

// LogMsg is the tentative interface for the structs returned on the log messages channel
type LogMsg interface {
	BuildID() string
	Level() int // type may change
	Msg() string
	Stream() Stream
}

// StatusMsg is the tentative interface for the structs returned on the status channel
type StatusMsg interface {
	BuildID() int
	Status() int // type may change
	Msg() string
	Error() error // should be checked for non-nil
}

// RunBuild runs a complete build for the provided Builderfile.  Currently, the
// channels argument is ignored but will be used in the future along with the
// LogMsg and StatusMsg interfaces
func RunBuild(unitConfig *builderfile.UnitConfig, contextDir string, channels ...chan interface{}) error {
	var err error
	var logger = logrus.New()
	var p *parser.Parser
	var bob *builder.Builder

	logger.Level = logrus.DebugLevel

	if err := envconfig.Process("build_runner", &conf.Config); err != nil {
		logger.WithField("err", err).Fatal("envconfig error")
	}

	if unitConfig == nil {
		return errors.New("unit config may not be nil")
	}
	if unitConfig.Version == 0 {
		unitConfig, err = builderfile.Convert0to1(unitConfig)
	}

	if err = unitConfig.HandleDeprecatedStanzas(); err != nil {
		return err
	}

	opts := parser.NewParserOptions{ContextDir: contextDir, Logger: logger}
	p = parser.NewParser(opts)

	instructionSet := p.InstructionSetFromBuilderfileStruct(unitConfig)
	commandSequence := p.CommandSequenceFromInstructionSet(instructionSet)

	bobOpts := builder.NewBuilderOptions{ContextDir: contextDir, Logger: logger}
	bob, err = builder.NewBuilder(bobOpts)
	if err != nil {
		return err
	}

	if buildErr := bob.BuildCommandSequence(commandSequence); buildErr != nil {
		return buildErr
	}

	return nil
}
