package dclient

import (
	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

/*
DockerClient is a wrapper for the go docker library.
*/
type DockerClient interface {
	LatestImageTaggedWithUUID(uuid string) (string, error)
	RemoveImage(name string) error
	LatestRepoTaggedWithUUID(uuid string) (string, error)
	TagImage(name string, opts docker.TagImageOptions) error
}

type realDockerClient struct {
	client *docker.Client
	host   string
	*logrus.Logger
}

// returns fixed output, used for testing
type nullDockerClient struct {
	*logrus.Logger
}

func (null *nullDockerClient) TagImage(name string, opts docker.TagImageOptions) error {
	return nil
}

func (null *nullDockerClient) LatestImageTaggedWithUUID(uuid string) (string, error) {
	return "abcdef0123456789", nil
}

func (null *nullDockerClient) RemoveImage(name string) error {
	return nil
}

func (null *nullDockerClient) LatestRepoTaggedWithUUID(uuid string) (string, error) {
	return "", nil
}
