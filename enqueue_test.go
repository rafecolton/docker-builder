package main_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-martini/martini"

	. "github.com/rafecolton/docker-builder"
	"github.com/rafecolton/docker-builder/git"
	"github.com/rafecolton/docker-builder/server"
	"github.com/rafecolton/docker-builder/server/webhook"
)

func init() {
	webhook.TestMode(true)
}

func testServer() (testServer *martini.ClassicMartini) {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	testServer = &martini.ClassicMartini{m, r}
	testServer.Post(server.BuildRoute, webhook.DockerBuild)
	return
}

var enqueuerHost = "http://foo:bar@docker-build-server.example.com:5000"

func testEnqueuer() *Enqueuer {
	return NewEnqueuer(EnqueueOptions{
		Bobfile: "Bobfile.foo",
		Host:    enqueuerHost,
		Top:     os.Getenv("PWD"),
	})
}

func TestEnqueueRequest(t *testing.T) {
	var testServer = testServer()
	var recorder = httptest.NewRecorder()
	var enqueuer = testEnqueuer()

	req, err := enqueuer.Request()
	if err != nil {
		t.Errorf("error making request: %q", err.Error())
	}

	testServer.ServeHTTP(recorder, req)
	if recorder.Code != webhook.AsyncSuccessCode {
		t.Errorf("expected response code %d, got %d", webhook.AsyncSuccessCode, recorder.Code)
	}
	if recorder.Body.String() != webhook.AsyncSuccessMessage {
		t.Errorf("expected response message %q, got %q", webhook.AsyncSuccessMessage, recorder.Body.String())
	}
}

func TestEnqueueRequestBody(t *testing.T) {
	var enqueuer = testEnqueuer()
	var expectedBody = fmt.Sprintf(`{"account":"rafecolton","bobfile":"Bobfile.foo","ref":"%s","repo":"docker-builder"}`, git.Branch(os.Getenv("PWD")))
	bodyBytes, err := enqueuer.BodyBytes()
	if err != nil {
		t.Error(err.Error())
	}
	if string(bodyBytes) != expectedBody {
		t.Errorf("expected request body %s, got %s", expectedBody, string(bodyBytes))
	}
}

func TestEnqueueRequestPath(t *testing.T) {
	var enqueuer = testEnqueuer()
	if enqueuer.RequestPath() != enqueuerHost+"/docker-build" {
		t.Errorf("expected request host %q, got %q", enqueuerHost+"/docker-build", enqueuer.RequestPath())
	}
}
