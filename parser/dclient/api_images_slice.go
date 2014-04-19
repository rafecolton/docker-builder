package dclient

import "github.com/fsouza/go-dockerclient"

type APIImagesSlice []docker.APIImages

func (slice APIImagesSlice) Len() int {
	return len(slice)
}

func (slice APIImagesSlice) Less(i, j int) bool {
	return slice[i].Created > slice[j].Created
}

func (slice APIImagesSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice APIImagesSlice) FirstID() string {
	return slice[0].ID
}
