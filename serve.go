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
	"github.com/martini-contrib/auth"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"
)

const SERVER_DESCRIPTION = `Start a small HTTP web server for receiving build requests.

Configure through the environment:

DOCKER_BUILDER_LOGLEVEL     =>     --log-level (global)
DOCKER_BUILDER_LOGFORMAT    =>     --log-format (global)
DOCKER_BUILDER_PORT         =>     --port
DOCKER_BUILDER_APITOKEN     =>     --api-token
DOCKER_BUILDER_SKIPPUSH     =>     --skip-push
DOCKER_BUILDER_USERNAME     =>     --username
DOCKER_BUILDER_PASSWORD     =>     --password

NOTE: If username and password are both empty (i.e. not provided), basic auth will not be used.
`

var apiToken string

// define the server
func serve(c *cli.Context) {
	// get vars from env and cli

	// set username and password (in helpers.go)
	setUnAndPwd(c)
	// get port
	port := c.Int("port")
	portString := fmt.Sprintf(":%d", port)
	// get api token
	apiToken = c.String("api-token")
	if apiToken == "" {
		apiToken = config.APIToken
	}
	// get "skip-push" option
	builder.SkipPush = c.Bool("skip-push") || config.SkipPush

	// set up server
	server := martini.Classic()

	// check for basic auth
	if un != "" && pwd != "" {
		server.Use(auth.Basic(un, pwd))
	}

	// establish routes
	server.Post("/docker-build", dockerBuild)

	// start server
	http.ListenAndServe(portString, server)
}

// handle a request
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
