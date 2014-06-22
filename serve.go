package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/job"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
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

	m := martini.Classic()

	m.Post("/docker-build", dockerBuild)

	http.ListenAndServe(portString, m) // instead of m.Run()
}

func dockerBuild(w http.ResponseWriter, req *http.Request) (int, string) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return 400, "400 bad request"
	}

	var spec = &job.JobSpec{}
	if err = json.Unmarshal([]byte(body), spec); err != nil {
		return 400, "400 bad request"
	}

	if err = spec.Validate(); err != nil {
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
		Logger:         Logger,
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
