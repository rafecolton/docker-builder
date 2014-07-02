package main

import (
	"net/http"

	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/webhook"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/modcloth/auth"
)

//ServerDescription is the help text for the `serer` command
const ServerDescription = `Start a small HTTP web server for receiving build requests.

Configure through the environment:

DOCKER_BUILDER_LOGLEVEL             =>     --log-level (global)
DOCKER_BUILDER_LOGFORMAT            =>     --log-format (global)
DOCKER_BUILDER_PORT                 =>     --port
DOCKER_BUILDER_APITOKEN             =>     --api-token
DOCKER_BUILDER_SKIPPUSH             =>     --skip-push
DOCKER_BUILDER_USERNAME             =>     --username
DOCKER_BUILDER_PASSWORD             =>     --password
DOCKER_BUILDER_TRAVISTOKEN          =>     --travis-token
DOCKER_BUILDER_NOTRAVIS             =>     --no-travis
DOCKER_BUILDER_GITHUBSECRET         =>     --github-secret
DOCKER_BUILDER_NOGITHUB             =>     --no-github

NOTE: If username and password are both empty (i.e. not provided), basic auth will not be used.
`

var apiToken, githubSecret, portString, pwd, travisToken, un string
var skipPush bool
var shouldTravis, shouldGitHub bool
var shouldBasicAuth, shouldTravisAuth, shouldGitHubAuth bool

var basicAuthFunc = func(http.ResponseWriter, *http.Request) {}
var travisAuthFunc = func(http.ResponseWriter, *http.Request) {}
var githubAuthFunc = func(http.ResponseWriter, *http.Request) {}

// define the server
func serve(c *cli.Context) {
	// set username and password (in helpers.go)
	setServerVars(c)

	// get "skip-push" option
	builder.SkipPush = skipPush
	// set up server
	server := martini.Classic()

	// check for basic auth
	if shouldBasicAuth {
		basicAuthFunc = auth.Basic(un, pwd)
	}

	// check for Travis auth
	if shouldTravisAuth {
		travisAuthFunc = auth.TravisCI(travisToken)
	}

	// check for GitHub auth
	if shouldGitHubAuth {
		githubAuthFunc = auth.GitHub(githubSecret)
	}

	// configure webhook globals
	webhook.Logger(Logger)
	webhook.APIToken(apiToken)

	// establish routes
	server.Get("/health", func() (int, string) { return 200, "200 OK" })
	server.Post("/docker-build", basicAuthFunc, webhook.DockerBuild)

	if shouldTravis {
		server.Post("/docker-build/travis", travisAuthFunc, webhook.Travis)
	}

	if shouldGitHub {
		server.Post("/docker-build/github", githubAuthFunc, webhook.Github)
	}

	// start server
	http.ListenAndServe(portString, server)
}
