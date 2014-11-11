package webhook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
)

type travisRepo struct {
	Owner string `json:"owner_name"`
	Name  string `json:"name"`
}

type travisRequest struct {
	Type       string     `json:"type"`
	Status     int        `json:"status"`
	Commit     string     `json:"commit"`
	Repository travisRepo `json:"repository"`
}

func makeTravisRequest(options *travisRequest, bodyString string) (*http.Request, error) {
	if bodyString == "" {
		results, err := json.Marshal(options)
		if err != nil {
			return nil, err
		}

		bodyString = string(results)
	}

	data := url.Values{}
	data.Set("payload", string(bodyString))
	body := []byte(data.Encode())

	req, err := http.NewRequest(
		"POST",
		"http://localhost:5000/docker-build/travis",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Length", strconv.Itoa(len(body)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

var _ = Describe("Travis", func() {
	Context("when travis request is invalid", func() {
		It("returns an error when status == 1", func() {
			var testServer = newTestServer()
			var recorder = httptest.NewRecorder()
			req, err := makeTravisRequest(&travisRequest{
				Status: 1,
			}, "")
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(409))
			Expect(recorder.Body.String()).To(Equal("409 conflict"))
		})

		It("returns an error when type is pull_request", func() {
			var testServer = newTestServer()
			var recorder = httptest.NewRecorder()
			req, err := makeTravisRequest(&travisRequest{
				Type: "pull_request",
			}, "")

			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(409))
			Expect(recorder.Body.String()).To(Equal("409 conflict"))
		})

		It("returns an error when JSON is invalid", func() {
			var testServer = newTestServer()
			var recorder = httptest.NewRecorder()
			req, err := makeTravisRequest(nil, `[this is not valid json}`)
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(400))
			Expect(recorder.Body.String()).To(Equal("400 bad request"))
		})
	})

	Context("when travis request is valid", func() {
		It("returns a valid spec", func() {
			var testServer = newTestServer()
			var recorder = httptest.NewRecorder()
			req, err := makeTravisRequest(&travisRequest{
				Type:   "push",
				Commit: "a427f16faa8e4d63f9fcaa4ec55e80765fd11b04",
				Repository: travisRepo{
					Owner: "testuser",
					Name:  "testrepo",
				},
			}, "")
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(202))
			Expect(recorder.Body.String()).To(Equal("202 accepted"))
		})
	})
})
