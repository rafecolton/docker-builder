package webhook_test

import (
	//. "github.com/modcloth/docker-builder/webhook"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
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

const (
	// hash of "username/repositoryTRAVIS_TOKEN"
	travisAuthHeader = "92a14eac5ccb1caf98a3623b60d284d77a6233cd20bb82b92a01fa53f7f58dd6"
)

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
	req.Header.Add("Travis-Repo-Slug",
		path.Join(options.Repository.Owner, options.Repository.Name))
	req.Header.Add("Authorization", travisAuthHeader)

	return req, nil
}

var _ = Describe("Travis", func() {

	var ()
	Context("when travis request is invalid", func() {
		It("returns an error when status == 1", func() {
			req, err := makeTravisRequest(&travisRequest{
				Status: 1,
				Repository: travisRepo{
					Owner: "username",
					Name:  "repository",
				},
			}, "")
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(409))
			Expect(recorder.Body.String()).To(Equal("409 conflict"))
		})

		It("returns an error when type is pull_request", func() {
			req, err := makeTravisRequest(&travisRequest{
				Type: "pull_request",
				Repository: travisRepo{
					Owner: "username",
					Name:  "repository",
				},
			}, "")
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(409))
			Expect(recorder.Body.String()).To(Equal("409 conflict"))
		})

		It("returns an error when JSON is invalid", func() {
			req, err := makeTravisRequest(&travisRequest{
				Repository: travisRepo{
					Owner: "username",
					Name:  "repository",
				},
			}, `[this is not valid json}`)

			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(400))
			Expect(recorder.Body.String()).To(Equal("400 bad request"))
		})
		It("returns an error when authentication is incorrent", func() {
			req, err := makeTravisRequest(&travisRequest{
				Repository: travisRepo{
					Owner: "wrong_username",
					Name:  "wrong_repo_name",
				},
			}, "")
			Expect(err).To(BeNil())
			Expect(req).ToNot(BeNil())

			testServer.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(401))
			Expect(recorder.Body.String()).To(Equal("Not Authorized\n"))
		})
	})
	Context("when travis request is valid", func() {
		It("returns success", func() {
			req, err := makeTravisRequest(&travisRequest{
				Type:   "push",
				Commit: "a427f16faa8e4d63f9fcaa4ec55e80765fd11b04",
				Repository: travisRepo{
					Owner: "username",
					Name:  "repository",
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
