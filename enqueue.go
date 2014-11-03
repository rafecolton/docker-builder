package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rafecolton/docker-builder/git"

	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

type RequestBody map[string]string

func enqueue(c *cli.Context) {
	var top = os.Getenv("PWD")
	var bobfile = c.Args().First()
	if bobfile == "" {
		bobfile = "Bobfile"
	}

	if !git.IsClean(top) {
		Logger.Error("cannot enqueue, working directory is dirty")
		gocleanup.Exit(1)
	}

	upToDate := git.UpToDate(top)
	if upToDate != 0 {
		switch upToDate {
		case 1:
			Logger.Warn("CAUTION: need to pull")
		case 2:
			Logger.Warn("CAUTION: need to push")
		case 3:
			Logger.Error("cannot enqueue, status has diverged from remote")
			gocleanup.Exit(1)
		}
	}

	//var host = os.Getenv("DOCKER_BUILDER_HOST") + "/jobs"
	var host = c.String("host") + "/jobs"
	var body = RequestBody(map[string]string{
		"account": git.RemoteAccount(top),
		"repo":    git.Repo(top),
		"ref":     git.Branch(top),
		"bobfile": bobfile,
	})
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		Logger.Errorf("error marshaling body to json: %q", err.Error())
		gocleanup.Exit(1)
	}
	resp, err := http.Post(host, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		Logger.Errorf("post error: %q", err.Error())
		gocleanup.Exit(1)
	}
	contentBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(contentBytes))
	gocleanup.Exit(0)
}
