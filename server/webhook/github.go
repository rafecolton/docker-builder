package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rafecolton/docker-builder/job"
)

var (
	githubSupportedEvents = map[string]bool{
		"push": true,
	}
)

type githubOwner struct {
	Name string `json:"name"`
}

type githubRepository struct {
	Name  string      `json:"name"`
	Owner githubOwner `json:"owner"`
}

type githubPushPayload struct {
	Repository githubRepository `json:"repository"`
	CommitSHA  string           `json:"after"`
}

/*
Github parses a Github webhook HTTP request and returns a job.Spec.
*/
func Github(w http.ResponseWriter, req *http.Request) (int, string) {
	event := req.Header.Get("X-Github-Event")
	if !githubSupportedEvents[event] {
		logger.Errorf("Github event type %s is not supported.", event)
		return http.StatusBadRequest, fmt.Sprintf("%d bad request", http.StatusBadRequest)
	}
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		logger.Error(err)
		return http.StatusBadRequest, fmt.Sprintf("%d bad request", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	var payload = &githubPushPayload{}
	if err := decoder.Decode(payload); err != nil {
		logger.Error(err)
		return http.StatusBadRequest, fmt.Sprintf("%d bad request", http.StatusBadRequest)
	}

	spec := &job.Spec{
		RepoOwner: payload.Repository.Owner.Name,
		RepoName:  payload.Repository.Name,
		GitRef:    payload.CommitSHA,
	}

	return processJobHelper(spec, w, req)
}
