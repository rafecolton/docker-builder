package webhook

import (
	"encoding/json"
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

func processJobHelper(spec *job.Spec, w http.ResponseWriter, req *http.Request) (int, string) {
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

	jobConfig := &job.Config{
		Logger:         logger,
		Workdir:        workdir,
		GitHubAPIToken: apiToken,
	}

	j := job.NewJob(jobConfig, spec)

	// if sync
	if spec.Sync || job.TestMode {
		if err = j.Process(); err != nil {
			logger.WithField("error", err).Error("unable to process job synchronously")
			return 417, `{"error": "` + err.Error() + `"}`
		}
		retBytes, err := json.Marshal(j)
		if err != nil {
			return 417, `{"error": "` + err.Error() + `"}`
		}

		return syncSuccessCode, string(retBytes)
	}

	// if async
	go j.Process()

	retBytes, err := json.Marshal(j)
	if err != nil {
		return 409, `{"error": "` + err.Error() + `"}`
	}

	return asyncSuccessCode, string(retBytes)
}
