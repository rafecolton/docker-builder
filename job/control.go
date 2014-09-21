package job

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strconv"

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

	retBytes, err := json.Marshal(job)
	if err != nil {
		return 409, `{"error": "` + err.Error() + `"}`
	}

	return 200, string(retBytes)
}

//GetAll gets all of the jobs as JSON.
func GetAll(params martini.Params, req *http.Request) (int, string) {
	var jobArr = make([]*Job, len(jobs)) //[l]*Job{}

	var count int
	for _, v := range jobs {
		jobArr[count] = v
		count++
	}
	sort.Sort(List(jobArr))

	retBytes, err := json.Marshal(jobArr)
	if err != nil {
		return 409, `{"error": "` + err.Error() + `"}`
	}

	return 200, string(retBytes)
}
