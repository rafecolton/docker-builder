package builder

import (
	"errors"
	"io"
	"io/ioutil"
	"regexp"

	l "github.com/Sirupsen/logrus"
	"github.com/moby/moby/pkg/archive"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"
	"github.com/rafecolton/go-dockerclient-quick"

	"github.com/winchman/builder-core/communication"
	"github.com/winchman/builder-core/filecheck"
	"github.com/winchman/builder-core/parser"
)

var (
	// SkipPush will, when set to true, override any behavior set by a Bobfile and
	// will cause builders *NOT* to run `docker push` commands.  SkipPush is also set
	// by the `--skip-push` option when used on the command line.
	SkipPush bool

	imageWithTagRegex = regexp.MustCompile("^(.*):(.*)$")
)

/*
A Builder is the struct that actually does the work of moving files around and
executing the commands that do the docker build.
*/
type Builder struct {
	dockerClient    dockerclient.DockerClient
	workdir         string
	nextSubSequence *parser.SubSequence
	Stdout          io.Writer
	reporter        *comm.Reporter
	Builderfile     string
	contextDir      string

	// KeepTemporaryTag instructs the builder to keep the temporary tag, which
	// takes the form registry/project:<random=uuid>
	KeepTemporaryTag bool
}

/*
SetNextSubSequence sets the next subsequence within bob to be processed. This
function is exported because it is used explicitly in tests, but in Build(), it
is intended to be used as a helper function.
*/
func (bob *Builder) SetNextSubSequence(subSeq *parser.SubSequence) {
	bob.nextSubSequence = subSeq
}

// NewBuilderOptions encapsulates all of the options necessary for creating a
// new builder
type NewBuilderOptions struct {
	Log          comm.LogChan
	Event        comm.EventChan
	ContextDir   string
	dockerClient dockerclient.DockerClient // default to nil for regular docker client
}

/*
NewBuilder returns an instance of a Builder struct.  The function exists in
case we want to initialize our Builders with something.
*/
func NewBuilder(opts NewBuilderOptions) *Builder {
	var ret = &Builder{
		reporter:   comm.NewReporter(opts.Log, opts.Event),
		contextDir: opts.ContextDir,
	}

	ret.dockerClient = opts.dockerClient

	if opts.Log != nil {
		ret.Stdout = comm.NewLogEntryWriter(opts.Log)
	} else {
		ret.Stdout = ioutil.Discard /* /dev/null */
	}

	return ret
}

// BuildCommandSequence performs a build from a parser-generated CommandSequence struct
func (bob *Builder) BuildCommandSequence(commandSequence *parser.CommandSequence) error {
	bob.reporter.Event(comm.EventOptions{EventType: comm.RequestedEvent})

	if bob.dockerClient == nil {
		client, err := dockerclient.NewDockerClient()
		if err != nil {
			return err
		}
		bob.dockerClient = client
	}

	for _, seq := range commandSequence.Commands {
		var imageID string
		var err error

		if err := bob.cleanWorkdir(); err != nil {
			return err
		}
		bob.SetNextSubSequence(seq)
		if err := bob.setup(); err != nil {
			return err
		}

		bob.reporter.Log(
			l.WithField("container_section", seq.Metadata.Name),
			"running commands for container section",
		)

		for _, cmd := range seq.SubCommand {
			opts := &parser.DockerCmdOpts{
				DockerClient: bob.dockerClient,
				Image:        imageID,
				ImageUUID:    seq.Metadata.UUID,
				SkipPush:     SkipPush,
				Stdout:       bob.Stdout,
				Workdir:      bob.workdir,
				Reporter:     bob.reporter,
			}
			cmd = cmd.WithOpts(opts)

			bob.reporter.Log(l.WithField("command", cmd.Message()), "running docker command")

			if imageID, err = cmd.Run(); err != nil {
				switch err.(type) {
				case parser.NilClientError:
					continue
				default:
					return err
				}
			}

			bob.reporter.Log(
				l.WithFields(l.Fields{
					"command":  cmd.Message(),
					"image_id": imageID,
				}),
				"finished running docker command",
			)
		}

		if !bob.KeepTemporaryTag {
			bob.attemptToDeleteTemporaryUUIDTag(seq.Metadata.UUID)
		}
	}

	bob.reporter.Event(comm.EventOptions{EventType: comm.CompletedEvent})

	return nil
}

func (bob *Builder) attemptToDeleteTemporaryUUIDTag(uuid string) {
	if bob.dockerClient == nil {
		return
	}

	regex := ":" + uuid + "$"
	image, err := bob.dockerClient.LatestImageByRegex(regex)
	if err != nil {
		bob.reporter.LogLevel(
			l.WithField("err", err),
			"error getting repo taggged with temporary tag",
			l.WarnLevel,
		)
	}

	for _, tag := range image.RepoTags {
		matched, err := regexp.MatchString(regex, tag)
		if err != nil {
			return
		}
		if matched {
			bob.reporter.LogLevel(
				l.WithFields(l.Fields{
					"image_id": image.ID,
					"tag":      tag,
				}),
				"deleting temporary tag",
				l.DebugLevel,
			)

			if err = bob.dockerClient.Client().RemoveImage(tag); err != nil {
				bob.reporter.LogLevel(
					l.WithField("err", err),
					"error deleting temporary tag",
					l.WarnLevel,
				)
			}
			return
		}
	}
}

/*
Setup moves all of the correct files into place in the temporary directory in
order to perform the docker build.
*/
func (bob *Builder) setup() error {
	var workdir = bob.workdir
	var pathToDockerfile *filecheck.TrustedFilePath
	var err error

	if bob.nextSubSequence == nil {
		return errors.New("no command sub sequence set, cannot perform setup")
	}

	meta := bob.nextSubSequence.Metadata
	dockerfile := meta.Dockerfile
	opts := filecheck.NewTrustedFilePathOptions{File: dockerfile, Top: bob.contextDir}
	pathToDockerfile, err = filecheck.NewTrustedFilePath(opts)
	if err != nil {
		return err
	}

	if pathToDockerfile.Sanitize(); pathToDockerfile.State != filecheck.OK {
		return pathToDockerfile.Error
	}

	contextDir := pathToDockerfile.Top()
	tarStream, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{"Dockerfile"},
	})
	if err != nil {
		return err
	}

	defer tarStream.Close()
	if err := archive.Untar(tarStream, workdir, nil); err != nil {
		return err
	}

	if err := fileutils.CpWithArgs(
		contextDir+"/"+meta.Dockerfile,
		workdir+"/Dockerfile",
		fileutils.CpArgs{PreserveModTime: true},
	); err != nil {
		return err
	}

	return nil
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
cleanWorkdir effectively does a rm -rf and mkdir -p on bob's workdir.  Intended
to be used before using the workdir (i.e. before new command groups).
*/
func (bob *Builder) cleanWorkdir() error {
	workdir := bob.generateWorkDir()
	bob.workdir = workdir

	if err := fileutils.RmRF(workdir); err != nil {
		return err
	}

	return fileutils.MkdirP(workdir, 0755)
}
