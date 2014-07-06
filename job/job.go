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

/*
TestMode monkeys with certain things for tests so bad things don't
happen
*/
var TestMode bool

const defaultTail = "100"

var gen uuid.UUIDGenerator
var logger *logrus.Logger

//KeepLogTimeInSeconds is the number of seconds to wait before deleting a job's
//logfile and marking it "archived".
var KeepLogTimeInSeconds = 600

/*
Job is the struct representation of a build job.  Intended to be
created with NewJob, but exported so it can be used for tests.
*/
type Job struct {
	Account        string         `json:"account,omitempty"`
	Archived       time.Time      `json:"archived,omitempty"`
	Completed      time.Time      `json:"completed,omitempty"`
	Created        time.Time      `json:"created"`
	Error          error          `json:"error,omitempty"`
	GitCloneDepth  string         `json:"clond_depth,omitempty"`
	GitHubAPIToken string         `json:"-"`
	ID             string         `json:"id,omitempty"`
	LogRoute       string         `json:"log_route,omitempty"`
	Logger         *logrus.Logger `json:"-"`
	Ref            string         `json:"ref,omitempty"`
	Repo           string         `json:"repo,omitempty"`
	Status         string         `json:"status"`
	Workdir        string         `json:"-"`
	infoRoute      string         `json:"-"`
	logDir         string         `json:"-"`
	logFile        *os.File       `json:"-"`
}

/*
JobConfig contains global configuration options for jobs
*/
type JobConfig struct {
	Workdir        string
	Logger         *logrus.Logger
	GitHubAPIToken string
}

//Logger sets the (global) logger for the server package
func Logger(l *logrus.Logger) {
	logger = l
}

/*
NewJob creates a new job from the config as well as a job spec.  After creating
the job, calling job.Process() will actually perform the work.
*/
func NewJob(cfg *JobConfig, spec *JobSpec) *Job {
	gen = uuid.NewUUIDGenerator(!TestMode)
	id, err := gen.NextUUID()
	if err != nil {
		cfg.Logger.WithField("error", err).Error("error creating uuid")
	}

	ret := &Job{
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
	}

	if !TestMode {
		ret.Created = time.Now()
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
			if TestMode {
				out = file
			} else {
				out = io.MultiWriter(os.Stdout, file)
			}
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

func (job *Job) clone() (string, error) {
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

func (job *Job) build(file string) error {

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
func (job *Job) Process() error {
	if TestMode {
		return job.processTestMode()
	}

	var archive = func() {
		time.Sleep(time.Duration(KeepLogTimeInSeconds) * time.Second)
		job.Status = "archived"
		job.Archived = time.Now()
		job.LogRoute = ""
		fileutils.Rm(fmt.Sprintf("%s/log.log", job.logDir))
	}

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
		job.Error = err
		go archive()
		return err
	}

	job.Status = "building"
	// step 2: build
	if err = job.build(filepath.Join(path, "Bobfile")); err != nil {
		job.Status = "errored"
		job.Error = err
		go archive()
		return err
	}

	job.Status = "completed"
	job.Completed = time.Now()
	fileutils.RmRF(path)
	go archive()
	return nil
}

func (job *Job) processTestMode() error {
	job.Logger.Warn("job.Process() called in test mode")
	job.Status = "completed"
	job.Completed = time.Now()

	// log something for test purposes
	levelBefore := job.Logger.Level
	job.Logger.Level = logrus.Debug
	job.Logger.Debug("FOO")
	job.Logger.Level = levelBefore

	return nil
}
