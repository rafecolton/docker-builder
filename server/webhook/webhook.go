package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"

	"github.com/modcloth/docker-builder/job"
)

const (
	asyncSuccessCode    = 202
	asyncSuccessMessage = "202 accepted"
	syncSuccessCode     = 201
	syncSuccessMessage  = "201 created"
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
		if spec.Sync {
			return syncSuccessCode, syncSuccessMessage
		}
		return asyncSuccessCode, asyncSuccessMessage
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

	// if sync
	if spec.Sync {
		if err = job.Process(); err != nil {
			logger.WithField("error", err).Error("unable to process job synchronously")
			return 417, "417 expectation failed"
		}
		return syncSuccessCode, syncSuccessMessage
	}

	// if async
	go job.Process()

	retBytes, err := json.Marshal(job)
	if err != nil {
		return 409, fmt.Sprintf(`{"error": %q}`, err)
	}

	return asyncSuccessCode, string(retBytes)
}
