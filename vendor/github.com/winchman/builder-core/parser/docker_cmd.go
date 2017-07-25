package parser

import (
	"io"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/rafecolton/go-dockerclient-quick"
	"github.com/winchman/builder-core/communication"
)

/*
DockerCmdOpts is an options struct for the options required by the various
structs that implement the DockerCmd interface
*/
type DockerCmdOpts struct {
	DockerClient dockerclient.DockerClient
	Image        string
	Workdir      string
	Stdout       io.Writer
	Stderr       io.Writer
	SkipPush     bool
	ImageUUID    string
	Reporter     *comm.Reporter
}

/*
DockerCmd is an interface that wraps the various docker command types.
*/
type DockerCmd interface {
	// Run() runs the underlying command. The string return value is expected
	// to be the ID of the image being operated on
	Run() (string, error)

	// Message() returns a string representation of the command if it were to
	// be run on the command line
	Message() string

	// WithOpts sets the options for the command. It is expected to return the
	// same DockerCmd in a state in which the Run() function can be called
	// immediately after without error (i.e.`dockerCmdInstance.WithOpts(opts).Run()`)
	WithOpts(opts *DockerCmdOpts) DockerCmd
}

//BuildCmd is a wrapper for the os/exec call for `docker build`
type BuildCmd struct {
	opts          *DockerCmdOpts
	buildOpts     docker.BuildImageOptions
	origBuildOpts []string
	reporter      *comm.Reporter
	test          bool
}

//WithOpts sets options required for the BuildCmd
func (b *BuildCmd) WithOpts(opts *DockerCmdOpts) DockerCmd {
	b.opts = opts
	b.reporter = opts.Reporter
	if opts.DockerClient.Client().HTTPClient == nil {
		b.test = true
	}
	return b
}

// NilClientError is the error returned by any Run() command if the underlying
// docker client is nil
type NilClientError struct{}

func (err NilClientError) Error() string {
	return "docker client is nil"
}

//Run is the command that actually calls docker build shell command.  Determine
//the image ID for the resulting image and return that as well.
func (b *BuildCmd) Run() (string, error) {
	var opts = b.opts
	if b.test {
		return opts.ImageUUID, NilClientError{}
	}

	b.reporter.Event(comm.EventOptions{EventType: comm.BuildEvent})

	buildOpts := b.buildOpts
	buildOpts.OutputStream = opts.Stdout
	buildOpts.ContextDir = opts.Workdir

	if err := opts.DockerClient.Client().BuildImage(buildOpts); err != nil {
		b.reporter.Event(comm.EventOptions{
			EventType: comm.BuildCompletedEvent,
			Data: map[string]interface{}{
				"uuid_tag": opts.ImageUUID,
				"error":    err,
			},
		})
		return "", err
	}

	image, err := opts.DockerClient.LatestImageByRegex(":" + opts.ImageUUID + "$")
	if err != nil {
		b.reporter.Event(comm.EventOptions{
			EventType: comm.BuildCompletedEvent,
			Data: map[string]interface{}{
				"uuid_tag": opts.ImageUUID,
				"error":    err,
			},
		})
		return "", err
	}

	b.reporter.Event(comm.EventOptions{
		EventType: comm.BuildCompletedEvent,
		Data: map[string]interface{}{
			"image_id": image.ID,
			"uuid_tag": opts.ImageUUID,
			"error":    nil,
		},
	})
	return image.ID, nil
}

//Message returns the shell command that gets run for docker build commands
func (b *BuildCmd) Message() string {
	ret := []string{"docker", "build", "-t", b.buildOpts.Name}
	ret = append(ret, b.origBuildOpts...)
	ret = append(ret, ".")
	return strings.Join(ret, " ")
}

//TagCmd is a wrapper for the docker TagImage functionality
type TagCmd struct {
	TagFunc  func(name string, opts docker.TagImageOptions) error
	Image    string
	Force    bool
	Tag      string
	Repo     string
	msg      string
	test     bool
	reporter *comm.Reporter
}

//WithOpts sets options required for the TagCmd
func (t *TagCmd) WithOpts(opts *DockerCmdOpts) DockerCmd {
	t.Image = opts.Image
	t.reporter = opts.Reporter
	if opts.DockerClient.Client().HTTPClient == nil {
		t.test = true
		return t
	}
	t.TagFunc = opts.DockerClient.Client().TagImage
	return t
}

//Run is the command that actually calls TagImage to do the tagging
func (t *TagCmd) Run() (string, error) {
	t.reporter.Event(comm.EventOptions{
		EventType: comm.TagEvent,
		Data: map[string]interface{}{
			"image_id": t.Image,
			"repo":     t.Repo,
			"tag":      t.Tag,
		},
	})
	var opts = &docker.TagImageOptions{
		Force: t.Force,
		Repo:  t.Repo,
		Tag:   t.Tag,
	}
	if t.test {
		return t.Image, NilClientError{}
	}
	err := t.TagFunc(t.Image, *opts)
	t.reporter.Event(comm.EventOptions{
		EventType: comm.TagCompletedEvent,
		Data: map[string]interface{}{
			"image_id": t.Image,
			"repo":     t.Repo,
			"tag":      t.Tag,
			"error":    err,
		},
	})

	return t.Image, err
}

//Message returns the shell command that would be equivalent to the TagImage command
func (t *TagCmd) Message() string {
	if t.msg == "" {
		msg := []string{"docker", "tag"}
		if t.Force {
			msg = append(msg, "--force")
		}
		msg = append(msg, t.Image)
		msg = append(msg, t.Repo+":"+t.Tag)
		t.msg = strings.Join(msg, " ")
	}

	return t.msg
}

//PushCmd is a wrapper for the docker PushImage functionality
type PushCmd struct {
	PushFunc     func(opts docker.PushImageOptions, auth docker.AuthConfiguration) error
	Image        string
	Tag          string
	Registry     string
	AuthUn       string
	AuthPwd      string
	AuthEmail    string
	OutputStream io.Writer

	skip     bool
	imageID  string
	test     bool
	reporter *comm.Reporter
}

//WithOpts sets options required for the PushCmd
func (p *PushCmd) WithOpts(opts *DockerCmdOpts) DockerCmd {
	if opts.DockerClient.Client().HTTPClient == nil {
		p.test = true
		return p
	}
	p.OutputStream = opts.Stdout
	p.PushFunc = opts.DockerClient.Client().PushImage
	p.skip = opts.SkipPush
	p.imageID = opts.Image
	p.reporter = opts.Reporter
	return p
}

//Run is the command that actually calls PushImage to do the pushing
func (p *PushCmd) Run() (string, error) {
	if p.skip || p.test {
		return p.imageID, nil
	}

	p.reporter.Event(comm.EventOptions{
		EventType: comm.PushEvent,
		Data: map[string]interface{}{
			"repo":     p.Image,
			"tag":      p.Tag,
			"registry": p.Registry,
		},
	})

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

	err := p.PushFunc(*opts, *auth)
	p.reporter.Event(comm.EventOptions{
		EventType: comm.PushEvent,
		Data: map[string]interface{}{
			"repo":     p.Image,
			"tag":      p.Tag,
			"registry": p.Registry,
			"error":    err,
		},
	})

	return p.imageID, err
}

//Message returns the shell command that would be equivalent to the PushImage command
func (p *PushCmd) Message() string {
	return "docker push " + p.Image + ":" + p.Tag
}
