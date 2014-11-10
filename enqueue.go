package main

import (
	"bufio"
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

type EnqueueOptions struct {
	Bobfile string
	Host    string
	Top     string
}

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

func NewEnqueuer(options EnqueueOptions) *Enqueuer {
	return &Enqueuer{
		account: git.RemoteAccount(options.Top),
		bobfile: options.Bobfile,
		host:    options.Host,
		ref:     git.Branch(options.Top),
		repo:    git.Repo(options.Top),
	}
}

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

func (enc *Enqueuer) RequestPath() string {
	return enc.host + server.BuildRoute
}

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

func (enc *Enqueuer) Enqueue() (string, error) {
	req, err := enc.Request()
	if err != nil {
		return "", err
	}
	reqBody, _ := ioutil.ReadAll(req.Body)
	Logger.Debugf("enqueueing request %s", reqBody)
	resp, err := http.ReadResponse(bufio.NewReader(nil), req)
	defer resp.Body.Close()
	contentBytes, _ := ioutil.ReadAll(resp.Body)
	return string(contentBytes), nil
}
