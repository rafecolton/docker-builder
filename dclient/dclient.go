package dclient

import (
	"github.com/rafecolton/bob/log"
)

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/wsxiaoys/terminal/color"
)

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
)

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
LatestImageTaggedWith(uuid) returns the image id of the most recently created
docker image that has been tagged with the specified uuid.
*/
func (rtoo *realDockerClient) LatestImageTaggedWithUUID(uuid string) (string, error) {
	/*
		LatestImage - figure out what this does...
	*/
	var images APIImagesSlice
	images, err := rtoo.client.ListImages(false)

	if err != nil {
		rtoo.Println(color.Sprintf("@{r!}Alas@{|}, docker images could not be listed on  %s\n----> %+v", rtoo.host, err))
		return "", err
	}

	// first is most recent
	sort.Sort(images)

	for _, image := range images {
		for _, tag := range image.RepoTags {
			matched, err := regexp.MatchString(fmt.Sprintf(":%s$", uuid), tag)
			if err != nil {
				return "", err
			}

			if matched {
				return image.ID, nil
			}
		}
	}

	return "", errors.New(color.Sprintf("@{r!}Alas@{|}, I am unable to find image tagged with uuid \"%s\"", uuid))
}
