package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/modcloth/docker-builder/builderfile"

	"github.com/modcloth/go-dockerclient"
)

/*
An InstructionSet is an intermediate datatype - once a Builderfile is parsed
and the TOML is validated, the parser parses the data into an InstructionSet.
The primary purpose of this step is to merge any global container options into
the sections for the individual containers.
*/
type InstructionSet struct {
	DockerBuildOpts []string
	DockerTagOpts   []string
	Containers      []builderfile.ContainerSection
}

/*
A CommandSequence is an intermediate data type in the parsing process. Once a
Builderfile is parsed into an InstructionSet, it is further parsed into a
CommandSequence, which is essential an array of strings where each string is a
command to be run.
*/
type CommandSequence struct {
	Commands []*SubSequence
}

/*
SubSequenceMetadata contains any important metadata about the container build
such as the name of the Dockerfile and which files/dirs to exclude.
*/
type SubSequenceMetadata struct {
	Name       string
	Dockerfile string
	Included   []string
	Excluded   []string
	UUID       string
	SkipPush   bool
}

/*
A SubSequence is a logical grouping of commands such as a sequence of build,
tag, and push commands.  In addition, the subsequence metadata contains any
important metadata about the container build such as the name of the Dockerfile
and which files/dirs to exclude.
*/
type SubSequence struct {
	Metadata   *SubSequenceMetadata
	SubCommand []interface{}
}

//TagCmd is a wrapper for the docker TagImage functionality
type TagCmd struct {
	TagFunc func(name string, opts docker.TagImageOptions) error
	Image   string
	Force   bool
	Tag     string
	Repo    string
	msg     string
}

//Run is the command that actually calls TagImage to do the tagging
func (t *TagCmd) Run() error {
	var opts = &docker.TagImageOptions{
		Force: t.Force,
		Repo:  t.Repo,
		Tag:   t.Tag,
	}
	return t.TagFunc(t.Image, *opts)
}

//Message returns the shell command that would be equivalent to the TagImage command
func (t *TagCmd) Message() string {
	if t.msg == "" {
		msg := []string{"docker", "tag"}
		if t.Force {
			msg = append(msg, "--force")
		}
		msg = append(msg, t.Image)
		msg = append(msg, fmt.Sprintf("%s:%s", t.Repo, t.Tag))
		t.msg = strings.Join(msg, " ")
	}

	return t.msg
}

type PushCmd struct {
	PushFunc     func(opts docker.PushImageOptions, auth docker.AuthConfiguration) error
	Image        string
	Tag          string
	Registry     string
	AuthUn       string
	AuthPwd      string
	AuthEmail    string
	OutputStream io.Writer
}

func (p *PushCmd) Run() error {
	auth := &docker.AuthConfiguration{
		Username: p.AuthUn,
		Password: p.AuthPwd,
		Email:    p.AuthEmail,
	}
	opts := &docker.PushImageOptions{
		Name:         p.Image,
		Tag:          p.Tag,
		Registry:     p.Registry,
		OutputStream: p.OutputStream,
	}
	return p.PushFunc(*opts, *auth)
}

func (p *PushCmd) Message() string {
	return fmt.Sprintf("docker push %s:%s", p.Image, p.Tag)
}
