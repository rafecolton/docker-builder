package webhook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/modcloth/docker-builder/job"
)

/*
DockerBuild parses a simple JSON blob returns a JobSpec
*/
func DockerBuild(w http.ResponseWriter, req *http.Request) (int, string) {
	// TODO: check content type

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return 400, "400 bad request"
	}

	var spec = &job.JobSpec{}
	if err = json.Unmarshal([]byte(body), spec); err != nil {
		return 400, "400 bad request"
	}

	return processJobHelper(spec, w, req)
}
