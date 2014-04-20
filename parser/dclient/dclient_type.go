package dclient

import (
	"github.com/rafecolton/bob/log"
)

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/wsxiaoys/terminal/color"
	"os"
)

/*
DockerClient is a wrapper for the go docker library.
*/
type DockerClient interface {
	Host() string
	LatestImageTaggedWith(uuid string) (string, error)
}

/*
NewDockerClient returns a new DockerClient (wrapper for a conneciton with a docker
daemon), properly initialized.  If you want a nullDockerClient for testing,
pass in nil as your logger and false for shouldBeReal.
*/
func NewDockerClient(logger log.Log, shouldBeReal bool) (DockerClient, error) {
	if logger == nil && !shouldBeReal {
		return &nullDockerClient{}, nil
	}

	var endpoint string

	defaultHost := os.Getenv("DOCKER_HOST")

	if defaultHost == "" {
		endpoint = "unix:///var/run/docker.sock"
	} else {
		// tcp endpoints cause a panic with this version of the go/docker library
		endpoint = defaultHost
	}

	dclient, err := docker.NewClient(endpoint)

	if err != nil {
		logger.Println(color.Sprintf("@{r!}Alas@{|}, docker host %s could not be reached\n----> %+v", endpoint, err))
		return nil, err
	}

	return &realDockerClient{
		client: dclient,
		host:   endpoint,
		Log:    logger,
	}, nil
}
