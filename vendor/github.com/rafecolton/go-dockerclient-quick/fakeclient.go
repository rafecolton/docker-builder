package dockerclient

import (
	"github.com/fsouza/go-dockerclient"
)

type fakeClient struct{}

func (client *fakeClient) Client() *docker.Client {
	return &docker.Client{}
}
func (client *fakeClient) LatestImageByRegex(regex string) (*docker.APIImages, error) {
	return &docker.APIImages{}, nil
}

// FakeClient returns a DockerClient implementation that is suitable for testing
func FakeClient() DockerClient {
	return &fakeClient{}
}
