package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modcloth/docker-builder/job"
)

const (
	TRAVIS_SUCCESS    = iota
	TRAVIS_NO_SUCCESS // pending or failure
)

const (
	TRAVIS_BUILD_TYPE_PUSH         = "push"
	TRAVIS_BUILD_TYPE_PULL_REQUEST = "pull_request"
)

type TravisRepository struct {
	Owner string `json:"owner_name"`
	Name  string `json:"name"`
}

type TravisPayload struct {
	Repository  TravisRepository `json:"repository"`
	CommitSHA   string           `json:"commit"`
	BuildStatus int              `json:"status"`
	BuildType   string           `json:"type"`
}

func Travis(req *http.Request) (spec *job.JobSpec, err error) {
	payloadBody := req.FormValue("payload")
	var payload = &TravisPayload{}
	if err = json.Unmarshal([]byte(payloadBody), payload); err != nil {
		return
	}

	if payload.BuildStatus != TRAVIS_SUCCESS {
		err = fmt.Errorf("build was not successful for %s/%s",
			payload.Repository.Owner, payload.Repository.Name)
		return
	}

	if payload.BuildType == TRAVIS_BUILD_TYPE_PULL_REQUEST {
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
