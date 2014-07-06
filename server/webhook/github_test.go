package webhook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

type githubOwner struct {
	Name string `json:"name"`
}

type githubRepo struct {
	Owner githubOwner `json:"owner"`
	Name  string      `json:"name"`
}

type githubRequest struct {
	Commit     string     `json:"after"`
	Repository githubRepo `json:"repository"`
	Event      string     `json:"-"`
	RawBody    string     `json:"-"`
}

func makeGithubRequest(options *githubRequest) (*http.Request, error) {
	var body []byte
	var err error

	if options.RawBody == "" {
		body, err = json.Marshal(options)
		if err != nil {
			return nil, err
		}
	} else {
		body = []byte(options.RawBody)
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:5000/docker-build/github",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Length", strconv.Itoa(len(body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Github-Event", options.Event)

	return req, nil
}

var _ = Describe("Github", func() {
	Context("when github request is unsupported", func() {
		It("returns an error when event is not push", func() {
			req, err := makeGithubRequest(&githubRequest{
				Event: "issue",
			})
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(400))
			Expect(recorder.Body.String()).To(Equal("400 bad request"))
		})
		It("returns an error when JSON is invalid", func() {
			req, err := makeGithubRequest(&githubRequest{
				RawBody: `[this is not valid json}`,
				Event:   "push",
			})
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(400))
			Expect(recorder.Body.String()).To(Equal("400 bad request"))
		})
	})
	Context("when Github request is correct", func() {
		It("returns a valid spec", func() {
			req, err := makeGithubRequest(&githubRequest{
				Event:  "push",
				Commit: "a427f16faa8e4d63f9fcaa4ec55e80765fd11b04",
				Repository: githubRepo{
					Owner: githubOwner{
						Name: "testuser",
					},
					Name: "testrepo",
				},
			})
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(202))
			Expect(recorder.Body.String()).To(Equal("202 accepted"))
		})
	})
})
