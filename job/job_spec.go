package job

import (
	"encoding/json"
	"errors"
)

/*
Spec contains the specs for a job, retrieved from parsed JSON
*/
type Spec struct {
	Bobfile        string `json:"bobfile"`
	RepoOwner      string `json:"account"`
	RepoName       string `json:"repo"`
	GitRef         string `json:"ref"`
	GitHubAPIToken string `json:"api_token"`
	Depth          string `json:"depth"`
	Sync           bool   `json:"sync"`
}

/*
NewSpec creates a new job spec from raw json data
*/
func NewSpec(raw []byte) (*Spec, error) {
	var ret = &Spec{}
	if err := json.Unmarshal(raw, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

/*
Validate checks that required fields are present in the spec.
*/
func (spec *Spec) Validate() error {

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
