package job

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
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
	var jobArr = []*Job{}

	for _, job := range jobs {
		var matches = true

		for attr, value := range req.URL.Query() {
			// `desiredValue` is the string the user wishes the field to match.
			// If multiple values are provided, the first will be used and the
			// rest will be ignored
			var desiredValue = value[0]

			// if the user provided an invalid attribute by which to filter, it
			// is ignored
			if parameterFieldMapping[attr] == "" {
				continue
			}

			// determine the value of the job struct for the requested attribute
			structField := parameterFieldMapping[attr]

			strukt := reflect.ValueOf(job).Elem()
			var actualValue = strukt.FieldByName(structField).String()

			// if it doesn't match, mark this one a dud and move on
			if actualValue != desiredValue {
				matches = false
				break
			}
		}

		// if we reach this point and have not called it a dud yet (i.e. never
		// marked `matches` as `false`), added it to the list we'll be returning
		if matches {
			jobArr = append(jobArr, job)
		}
	}

	// sort the final (possibly shortened) list
	sort.Sort(ByCreatedDescending(jobArr))

	retBytes, err := json.Marshal(jobArr)
	if err != nil {
		return 409, `{"error": "` + err.Error() + `"}`
	}

	return 200, string(retBytes)
}

var parameterFieldMapping = map[string]string{
	"account": "Account",
	"bobfile": "Bobfile",
	"id":      "ID",
	"ref":     "Ref",
	"repo":    "Repo",
	"status":  "Status",
}
