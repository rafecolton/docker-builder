package dclient

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

/*
NewDockerClient returns a new DockerClient (wrapper for a conneciton with a docker
daemon), properly initialized.  If you want a nullDockerClient for testing,
pass in nil as your logger and false for shouldBeReal.
*/
func NewDockerClient(logger *logrus.Logger, shouldBeReal bool) (DockerClient, error) {
	if logger == nil && !shouldBeReal {
		quietLogger := logrus.New()
		quietLogger.Level = logrus.Panic

		return &nullDockerClient{
			Logger: quietLogger,
		}, nil
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
		logger.WithFields(logrus.Fields{
			"docker_host": endpoint,
			"error":       err,
		}).Error("docker host %q could not be reached")

		return nil, err
	}

	return &realDockerClient{
		client: dclient,
		host:   endpoint,
		Logger: logger,
	}, nil
}

func (rtoo *realDockerClient) RemoveImage(name string) error {
	return rtoo.client.RemoveImage(name)
}

/*
LatestImageTaggedWith(uuid) returns the image id of the most recently created
docker image that has been tagged with the specified uuid.
*/
func (rtoo *realDockerClient) LatestImageTaggedWithUUID(uuid string) (string, error) {
	images, err := rtoo.sortedImages()
	if err != nil {
		return "", err
	}

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

	return "", fmt.Errorf("unable to find image tagged with uuid %q", uuid)
}

func (rtoo *realDockerClient) LatestRepoTaggedWithUUID(uuid string) (string, error) {
	images, err := rtoo.sortedImages()
	if err != nil {
		return "", err
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			matched, err := regexp.MatchString(fmt.Sprintf(":%s$", uuid), tag)
			if err != nil {
				return "", err
			}

			if matched {
				return tag, nil
			}
		}
	}

	return "", fmt.Errorf("unable to find image tagged with uuid %q", uuid)
}

func (rtoo *realDockerClient) sortedImages() (APIImagesSlice, error) {
	/*
		LatestImage - figure out what this does...
	*/
	var images APIImagesSlice
	images, err := rtoo.client.ListImages(false)

	if err != nil {
		rtoo.WithFields(logrus.Fields{
			"docker_host": rtoo.host,
			"error":       err,
		}).Error("docker images could not be listed")

		return nil, err
	}

	// first is most recent
	sort.Sort(images)

	return images, nil
}

func (rtoo *realDockerClient) TagImage(name string, opts docker.TagImageOptions) error {
	return rtoo.client.TagImage(name, opts)
}
