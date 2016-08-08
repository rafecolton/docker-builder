package vauth

import (
	"bytes"
	"github.com/go-martini/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GitHubAuth(t *testing.T) {
	recorder := httptest.NewRecorder()

	secret := "secret"
	signature := "sha1=5d61605c3feea9799210ddcb71307d4ba264225f"
	body := "{}"

	m := martini.New()
	m.Use(GitHub(secret))
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
	})

	r, _ := http.NewRequest("GET", "foo", bytes.NewReader([]byte(body)))
	r.Header.Set("X-Hub-Signature", signature)
	m.ServeHTTP(recorder, r)

	if recorder.Code == 401 {
		t.Error("Response is 401")
	}

	if recorder.Body.String() != "hello" {
		t.Error("Auth failed, got: ", recorder.Body.String())
	}
}
