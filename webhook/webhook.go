package webhook

import (
	"net/http"

	"github.com/modcloth/docker-builder/job"
)

/*
Parser is any function that parses a webhook request and returns a job spec.
*/
type Parser func(*http.Request) (*job.JobSpec, error)
