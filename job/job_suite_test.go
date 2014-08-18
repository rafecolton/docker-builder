package job_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rafecolton/docker-builder/job"
	"testing"

	"net/http/httptest"

	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"

	"github.com/rafecolton/docker-builder/server/webhook"
)

var recorder *httptest.ResponseRecorder
var testServer *martini.ClassicMartini

func init() {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	testServer = &martini.ClassicMartini{m, r}
	l := &logrus.Logger{
		Level:     logrus.PanicLevel,
		Formatter: &logrus.JSONFormatter{},
	}
	Logger(l)
	webhook.Logger(l)

	// job control routes
	testServer.Group("/jobs", func(r martini.Router) {
		r.Get("/:id", Get)
		r.Get("/:id/tail", TailN)
		r.Post("", webhook.DockerBuild)
		r.Get("", GetAll)
	})
}

func TestBuilder(t *testing.T) {
	TestMode = true
	RegisterFailHandler(Fail)
	RunSpecs(t, "Job & Job Spec Specs")
}
