package runner

import (
	"errors"

	"github.com/Sirupsen/logrus"
	b "github.com/winchman/builder-core/builder"
	"github.com/winchman/builder-core/communication"
	p "github.com/winchman/builder-core/parser"
	"github.com/winchman/builder-core/unit-config"
)

// Flag is an option type for the RunBuild command
type Flag uint8

const (
	// KeepTemporaryTag instructs the builder not to delete the random uuid tag
	KeepTemporaryTag Flag = 1 << iota // KeepTemporaryTag == 1 (iota has been reset)

	noop1 Flag = 1 << iota // noop1 == 2 // here for example/testing purposes
	noop2 Flag = 1 << iota // noop2 == 4 // here for example/testing purposes
)

// Options encapsulates the options for RunBuild/RunBuildSynchronously
type Options struct {
	UnitConfig *unitconfig.UnitConfig
	ContextDir string

	// LogLevel is only used for RunBuildSynchronously, ignored for RunBuild
	// LogLevel defaults to PanicLevel if not set
	LogLevel logrus.Level
}

func shouldKeepTemporaryTag(flags []Flag) bool {
	var total Flag
	for _, flag := range flags {
		total |= flag
	}
	return total&KeepTemporaryTag == KeepTemporaryTag
}

// RunBuild runs a complete build for the provided unit config.  Currently, the
// channels argument is ignored but will be used in the future along with the
// LogMsg and StatusMsg interfaces
func RunBuild(opts Options, flags ...Flag) (comm.LogChan, comm.EventChan, comm.ExitChan) {

	var unitConfig = opts.UnitConfig
	var contextDir = opts.ContextDir

	var log = make(chan comm.LogEntry, 1)
	var event = make(chan comm.Event, 1)
	var exit = make(chan error)

	go func() {
		var err error

		if unitConfig == nil {
			exit <- errors.New("unit config may not be nil")
			return
		}

		parser := p.NewParser(p.NewParserOptions{
			ContextDir: contextDir,
			Log:        log,
			Event:      event,
		})
		commandSequence := parser.Parse(unitConfig)

		builder := b.NewBuilder(b.NewBuilderOptions{
			ContextDir: contextDir,
			Log:        log,
			Event:      event,
		})
		builder.KeepTemporaryTag = shouldKeepTemporaryTag(flags)
		if err = builder.BuildCommandSequence(commandSequence); err != nil {
			exit <- err
			return
		}

		exit <- nil
	}()

	return log, event, exit
}

// RunBuildSynchronously - run a build, wait for it to finish, log to stdout
func RunBuildSynchronously(opts Options, flags ...Flag) error {
	var logger = logrus.New()
	logger.Level = opts.LogLevel
	log, status, done := RunBuild(opts, flags...) // make sure to update this as needed
	for {
		select {
		case e, ok := <-log:
			if !ok {
				return errors.New("log channel closed prematurely")
			}
			e.LogWithLogger(logger)
		case event, ok := <-status:
			if !ok {
				return errors.New("status channel closed prematurely")
			}
			logger.WithFields(event.Data()).Debugf("status event (type %s)", event.EventType())
		case err, ok := <-done:
			if !ok {
				return errors.New("exit channel closed prematurely")
			}
			return err
		}
	}
}
