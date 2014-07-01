package webhook_test

//import (
//. "github.com/modcloth/docker-builder/webhook"

//. "github.com/onsi/ginkgo"
//. "github.com/onsi/gomega"

//"bytes"
//"encoding/json"
//"strconv"
//)

//type githubOwner struct {
//Name string `json:"name"`
//}

//type githubRepo struct {
//Owner githubOwner `json:"owner"`
//Name  string      `json:"name"`
//}

//type githubRequest struct {
//Commit     string     `json:"after"`
//Repository githubRepo `json:"repository"`
//Event      string     `json:"-"`
//}

//func makeGithubRequest(options *githubRequest) {
//result, err := json.Marshal(options)
//if err != nil {
//Fail(err.Error())
//}
//PostRequest("POST",
//"/docker-builder/github",
//HandleWebhook(Github),
//bytes.NewReader(result),
//map[string]string{
//"Content-Length": strconv.Itoa(len(result)),
//"Content-Type":   "application/json",
//"X-Github-Event": options.Event,
//},
//)
//}

//var _ = Describe("Github", func() {
//Context("when github request is unsupported", func() {
//It("returns an error when event is not push", func() {
//makeGithubRequest(&githubRequest{
//Event: "issue",
//})
//Expect(response.Code).To(Equal(400))
//Expect(spec).To(BeNil())
//Expect(err).ToNot(BeNil())
//})
//It("returns an error when JSON is invalid", func() {
//body := []byte(`[this is not valid json}`)
//PostRequest("POST",
//"/docker-builder/github",
//HandleWebhook(Github),
//bytes.NewReader(body),
//map[string]string{
//"Content-Length": strconv.Itoa(len(body)),
//"Content-Type":   "application/x-www-form-urlencoded",
//"X-Github-Event": "push",
//},
//)

//Expect(response.Code).To(Equal(400))
//Expect(spec).To(BeNil())
//Expect(err).ToNot(BeNil())
//})
//})
//Context("when Github request is correct", func() {
//It("returns a valid spec", func() {
//makeGithubRequest(&githubRequest{
//Event:  "push",
//Commit: "a427f16faa8e4d63f9fcaa4ec55e80765fd11b04",
//Repository: githubRepo{
//Owner: githubOwner{
//Name: "testuser",
//},
//Name: "testrepo",
//},
//})
//Expect(response.Code).To(Equal(202))
//Expect(spec).ToNot(BeNil())
//Expect(err).To(BeNil())
//err = spec.Validate()
//Expect(err).To(BeNil())
//})
//})
//})
