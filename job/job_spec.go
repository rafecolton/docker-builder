package job

import (
	"encoding/json"
	"errors"
)

/*
JobSpec contains the specs for a job, retrieved from parsed JSON
*/
type JobSpec struct {
	RepoOwner      string `json:"account"`
	RepoName       string `json:"repo"`
	GitRef         string `json:"ref"`
	GitHubAPIToken string `json:"api_token"`
	Depth          string `json:"depth"`
}

/*
NewJobSpec creates a new job spec based on the arguments that would be passed
along from the job goworker picks up from Redis.
*/
func NewJobSpec(args ...interface{}) (*JobSpec, error) {
	return extractJobSpecFromRawArgs(args...)
}

/*
Validate checks that required fields are present in the spec.
*/
func (spec *JobSpec) Validate() error {

	if spec.RepoOwner == "" {
		return errors.New("account must be provided for job spec")
	}

	if spec.RepoName == "" {
		return errors.New("repo must be provided for job spec")
	}

	if spec.GitRef == "" {
		return errors.New("ref must be provided for job spec")
	}

	return nil
}

func extractJobSpecFromRawArgs(args ...interface{}) (*JobSpec, error) {
	var ret = &JobSpec{}

	if len(args) < 1 {
		return nil, errors.New("a single build spec object argument is required")
	}

	rawBuildJobSpec, ok := args[0].(interface{})
	if !ok {
		return nil, errors.New("build spec args must be an object")
	}

	argJSONBytes, err := json.Marshal(rawBuildJobSpec)
	if err != nil {
		return nil, errors.New("failed to re-serialize build job spec object")
	}

	err = json.Unmarshal(argJSONBytes, ret)
	if err != nil {
		return nil, errors.New("failed to deserialize build job spec as JSON")
	}

	if err = ret.Validate(); err != nil {
		return nil, err
	}

	return ret, nil
}
