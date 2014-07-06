package job_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/modcloth/docker-builder/job"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func makeRequest(method, path string, body []byte) (req *http.Request, err error) {
	if len(body) != 0 {
		req, err = http.NewRequest(method, fmt.Sprintf("http://localhost:5000/%s", path), bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("http://localhost:5000/%s", path), nil)
	}

	return req, nil
}

var (
	jobID       = "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3"
	validBody   = `{"account": "foo", "repo": "bar", "ref": "baz"}`
	data        = []byte(validBody)
	recorder2   *httptest.ResponseRecorder
	job         = &Job{}
	jobMap      = map[string]Job{}
	expectedJob = &Job{
		Account:  "foo",
		ID:       jobID,
		LogRoute: fmt.Sprintf("/jobs/%s/tail?n=100", jobID),
		Ref:      "baz",
		Repo:     "bar",
		Status:   "created",
	}
)

var _ = Describe("POST /jobs", func() {

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	It("receives the correct response code and JSON", func() {
		post, _ := makeRequest("POST", "jobs", data)
		testServer.ServeHTTP(recorder, post)
		json.Unmarshal(recorder.Body.Bytes(), &job)

		Expect(job.Account).To(Equal(expectedJob.Account))
		Expect(job.ID).To(Equal(expectedJob.ID))
		Expect(job.LogRoute).To(Equal(expectedJob.LogRoute))
		Expect(job.Ref).To(Equal(expectedJob.Ref))
		Expect(job.Repo).To(Equal(expectedJob.Repo))
		Expect(recorder.Code).To(Equal(201))
	})
})

var _ = Describe("GET /jobs", func() {

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder2 = httptest.NewRecorder()
		post, _ := makeRequest("POST", "jobs", data)
		testServer.ServeHTTP(recorder, post)
	})

	It("receives the correct response code and JSON", func() {
		get, _ := makeRequest("GET", "jobs", nil)
		testServer.ServeHTTP(recorder2, get)
		json.Unmarshal(recorder2.Body.Bytes(), &jobMap)
		job := jobMap[jobID]

		Expect(job.Account).To(Equal(expectedJob.Account))
		Expect(job.ID).To(Equal(expectedJob.ID))
		Expect(job.LogRoute).To(Equal(expectedJob.LogRoute))
		Expect(job.Ref).To(Equal(expectedJob.Ref))
		Expect(job.Repo).To(Equal(expectedJob.Repo))
		Expect(recorder2.Code).To(Equal(200))
	})
})

var _ = Describe("GET /jobs/:id", func() {
	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder2 = httptest.NewRecorder()
		post, _ := makeRequest("POST", "jobs", data)
		testServer.ServeHTTP(recorder, post)
	})

	It("receives the correct response code and JSON", func() {
		get, _ := makeRequest("GET", "jobs/035c4ea0-d73b-5bde-7d6f-c806b04f2ec3", nil)
		testServer.ServeHTTP(recorder2, get)
		json.Unmarshal(recorder2.Body.Bytes(), &job)

		Expect(job.Account).To(Equal(expectedJob.Account))
		Expect(job.ID).To(Equal(expectedJob.ID))
		Expect(job.LogRoute).To(Equal(expectedJob.LogRoute))
		Expect(job.Ref).To(Equal(expectedJob.Ref))
		Expect(job.Repo).To(Equal(expectedJob.Repo))
		Expect(recorder2.Code).To(Equal(200))
	})

})
