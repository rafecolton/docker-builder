package webhook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
)

func makeJSONRequest(rawBody string) (*http.Request, error) {
	body := []byte(rawBody)

	req, err := http.NewRequest(
		"POST",
		"http://localhost:5000/docker-build",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Length", strconv.Itoa(len(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

var _ = Describe("DockerBuild", func() {
	var (
		validBody     = `{"account": "foo"}`
		validBodySync = `{"account": "foo", "sync":true}`
		invalidBody   = `{"account": 5}`
	)

	Context("when JSON data is invalid", func() {
		It("returns an error", func() {
			var testServer = newTestServer()
			var recorder = httptest.NewRecorder()
			req, err := makeJSONRequest(invalidBody)
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(400))
			Expect(recorder.Body.String()).To(Equal("400 bad request"))
		})
	})

	Context("when JSON data is valid", func() {
		It("returns a 202", func() {
			var testServer = newTestServer()
			var recorder = httptest.NewRecorder()
			req, err := makeJSONRequest(validBody)
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(202))
			Expect(recorder.Body.String()).To(Equal("202 accepted"))
		})
		It("returns a 201 if sync == true", func() {
			var testServer = newTestServer()
			var recorder = httptest.NewRecorder()
			req, err := makeJSONRequest(validBodySync)
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(201))
			Expect(recorder.Body.String()).To(Equal("201 created"))
		})
	})
})
