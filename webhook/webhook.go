package webhook

import (
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"

	"github.com/modcloth/docker-builder/job"
)

const (
	processJobSuccessCode    = 202
	processJobSuccessMessage = "202 accepted"
)

var logger *logrus.Logger
var apiToken string
var testMode bool

//Logger sets the (global) logger for the webhook package
func Logger(l *logrus.Logger) {
	logger = l
}

//APIToken sets the (global) apiToken for the webhook package
func APIToken(t string) {
	apiToken = t
}

//TestMode sets the (global) testMode variable for the webhook package
func TestMode(b bool) {
	testMode = b
}

func processJobHelper(spec *job.JobSpec, w http.ResponseWriter, req *http.Request) (int, string) {
	// If tests are running, don't actually attempt to build containers, just return success.
	// This is meant to allow testing ot the HTTP interactions for the webhooks
	if testMode {
		return processJobSuccessCode, processJobSuccessMessage
	}

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
	var sync bool

	// if sync
	if sync {
		if err = job.Process(); err != nil {
			return 417, "417 expectation failed"
		}

		return 201, "201 created"
	}

	// if async
	go job.Process()
	return processJobSuccessCode, processJobSuccessMessage
}
