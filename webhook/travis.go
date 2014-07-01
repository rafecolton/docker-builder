package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modcloth/docker-builder/job"
)

const (
	travisSuccess = iota
)

const (
	travisBuildTypePullRequest = "pull_request"
)

type travisRepository struct {
	Owner string `json:"owner_name"`
	Name  string `json:"name"`
}

type travisPayload struct {
	Repository  travisRepository `json:"repository"`
	CommitSHA   string           `json:"commit"`
	BuildStatus int              `json:"status"`
	BuildType   string           `json:"type"`
}

/*
Travis parses a webhook HTTP request from Travis CI and returns a JobSpec.
*/
func Travis(req *http.Request) (spec *job.JobSpec, err error) {
	payloadBody := req.FormValue("payload")
	var payload = &travisPayload{}
	if err = json.Unmarshal([]byte(payloadBody), payload); err != nil {
		return
	}

	if payload.BuildStatus != travisSuccess {
		err = fmt.Errorf("build was not successful for %s/%s",
			payload.Repository.Owner, payload.Repository.Name)
		return
	}

	if payload.BuildType == travisBuildTypePullRequest {
		err = fmt.Errorf("won't build for pull request on %s/%s",
			payload.Repository.Owner, payload.Repository.Name)
		return
	}

	spec = &job.JobSpec{
		RepoOwner: payload.Repository.Owner,
		RepoName:  payload.Repository.Name,
		GitRef:    payload.CommitSHA,
	}

	return
}
