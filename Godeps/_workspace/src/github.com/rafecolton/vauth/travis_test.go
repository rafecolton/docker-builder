package vauth

import (
	"github.com/go-martini/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_TravisAuth(t *testing.T) {
	recorder := httptest.NewRecorder()

	auth := "92a14eac5ccb1caf98a3623b60d284d77a6233cd20bb82b92a01fa53f7f58dd6"

	m := martini.New()
	m.Use(TravisCI("TRAVIS_TOKEN"))
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
	})

	r, _ := http.NewRequest("GET", "foo", nil)
	r.Header.Set("Authorization", auth)
	r.Header.Set("Travis-Repo-Slug", "username/repository")
	m.ServeHTTP(recorder, r)

	if recorder.Code == 401 {
		t.Error("Response is 401")
	}

	if recorder.Body.String() != "hello" {
		t.Error("Auth failed, got: ", recorder.Body.String())
	}
}
