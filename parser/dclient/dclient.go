package dclient

import (
	"github.com/rafecolton/bob/log"
)

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/wsxiaoys/terminal/color"
	"os"
	//"sort"
)

type realDockerClient struct {
	client *docker.Client
	host   string
	log.Log
}

/*
DockerClient is a wrapper for the go docker library.
*/
type DockerClient interface {
	Host() string
	LatestImageTaggedWith(uuid string) (string, error)
}

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

/*
Host returns the name of the docker host we are trying to contact.
*/
func (rtoo *realDockerClient) Host() string {
	return rtoo.host
}

/*
LatestImageTaggedWith(uuid) returns the image id of the most recently created
docker image that has been tagged with the specified uuid.
*/
func (rtoo *realDockerClient) LatestImageTaggedWith(uuid string) (string, error) {
	return "abc", nil
	//[>
	//LatestImage - figure out what this does...
	//*/
	//func (dclient *realDclient) LatestImage() (string, error) {
	//var images APIImagesSlice
	//images, err := dclient.client.ListImages(false)

	//if err != nil {
	//dclient.Println(color.Sprintf("@{r!}Alas@{|}, docker images could not be listed on  %s\n----> %+v", dclient.host, err))
	//return "", err
	//}

	//// first is most recent
	//sort.Sort(images)

	//return images.FirstID(), nil
	//}
}
