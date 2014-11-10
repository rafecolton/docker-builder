package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/rafecolton/docker-builder/git"
	"github.com/rafecolton/docker-builder/server"

	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

// EnqueueOptions is a struct for sending options to an Enqueuer
type EnqueueOptions struct {
	Bobfile string
	Host    string
	Top     string
}

// Enqueuer is a struct that handles parsing the repo data and making the
// actual enqueue request for the `docker-builder enqueue` feature
type Enqueuer struct {
	account string
	bobfile string
	host    string
	ref     string
	repo    string
}

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
	if upToDate != git.StatusUpToDate {
		switch upToDate {
		case git.StatusNeedToPull:
			Logger.Warn("CAUTION: need to pull")
		case git.StatusNeedToPush:
			Logger.Warn("CAUTION: need to push")
		case git.StatusDiverged:
			Logger.Error("cannot enqueue, status has diverged from remote")
			gocleanup.Exit(1)
		}
	}

	var host = c.String("host")
	opts := EnqueueOptions{
		Host:    host,
		Bobfile: bobfile,
		Top:     top,
	}
	enqueuer := NewEnqueuer(opts)
	result, err := enqueuer.Enqueue()
	if err != nil {
		Logger.Errorf("error enqueueing build: %q", err.Error())
		gocleanup.Exit(1)
	}
	Logger.Debugln(result)
	gocleanup.Exit(0)
}

// NewEnqueuer returns an Enqueuer with data populated from the repo
// information
func NewEnqueuer(options EnqueueOptions) *Enqueuer {
	return &Enqueuer{
		account: git.RemoteAccount(options.Top),
		bobfile: options.Bobfile,
		host:    options.Host,
		ref:     git.Branch(options.Top),
		repo:    git.Repo(options.Top),
	}
}

// BodyBytes returns the byte slice that enc would send in an enqueue request
func (enc *Enqueuer) BodyBytes() ([]byte, error) {
	var body = map[string]string{
		"account": enc.account,
		"repo":    enc.repo,
		"ref":     enc.ref,
		"bobfile": enc.bobfile,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

// RequestPath returns the path to which enqueue requests are sent.  This
// includes both the host and the route
func (enc *Enqueuer) RequestPath() string {
	return enc.host + server.BuildRoute
}

// Request returns the http request that will be sent for enqueueing
func (enc *Enqueuer) Request() (*http.Request, error) {
	bodyBytes, err := enc.BodyBytes()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		enc.RequestPath(),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Add("Content-Length", strconv.Itoa(len(bodyBytes)))
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// Enqueue performs an actual http request using the result of Request().
// Enqueue() should not be called during tests
func (enc *Enqueuer) Enqueue() (string, error) {
	req, err := enc.Request()
	if err != nil {
		return "", err
	}
	reqBody, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewReader(reqBody)) // reset body after reading
	Logger.Debugf("enqueueing request %s", reqBody)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	contentBytes, _ := ioutil.ReadAll(resp.Body)
	return string(contentBytes), nil
}
