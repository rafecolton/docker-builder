package dockerclient

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"

	"github.com/fsouza/go-dockerclient"
	"github.com/rafecolton/go-dockerclient-sort"
)

var client *dockerClient

type dockerClient docker.Client

// DockerClient wraps docker.Client, adding a few handy functions
type DockerClient interface {
	// Client returns the underlying *docker.Client
	Client() *docker.Client

	// LatestImageByRegex returns the docker api image object if that object is
	// tagged with at least one tag that matches "regex"
	LatestImageByRegex(regex string) (*docker.APIImages, error)
}

// NewDockerClient returns the dockerclient used by the artifactory package
func NewDockerClient() (DockerClient, error) {
	if client != nil {
		return client, nil
	}

	endpoint, err := getEndpoint()
	if err != nil {
		return nil, err
	}
	certPath := os.Getenv("DOCKER_CERT_PATH")
	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY") != ""

	var dclient *docker.Client
	if endpoint.Scheme == "https" {
		cert := path.Join(certPath, "cert.pem")
		key := path.Join(certPath, "key.pem")
		ca := ""
		if tlsVerify {
			ca = path.Join(certPath, "ca.pem")
		}

		dclient, err = docker.NewTLSClient(endpoint.String(), cert, key, ca)
		if err != nil {
			return nil, err
		}
	} else {
		dclient, err = docker.NewClient(endpoint.String())
		if err != nil {
			return nil, err
		}
	}
	client = (*dockerClient)(dclient)
	return client, nil
}

func getEndpoint() (*url.URL, error) {
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		endpoint = "unix:///var/run/docker.sock"
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse endpoint %s as URL", endpoint)
	}
	if u.Scheme == "tcp" {
		_, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s for port", u.Host)
		}

		// Only reliable way to determine if we should be using HTTPS appears to be via port
		if os.Getenv("DOCKER_HOST_SCHEME") != "" {
			u.Scheme = os.Getenv("DOCKER_HOST_SCHEME")
		} else if port == "2376" {
			u.Scheme = "https"
		} else {
			u.Scheme = "http"
		}
	}
	return u, nil
}

func (client *dockerClient) LatestImageByRegex(regex string) (*docker.APIImages, error) {
	images, err := (*docker.Client)(client).ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}
	sort.Sort(dockersort.ByCreatedDescending(images))
	for _, image := range images {
		for _, tag := range image.RepoTags {
			matched, err := regexp.MatchString(regex, tag)
			if err != nil {
				return nil, err
			}
			if matched {
				return &image, nil
			}
		}
	}
	return nil, fmt.Errorf("unable to find image matching %q", regex)
}

// Client returns the underlying *docker.Client for calling all of its functions
func (client *dockerClient) Client() *docker.Client {
	return (*docker.Client)(client)
}
