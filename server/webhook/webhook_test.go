package webhook_test

import (
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/modcloth/docker-builder/server/webhook"

	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
)

var recorder = httptest.NewRecorder()
var testServer *martini.ClassicMartini

func init() {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	testServer = &martini.ClassicMartini{m, r}

	testServer.Post("/docker-build/github", Github)
	testServer.Post("/docker-build/travis", Travis)
	testServer.Post("/docker-build", DockerBuild)

	l := &logrus.Logger{Level: logrus.Panic}

	Logger(l)
}

func TestMain(t *testing.T) {
	TestMode(true)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Specs")
}
