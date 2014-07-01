package webhook

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/martini-contrib/auth"

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
func Travis(w http.ResponseWriter, req *http.Request) (int, string) {
	payloadBody := req.FormValue("payload")
	var payload = &travisPayload{}

	if err := json.Unmarshal([]byte(payloadBody), payload); err != nil {
		logger.WithField("error", err).Error("error unmarshaling json")
		return 400, "400 bad request"
	}

	if payload.BuildStatus != travisSuccess {
		logger.WithFields(logrus.Fields{
			"owner": payload.Repository.Owner,
			"repo":  payload.Repository.Name,
		}).Error("travis build was not a success, not build")
		return 409, "409 conflict"
	}

	if payload.BuildType == travisBuildTypePullRequest {
		logger.WithFields(logrus.Fields{
			"owner": payload.Repository.Owner,
			"repo":  payload.Repository.Name,
		}).Error("build type is a pull request, not building")
		return 409, "409 conflict"
	}

	spec := &job.JobSpec{
		RepoOwner: payload.Repository.Owner,
		RepoName:  payload.Repository.Name,
		GitRef:    payload.CommitSHA,
	}

	return processJobHelper(spec, w, req)
}

/*
TravisAuth returns a Handler that authenticates via Travis's Authorization for
Webhooks scheme (http://docs.travis-ci.com/user/notifications/#Authorization-for-Webhooks)

Writes a http.StatusUnauthorized if authentication fails
*/
func TravisAuth(token string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		providedAuth := req.Header.Get("Authorization")

		travisRepoSlug := req.Header.Get("Travis-Repo-Slug")
		calculatedAuth := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s%s", travisRepoSlug, token))))

		if !auth.SecureCompare(providedAuth, calculatedAuth) {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}
	}

}
