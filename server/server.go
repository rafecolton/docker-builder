package server

import (
	"net/http"

	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/conf"
	"github.com/modcloth/docker-builder/webhook"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/modcloth/auth"
)

var logger *logrus.Logger
var config = conf.Config
var server = martini.Classic()

func Logger(l *logrus.Logger) {
	logger = l
}

func Serve(context *cli.Context) {
	// set vars
	setVarsFromContext(context)

	// get "skip-push" option
	builder.SkipPush = skipPush

	// set up auth functions
	if shouldBasicAuth {
		basicAuthFunc = auth.Basic(un, pwd)
	}
	if shouldTravisAuth {
		travisAuthFunc = auth.TravisCI(travisToken)
	}
	if shouldGitHubAuth {
		githubAuthFunc = auth.GitHub(githubSecret)
	}

	// configure webhooks
	webhook.Logger(logger)
	webhook.APIToken(apiToken)
	if shouldTravis {
		server.Post("/docker-build/travis", travisAuthFunc, webhook.Travis)
	}
	if shouldGitHub {
		server.Post("/docker-build/github", githubAuthFunc, webhook.Github)
	}

	// establish routes
	server.Get("/health", func() (int, string) { return 200, "200 OK" })
	server.Post("/docker-build", basicAuthFunc, webhook.DockerBuild)

	// start server
	http.ListenAndServe(portString, server)
}
