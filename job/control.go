package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/go-martini/martini"
)

var jobs = map[string]*job{}

/*
TailN is the handler function for the job log tailing route.
*/
func TailN(params martini.Params, req *http.Request) (int, string) {
	id := params["id"]
	n := req.FormValue("n")

	if n == "" {
		n = defaultTail
	}

	nInt, err := strconv.Atoi(n)
	if err != nil {
		return 400, fmt.Sprintf("%s is not a valid number", n)
	}

	out, err := tailN(nInt, id)
	if err != nil {
		return 412, err.Error()
	}

	return 200, out
}

func tailN(n int, id string) (string, error) {
	job := jobs[id]

	if job.Status == "archived" {
		return "", errors.New("job is archived, no log file available")
	}

	logFilePath := fmt.Sprintf("%s/log.log", job.logDir)
	out, err := exec.Command("tail", "-n", fmt.Sprintf("%d", n), logFilePath).Output()

	return string(out), err
}

//Get gets the requested job as JSON.
func Get(params martini.Params, req *http.Request) (int, string) {
	id := params["id"]
	job := jobs[id]

	retBytes, err := json.Marshal(job)
	if err != nil {
		return 409, fmt.Sprintf(`{"error": %q}`, err)
	}

	return 200, string(retBytes)
}

//GetAll gets all of the jobs as JSON.
func GetAll(params martini.Params) (int, string) {
	retBytes, err := json.Marshal(jobs)
	if err != nil {
		return 409, fmt.Sprintf(`{"error": %q}`, err)
	}

	return 200, string(retBytes)
}
