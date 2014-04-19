package dclient

import (
	"github.com/rafecolton/bob/log"
)

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"sort"
)

//for _, img := range imgs {
////runtime.Println("ID: ", img.ID)
////runtime.Println("RepoTags: ", img.RepoTags)
//runtime.Println("Created: ", img.Created)
////runtime.Println("Size: ", img.Size)
////runtime.Println("VirtualSize: ", img.VirtualSize)
////runtime.Println("ParentId: ", img.ParentId)
////runtime.Println("Repository: ", img.Repository)
//}

/*
 */
type Dclient struct {
	client *docker.Client
	host   string
	log.Log
}

func NewDclient(logger log.Log) (*Dclient, error) {
	var endpoint string

	defaultHost := os.Getenv("DOCKER_HOST")

	if defaultHost == "" {
		endpoint = "unix:///var/run/docker.sock"
	} else {
		endpoint = defaultHost
	}

	dclient, err := docker.NewClient(endpoint)

	if err != nil {
		logger.Println(color.Sprintf("@{r!}Alas@{|}, docker host %s could not be reached\n----> %+v", endpoint, err))
		return nil, err
	}

	return &Dclient{
		client: dclient,
		host:   endpoint,
		Log:    logger,
	}, nil
}

func (dclient *Dclient) LatestImage() (string, error) {
	var images APIImagesSlice
	images, err := dclient.client.ListImages(false)

	if err != nil {
		dclient.Println(color.Sprintf("@{r!}Alas@{|}, docker images could not be listed on  %s\n----> %+v", dclient.host, err))
		return "", err
	}

	// first is most recent
	sort.Sort(images)

	return images.FirstID(), nil
}
