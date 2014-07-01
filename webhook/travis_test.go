package webhook_test

import (
	. "github.com/modcloth/docker-builder/webhook"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"encoding/json"
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

func makeTravisRequest(options *travisRequest) {
	result, err := json.Marshal(options)
	if err != nil {
		Fail(err.Error())
	}
	data := url.Values{}
	data.Set("payload", string(result))
	PostRequest("POST",
		"/docker-builder/travis",
		HandleWebhook(Travis),
		bytes.NewReader([]byte(data.Encode())),
		map[string]string{
			"Content-Length": strconv.Itoa(len(data.Encode())),
			"Content-Type":   "application/x-www-form-urlencoded",
		},
	)
}

var _ = Describe("Travis", func() {
	Context("when travis request is unsupported", func() {
		It("returns an error when status == 1", func() {
			makeTravisRequest(&travisRequest{
				Status: 1,
			})
			Expect(response.Code).To(Equal(400))
			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
		It("returns an error when type == pull_request", func() {
			makeTravisRequest(&travisRequest{
				Type: "pull_request",
			})
			Expect(response.Code).To(Equal(400))
			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
		It("returns an error when JSON is invalid", func() {
			body := `[this is not valid json}`
			data := url.Values{}
			data.Set("payload", string(body))
			PostRequest("POST",
				"/docker-builder/travis",
				HandleWebhook(Travis),
				bytes.NewReader([]byte(data.Encode())),
				map[string]string{
					"Content-Length": strconv.Itoa(len(data.Encode())),
					"Content-Type":   "application/x-www-form-urlencoded",
				},
			)

			Expect(response.Code).To(Equal(400))
			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})
	Context("when travis request is correct", func() {
		It("returns a valid spec", func() {
			makeTravisRequest(&travisRequest{
				Type:   "push",
				Commit: "a427f16faa8e4d63f9fcaa4ec55e80765fd11b04",
				Repository: travisRepo{
					Owner: "testuser",
					Name:  "testrepo",
				},
			})
			Expect(response.Code).To(Equal(202))
			Expect(spec).ToNot(BeNil())
			Expect(err).To(BeNil())
			err = spec.Validate()
			Expect(err).To(BeNil())
		})
	})
})
