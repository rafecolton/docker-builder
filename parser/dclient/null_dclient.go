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
LatestImageTaggedWithUUID is a mandatory method of the DockerClient interface.
*/
func (null *nullDockerClient) LatestImageTaggedWithUUID(uuid string) (string, error) {
	return "abcdef0123456789", nil
}
