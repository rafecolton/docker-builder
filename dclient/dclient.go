package dclient

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
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
		quietLogger.Level = logrus.PanicLevel

		return &nullDockerClient{
			Logger: quietLogger,
		}, nil
	}

	endpoint, err := getEndpoint()
	if err != nil {
		return nil, err
	}
	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY") != ""
	certPath := os.Getenv("DOCKER_CERT_PATH")

	var dclient *docker.Client
	if endpoint.Scheme == "https" {
		if certPath == "" {
			return nil, fmt.Errorf("Using TLS, but DOCKER_CERT_PATH is empty")
		}

		cert := path.Join(certPath, "cert.pem")
		key := path.Join(certPath, "key.pem")
		ca := ""
		if tlsVerify {
			ca = path.Join(certPath, "ca.pem")
		}

		dclient, err = docker.NewTLSClient(endpoint.String(), cert, key, ca)
	} else {
		dclient, err = docker.NewClient(endpoint.String())
	}

	if err != nil {
		logger.WithFields(logrus.Fields{
			"docker_host": endpoint,
			"error":       err,
		}).Error("docker host %q could not be reached")

		return nil, err
	}

	return &realDockerClient{
		client: dclient,
		host:   endpoint.String(),
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
			matched, err := regexp.MatchString(":"+uuid+"$", tag)
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
			matched, err := regexp.MatchString(":"+uuid+"$", tag)
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

func (rtoo *realDockerClient) PushImage(opts docker.PushImageOptions, auth docker.AuthConfiguration) error {
	return rtoo.client.PushImage(opts, auth)
}

func (rtoo *realDockerClient) BuildImage(opts docker.BuildImageOptions) error {
	return rtoo.client.BuildImage(opts)
}

// Workaround since DOCKER_HOST typically is tcp:// but we need to vary whether
// we use HTTP/HTTPS when interacting with the API
// Can be removed if https://github.com/fsouza/go-dockerclient/issues/173 is
// resolved
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
