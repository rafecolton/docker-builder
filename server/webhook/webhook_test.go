package webhook_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rafecolton/docker-builder/server/webhook"

	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
)

func newTestServer() (testServer *martini.ClassicMartini) {
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
	return
}

func init() {
	l := &logrus.Logger{Level: logrus.PanicLevel}
	Logger(l)
}

func TestMain(t *testing.T) {
	TestMode(true)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Specs")
}
