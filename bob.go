package bob

import (
	"github.com/modcloth/bob/dclient"
	"github.com/modcloth/bob/log"
	"github.com/modcloth/bob/parser"
)

import (
	"github.com/hishboy/gocommons/lang"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"
	"github.com/wsxiaoys/terminal/color"
)

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

/*
A Builder is the struct that actually does the work of moving files around and
executing the commands that do the docker build.
*/
type Builder struct {
	dockerClient dclient.DockerClient
	log.Logger
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
func NewBuilder(logger log.Logger, shouldBeRegular bool) (*Builder, error) {
	if logger == nil {
		logger = &log.NullLogger{}
	}

	client, err := dclient.NewDockerClient(logger, shouldBeRegular)

	if err != nil {
		return nil, err
	}

	return &Builder{
		dockerClient: client,
		Logger:       logger,
		isRegular:    shouldBeRegular,
		Stdout:       log.NewOutWriter(logger, "         @{g}%s@{|}"),
		Stderr:       log.NewOutWriter(logger, "         @{r}%s@{|}"),
	}, nil
}

/*
Build is currently a placeholder function but will eventually be used to do the
actual work of building.
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

		bob.Println(color.Sprintf("@{w!}  ----->  Running commands for %q @{|}", seq.Metadata.Name))

		var imageID string
		var err error

		for _, cmd := range seq.SubCommand {
			cmd.Stdout = bob.Stdout
			cmd.Stderr = bob.Stderr
			cmd.Dir = workdir

			if cmd.Path == "docker" {
				path, err := exec.LookPath("docker")
				if err != nil {
					return err
				}

				cmd.Path = path
			}

			switch cmd.Args[1] {
			case "build":
				bob.Println(color.Sprintf("@{w!}  ----->  Running command %s @{|}", cmd.Args))
				if err := cmd.Run(); err != nil {
					return err
				}

				imageID, err = bob.LatestImageTaggedWithUUID(seq.Metadata.UUID)
				if err != nil {
					return err
				}
			case "tag":
				cmd.Args[2] = imageID
				bob.Println(color.Sprintf("@{w!}  ----->  Running command %s @{|}", cmd.Args))

				if err := cmd.Run(); err != nil {
					return err
				}
			case "push":
				bob.Println(color.Sprintf("@{w!}  ----->  Running command %s @{|}", cmd.Args))

				if err := cmd.Run(); err != nil {
					return err
				}
			default:
				return errors.New(
					color.Sprintf(
						"@{r!}oops, looks like the command you're asking me to run is improperly formed:@{|} %s\n",
						cmd.Args,
					),
				)
			}
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
		os.RemoveAll(tmp)
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

	if err := os.RemoveAll(workdir); err != nil {
		return err
	}

	if err := os.MkdirAll(workdir, 0755); err != nil {
		return err
	}

	return nil
}

/*
LatestImageTaggedWithUUID accepts a uuid and invokes the underlying utility
DockerClient to determine the id of the most recently created image tagged with
the provided uuid.
*/
func (bob *Builder) LatestImageTaggedWithUUID(uuid string) (string, error) {
	id, err := bob.dockerClient.LatestImageTaggedWithUUID(uuid)
	if err != nil {
		bob.Println(err)
		return "", err
	}

	return id, nil
}
