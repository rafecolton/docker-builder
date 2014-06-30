package webhook

import (
	"net/http"

	"github.com/modcloth/docker-builder/job"
)

type Parser func(*http.Request) (*job.JobSpec, error)
