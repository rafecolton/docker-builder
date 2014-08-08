package builder

import (
	"github.com/modcloth/docker-builder/dclient"
	"github.com/modcloth/docker-builder/log"
	"github.com/modcloth/docker-builder/parser"
)

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/hishboy/gocommons/lang"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"
)

var imageWithTagRegex = regexp.MustCompile("^(.*):(.*)$")

/*
SkipPush, when set to true, will override any behavior set by a Bobfile and
will cause builders *NOT* to run `docker push` commands.  SkipPush is also set
by the `--skip-push` option when used on the command line.
*/
var SkipPush bool

/*
A Builder is the struct that actually does the work of moving files around and
executing the commands that do the docker build.
*/
type Builder struct {
	dockerClient dclient.DockerClient
	*logrus.Logger
	workdir         string
	isRegular       bool
	nextSubSequence *parser.SubSequence
	Stderr          io.Writer
	Stdout          io.Writer
	Builderfile     string
}

/*
SetNextSubSequence sets the next subsequence within bob to be processed. This
function is exported because it is used explicitly in tests, but in Build(), it
is intended to be used as a helper function.
*/
func (bob *Builder) SetNextSubSequence(subSeq *parser.SubSequence) {
	bob.nextSubSequence = subSeq
}

/*
NewBuilder returns an instance of a Builder struct.  The function exists in
case we want to initialize our Builders with something.
*/
func NewBuilder(logger *logrus.Logger, shouldBeRegular bool) (*Builder, error) {
	if logger == nil {
		logger = logrus.New()
		logger.Level = logrus.PanicLevel
	}

	client, err := dclient.NewDockerClient(logger, shouldBeRegular)

	if err != nil {
		return nil, err
	}

	if logrus.IsTerminal() {
		return &Builder{
			dockerClient: client,
			Logger:       logger,
			isRegular:    shouldBeRegular,
			Stdout:       log.NewOutWriter(logger, "         @{g}%s@{|}"),
			Stderr:       log.NewOutWriter(logger, "         @{r}%s@{|}"),
		}, nil
	}

	return &Builder{
		dockerClient: client,
		Logger:       logger,
		isRegular:    shouldBeRegular,
		Stdout:       log.NewOutWriter(logger, "         %s"),
		Stderr:       log.NewOutWriter(logger, "         %s"),
	}, nil
}

/*
BuildFromFile combines Build() with parser.Parse() to reduce the number of
steps needed to build with bob programatically.
*/
func (bob *Builder) BuildFromFile(file string) error {
	par, err := parser.NewParser(file, bob.Logger)
	if err != nil {
		return err
	}

	commandSequence, err := par.Parse()
	if err != nil {
		return err
	}

	bob.Builderfile = file

	return bob.Build(commandSequence)
}

/*
Build does the building!
*/
func (bob *Builder) Build(commandSequence *parser.CommandSequence) error {
	for _, seq := range commandSequence.Commands {
		if err := bob.CleanWorkdir(); err != nil {
			return err
		}
		bob.SetNextSubSequence(seq)
		if err := bob.Setup(); err != nil {
			return err
		}

		workdir := bob.Workdir()

		bob.WithFields(logrus.Fields{
			"container_section": seq.Metadata.Name,
		}).Info("running commands for container section")

		var imageID string
		var err error

		for _, cmd := range seq.SubCommand {
			opts := &parser.DockerCmdOpts{
				DockerClient: bob.dockerClient,
				Image:        imageID,
				ImageUUID:    seq.Metadata.UUID,
				SkipPush:     SkipPush,
				Stderr:       bob.Stderr,
				Stdout:       bob.Stdout,
				Workdir:      workdir,
			}

			cmd = cmd.WithOpts(opts)

			bob.WithField("command", cmd.Message()).Info("running docker command")

			if imageID, err = cmd.Run(); err != nil {
				return err
			}
		}

		repoWithTag, err := bob.dockerClient.LatestRepoTaggedWithUUID(seq.Metadata.UUID)
		if err != nil {
			bob.WithField("err", err).Warn("error getting repo taggged with temporary tag")
		}

		bob.WithField("tag", repoWithTag).Info("deleting temporary tag")

		if err = bob.dockerClient.RemoveImage(repoWithTag); err != nil {
			bob.WithField("err", err).Warn("error deleting temporary tag")

		}
	}

	return nil
}

/*
Setup moves all of the correct files into place in the temporary directory in
order to perform the docker build.
*/
func (bob *Builder) Setup() error {
	if bob.nextSubSequence == nil {
		return errors.New("no command sub sequence set, cannot perform setup")
	}

	meta := bob.nextSubSequence.Metadata
	fileSet := lang.NewHashSet()

	if len(meta.Included) == 0 {
		files, err := ioutil.ReadDir(bob.Repodir())
		if err != nil {
			return err
		}

		for _, v := range files {
			fileSet.Add(v.Name())
		}
	} else {
		for _, v := range meta.Included {
			fileSet.Add(v)
		}
	}

	// subtract any excludes from fileSet
	for _, exclude := range meta.Excluded {
		if fileSet.Contains(exclude) {
			fileSet.Remove(exclude)
		}
	}

	if fileSet.Contains("Dockerfile") {
		fileSet.Remove("Dockerfile")
	}

	// add the Dockerfile
	fileSet.Add(meta.Dockerfile)

	workdir := bob.Workdir()
	repodir := bob.Repodir()

	// copy the actual files over
	for _, file := range fileSet.ToSlice() {
		src := fmt.Sprintf("%s/%s", repodir, file)
		dest := fmt.Sprintf("%s/%s", workdir, file)

		if file == meta.Dockerfile {
			dest = fmt.Sprintf("%s/%s", workdir, "Dockerfile")
		}

		fileInfo, err := os.Stat(src)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			err = fileutils.CpR(src, dest)
		} else {
			err = fileutils.Cp(src, dest)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

/*
Repodir is the dir from which we are using files for our docker builds.
*/
func (bob *Builder) Repodir() string {
	if !bob.isRegular {
		repoDir := "spec/fixtures/repodir"
		return fmt.Sprintf("%s/%s", os.Getenv("PWD"), repoDir)
	}
	return filepath.Dir(bob.Builderfile)
}

/*
Workdir returns bob's working directory.
*/
func (bob *Builder) Workdir() string {
	return bob.workdir
}

func (bob *Builder) generateWorkDir() string {
	tmp, err := ioutil.TempDir("", "bob")
	if err != nil {
		return ""
	}

	gocleanup.Register(func() {
		fileutils.RmRF(tmp)
	})

	return tmp
}

/*
CleanWorkdir effectively does a rm -rf and mkdir -p on bob's workdir.  Intended
to be used before using the workdir (i.e. before new command groups).
*/
func (bob *Builder) CleanWorkdir() error {
	workdir := bob.generateWorkDir()
	bob.workdir = workdir

	if err := fileutils.RmRF(workdir); err != nil {
		return err
	}

	return fileutils.MkdirP(workdir, 0755)
}
