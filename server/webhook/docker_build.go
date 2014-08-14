package webhook

import (
	"io/ioutil"
	"net/http"

	"github.com/modcloth/docker-builder/job"
)

/*
DockerBuild parses a simple JSON blob returns a job.Spec
*/
func DockerBuild(w http.ResponseWriter, req *http.Request) (int, string) {
	// TODO: check content type

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return 400, "400 bad request"
	}

	spec, err := job.NewSpec(body)
	if err != nil {
		return 400, "400 bad request"
	}

	return processJobHelper(spec, w, req)
}
