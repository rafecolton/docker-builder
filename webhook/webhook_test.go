package webhook_test

import (
	//"fmt"
	//"io"
	//"net/http"
	//"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//"github.com/modcloth/docker-builder/job"
	. "github.com/modcloth/docker-builder/webhook"

	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
	"github.com/Sirupsen/logrus"
)

//var (
//response *httptest.ResponseRecorder
//spec     *job.JobSpec
//err      error
//)

func init() {
	Logger(logrus.New())
}

func TestMain(t *testing.T) {
	TestMode(true)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Specs")
}

//func PostRequest(method string, route string, handler martini.Handler, body io.Reader, extraHeaders map[string]string) {
//m := martini.Classic()
//m.Post(route, handler)
//m.Use(render.Renderer())
//request, _ := http.NewRequest(method, route, body)
//for key, value := range extraHeaders {
//request.Header.Add(key, value)
//}
//response = httptest.NewRecorder()
//m.ServeHTTP(response, request)
//}

//func HandleWebhook(parseJobSpec Parser) func(http.ResponseWriter, *http.Request) (int, string) {
//return func(w http.ResponseWriter, req *http.Request) (int, string) {
//spec, err = parseJobSpec(req)
//if err != nil {
//fmt.Println(err)
//return 400, "400 bad request"
//}
//return 202, "202 accepted"
//}
//}
