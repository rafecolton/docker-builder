package dclient

import (
	"github.com/rafecolton/bob/log"
)

import (
	"github.com/fsouza/go-dockerclient"
)

/*
DockerClient is a wrapper for the go docker library.
*/
type DockerClient interface {
	LatestImageTaggedWithUUID(uuid string) (string, error)
}

type realDockerClient struct {
	client *docker.Client
	host   string
	log.Log
}

// returns fixed output, used for testing
type nullDockerClient struct {
	log.Log
}

/*
LatestImageTaggedWithUUID is a mandatory method of the DockerClient interface.
*/
func (null *nullDockerClient) LatestImageTaggedWithUUID(uuid string) (string, error) {
	return "abcdef0123456789", nil
}
