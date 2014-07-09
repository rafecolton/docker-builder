package server

import (
	"net/http"
	"time"

	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/job"
	"github.com/modcloth/docker-builder/server/webhook"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/modcloth/auth"
)

var logger *logrus.Logger
var server *martini.ClassicMartini
var skipLogging = map[string]bool{
	"/health": true,
}

//Logger sets the (global) logger for the server package
func Logger(l *logrus.Logger) {
	logger = l
}

// Serve sets everything up and runs the docker-builder server
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

	server = setupServer()

	if shouldTravis {
		server.Post("/docker-build/travis", travisAuthFunc, webhook.Travis)
	}
	if shouldGitHub {
		server.Post("/docker-build/github", githubAuthFunc, webhook.Github)
	}

	// base routes
	server.Get("/health", func() (int, string) { return 200, "200 OK" })
	server.Post("/docker-build", basicAuthFunc, webhook.DockerBuild)

	// job control routes
	server.Group("/jobs", func(r martini.Router) {
		r.Get("/:id", job.Get)
		r.Get("/:id/tail", job.TailN)
		r.Post("", webhook.DockerBuild)
		r.Get("", job.GetAll)
	}, basicAuthFunc)

	// start server
	http.ListenAndServe(portString, server)
}

func setupServer() *martini.ClassicMartini {
	router := martini.NewRouter()
	server := martini.New()
	server.Use(martini.Recovery())
	server.Use(requestLogger)
	server.MapTo(router, (*martini.Routes)(nil))
	server.Action(router.Handle)
	return &martini.ClassicMartini{server, router}
}

func requestLogger(res http.ResponseWriter, req *http.Request, c martini.Context) {
	if skipLogging[req.URL.Path] {
		return
	}

	start := time.Now()

	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = req.RemoteAddr
		}
	}

	logger.Printf("Started %s %s for %s", req.Method, req.URL.Path, addr)

	rw := res.(martini.ResponseWriter)
	c.Next()

	logger.Printf("Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
}
