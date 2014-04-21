package dclient

import (
	"github.com/rafecolton/bob/log"
)

import (
	"github.com/fsouza/go-dockerclient"
	//"github.com/wsxiaoys/terminal/color"
	//"os"
	//"sort"
)

type realDockerClient struct {
	client *docker.Client
	host   string
	log.Log
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
func (rtoo *realDockerClient) LatestImageTaggedWithUUID(uuid string) (string, error) {
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
