package dclient

import "github.com/modcloth/go-dockerclient"

/*
APIImagesSlice is a type that wraps []docker.APIImages so that it can be sorted
using the "sort" package Interface.
*/
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

/*
FirstID - figure out what this does.
*/
func (slice APIImagesSlice) FirstID() string {
	return slice[0].ID
}
