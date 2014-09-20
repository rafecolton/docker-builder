package job

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/pkg/tailfile"
	"github.com/go-martini/martini"
)

var jobs = map[string]*Job{}

/*
TailN is the handler function for the job log tailing route.
*/
func TailN(params martini.Params, req *http.Request) (int, string) {
	id := params["id"]
	n := req.FormValue("n")

	if n == "" {
		n = defaultTail
	}

	intN, err := strconv.Atoi(n)
	if err != nil {
		return 400, n + " is not a valid number"
	}

	out, err := tailN(intN, id)
	if err != nil {
		return 412, err.Error()
	}

	return 200, out
}

func tailN(n int, id string) (string, error) {
	job := jobs[id]
	logFilePath := job.logDir + "/log.log"

	file, err := os.Open(logFilePath)
	if err != nil {
		return "", err
	}

	byteMatrix, err := tailfile.TailFile(file, n)
	if err != nil {
		return "", err
	}

	out := bytes.Join(byteMatrix, []byte("\n"))

	return string(out) + "\n", err
}

//Get gets the requested job as JSON.
func Get(params martini.Params, req *http.Request) (int, string) {
	id := params["id"]
	job := jobs[id]
	job.setLogRouteHost(req)

	retBytes, err := json.Marshal(job)
	if err != nil {
		return 409, `{"error": "` + err.Error() + `"}`
	}

	return 200, string(retBytes)
}

//GetAll gets all of the jobs as JSON.
func GetAll(params martini.Params, req *http.Request) (int, string) {
	for _, job := range jobs {
		job.setLogRouteHost(req)
	}

	retBytes, err := json.Marshal(jobs)
	if err != nil {
		return 409, `{"error": "` + err.Error() + `"}`
	}

	return 200, string(retBytes)
}

func (j *Job) setLogRouteHost(req *http.Request) {
	var host string

	if req.Host != "" {
		host = req.Host
	} else {
		host = req.URL.Host
	}

	var scheme string
	if req.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	if strings.HasPrefix(j.LogRoute, "/") {
		if scheme != "" {
			scheme = scheme + "://"
		}

		j.LogRoute = scheme + host + j.LogRoute
	}
}
