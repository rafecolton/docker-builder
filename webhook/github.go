package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/modcloth/docker-builder/job"
)

var (
	githubSupportedEvents = map[string]bool{
		"push": true,
	}
)

type GithubOwner struct {
	Name string `json:"name"`
}

type GithubRepository struct {
	Name  string      `json:"name"`
	Owner GithubOwner `json:"owner"`
}

type GithubPushPayload struct {
	Repository GithubRepository `json:"repository"`
	CommitSHA  string           `json:"after"`
}

func Github(req *http.Request) (spec *job.JobSpec, err error) {
	event := req.Header.Get("X-Github-Event")
	if !githubSupportedEvents[event] {
		err = fmt.Errorf("Github event type %s is not supported.", event)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return
	}
	var payload = &GithubPushPayload{}
	if err = json.Unmarshal([]byte(body), payload); err != nil {
		return
	}

	spec = &job.JobSpec{
		RepoOwner: payload.Repository.Owner.Name,
		RepoName:  payload.Repository.Name,
		GitRef:    payload.CommitSHA,
	}

	return
}
