package dclient

// returns fixed output, used for testing
type nullDockerClient struct{}

/*
Host is a mandatory method of the DockerClient interface
*/
func (null *nullDockerClient) Host() string {
	return "null"
}

/*
LatestImageTaggedWith(uuid) is a mandatory method of the DockerClient interface.
*/
func (null *nullDockerClient) LatestImageTaggedWith(uuid string) (string, error) {
	return "abcdef0123456789", nil
}
