package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/job"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"
)

var apiToken string

func serve(c *cli.Context) {
	port := c.Int("port")
	portString := fmt.Sprintf(":%d", port)

	apiToken = c.String("api-token")

	if apiToken == "" {
		apiToken = config.APIToken
	}

	builder.SkipPush = c.Bool("skip-push") || config.SkipPush

	mux := http.NewServeMux()
	mux.HandleFunc("/docker-build", dockerBuild)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Use(loggingMiddleware())
	n.Run(portString)
}

func loggingMiddleware() *negronilogrus.Middleware {
	return negronilogrus.NewCustomMiddleware(Logger.Level, Logger.Formatter)
}

func dockerBuild(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	var spec = &job.JobSpec{}
	if err = json.Unmarshal([]byte(body), spec); err != nil {
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	if err = spec.Validate(); err != nil {
		http.Error(w, "412 precondition failed", http.StatusPreconditionFailed)
		return
	}

	workdir, err := ioutil.TempDir("", "docker-build-worker")
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	gocleanup.Register(func() {
		fileutils.RmRF(workdir)
	})

	jobConfig := &job.JobConfig{
		Logger:         Logger,
		Workdir:        workdir,
		GitHubAPIToken: apiToken,
	}

	job := job.NewJob(jobConfig, spec)

	if err = job.Process(); err != nil {
		http.Error(w, "417 expectation failed", http.StatusExpectationFailed)
		return
	}
}
