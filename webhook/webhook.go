package webhook

import (
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"

	"github.com/modcloth/docker-builder/job"
)

var logger *logrus.Logger
var apiToken string

//Logger sets the (global) logger for the webhook package
func Logger(l *logrus.Logger) {
	logger = l
}

//APIToken sets the (global) apiToken for the webhook package
func APIToken(t string) {
	apiToken = t
}

func processJobHelper(spec *job.JobSpec, w http.ResponseWriter, req *http.Request) (int, string) {
	if err := spec.Validate(); err != nil {
		return 412, "412 precondition failed"
	}

	workdir, err := ioutil.TempDir("", "docker-build-worker")
	if err != nil {
		return 500, "500 internal server error"
	}

	gocleanup.Register(func() {
		fileutils.RmRF(workdir)
	})

	jobConfig := &job.JobConfig{
		Logger:         logger,
		Workdir:        workdir,
		GitHubAPIToken: apiToken,
	}

	job := job.NewJob(jobConfig, spec)

	// TODO: set this from somewhere
	var async bool

	// if async
	if async {
		if err = job.Process(); err != nil {
			return 417, "417 expectation failed"
		}

		return 201, "201 created"
	}

	// if not async
	go job.Process()
	return 202, "202 accepted"
}
