package job

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/docker-builder/parser/uuid"

	"github.com/Sirupsen/logrus"
	"github.com/modcloth/go-fileutils"
	"github.com/modcloth/kamino"
)

const defaultTail = "100"

var gen = uuid.NewUUIDGenerator(true)

//KeepLogTimeInSeconds is the number of seconds to wait before deleting a job's
//logfile and marking it "archived".
var KeepLogTimeInSeconds = 600

type job struct {
	Workdir        string         `json:"-"`
	GitHubAPIToken string         `json:"-"`
	GitCloneDepth  string         `json:"clond_depth,omitempty"`
	Account        string         `json:"account,omitempty"`
	Ref            string         `json:"ref,omitempty"`
	Repo           string         `json:"repo,omitempty"`
	Logger         *logrus.Logger `json:"-"`
	ID             string         `json:"id,omitempty"`
	infoRoute      string         `json:"-"`
	LogRoute       string         `json:"log_route,omitempty"`
	Status         string         `json:"status"`
	logDir         string         `json:"-"`
	logFile        *os.File       `json:"-"`
	Created        time.Time      `json:"created"`
	Completed      time.Time      `json:"completed,omitempty"`
	Archived       time.Time      `json:"archived,omitempty"`
}

/*
JobConfig contains global configuration options for jobs
*/
type JobConfig struct {
	Workdir        string
	Logger         *logrus.Logger
	GitHubAPIToken string
}

/*
NewJob creates a new job from the config as well as a job spec.  After creating
the job, calling job.Process() will actually perform the work.
*/
func NewJob(cfg *JobConfig, spec *JobSpec) *job {
	id, err := gen.NextUUID()
	if err != nil {
		cfg.Logger.WithField("error", err).Error("error creating uuid")
	}

	ret := &job{
		ID:             id,
		Account:        spec.RepoOwner,
		GitHubAPIToken: spec.GitHubAPIToken,
		Ref:            spec.GitRef,
		Repo:           spec.RepoName,
		Workdir:        cfg.Workdir,
		infoRoute:      fmt.Sprintf("/jobs/%s", id),
		LogRoute:       fmt.Sprintf("/jobs/%s/tail?n=%s", id, defaultTail),
		logDir:         fmt.Sprintf("%s/%s", cfg.Workdir, id),
		Status:         "created",
		Created:        time.Now(),
	}

	out := io.MultiWriter(os.Stdout)

	if err = fileutils.MkdirP(ret.logDir, 0755); err != nil {
		cfg.Logger.WithField("error", err).Error("error creating log dir")
		id = ""
	} else {
		file, err := os.Create(fmt.Sprintf("%s/log.log", ret.logDir))
		if err != nil {
			cfg.Logger.WithField("error", err).Error("error creating log file")
			id = ""
		} else {
			out = io.MultiWriter(os.Stdout, file)
			ret.logFile = file
		}
	}

	l := &logrus.Logger{
		Formatter: cfg.Logger.Formatter,
		Level:     cfg.Logger.Level,
		Out:       out,
	}

	ret.Logger = l

	if ret.GitHubAPIToken == "" {
		ret.GitHubAPIToken = cfg.GitHubAPIToken
	}

	if id != "" {
		jobs[id] = ret
	}

	return ret
}

func (job *job) clone() (string, error) {
	job.Logger.WithFields(logrus.Fields{
		"api_token_present":  job.GitHubAPIToken != "",
		"account":            job.Account,
		"ref":                job.Ref,
		"repo":               job.Repo,
		"clone_cache_option": kamino.No,
	}).Info("starting clone process")

	genome := &kamino.Genome{
		APIToken: job.GitHubAPIToken,
		Account:  job.Account,
		Ref:      job.Ref,
		Repo:     job.Repo,
		UseCache: kamino.No,
	}

	factory, err := kamino.NewCloneFactory(job.Workdir)
	if err != nil {
		job.Logger.WithFields(logrus.Fields{
			"api_token_present":  job.GitHubAPIToken != "",
			"account":            job.Account,
			"ref":                job.Ref,
			"repo":               job.Repo,
			"clone_cache_option": kamino.No,
			"error":              err,
		}).Error("issue creating clone factory")

		return "", err
	}

	job.Logger.Debug("attempting to clone")

	path, err := factory.Clone(genome)
	if err != nil {
		job.Logger.WithFields(logrus.Fields{
			"api_token_present":  job.GitHubAPIToken != "",
			"account":            job.Account,
			"ref":                job.Ref,
			"repo":               job.Repo,
			"clone_cache_option": kamino.No,
			"error":              err,
		}).Error("issue cloning")

		return "", err
	}

	job.Logger.WithFields(logrus.Fields{
		"api_token_present":  job.GitHubAPIToken != "",
		"account":            job.Account,
		"ref":                job.Ref,
		"repo":               job.Repo,
		"clone_cache_option": kamino.No,
	}).Info("cloning successful")

	return path, nil
}

func (job *job) build(file string) error {

	job.Logger.Debug("attempting to create a builder")

	bob, err := builder.NewBuilder(job.Logger, true)
	if err != nil {
		job.Logger.WithField("error", err).Error("issue creating a builder")
		return err
	}

	job.Logger.WithField("file", file).Info("building from file")

	return bob.BuildFromFile(file)
}

/*
Process does the actual job processing work, including:

	1. clone the repo
	2. build from the Bobfile at the top level
	3. clean up the cloned repo
*/
func (job *job) Process() error {
	defer func() {
		if job.logFile != nil {
			job.logFile.Close()
		}
	}()

	job.Status = "cloning"
	// step 1: clone
	path, err := job.clone()
	if err != nil {
		job.Status = "errored"
		return err
	}

	job.Status = "building"
	// step 2: build
	if err = job.build(filepath.Join(path, "Bobfile")); err != nil {
		job.Status = "errored"
		return err
	}
	job.Status = "completed"
	job.Completed = time.Now()

	go func() {
		time.Sleep(time.Duration(KeepLogTimeInSeconds) * time.Second)
		job.Status = "archived"
		job.Archived = time.Now()
		job.LogRoute = ""
		fileutils.Rm(fmt.Sprintf("%s/log.log", job.logDir))
	}()

	// step 3: cleanup
	fileutils.RmRF(path)

	return nil
}
