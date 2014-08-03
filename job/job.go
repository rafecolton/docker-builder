package job

import (
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

const (
	defaultTail         = "100"
	defaultBobfile      = "Bobfile"
	specFixturesRepoDir = "./Specs/fixtures/repodir"
)

var (
	// TestMode monkeys with certain things for tests so bad things don't happen
	TestMode bool

	gen    uuid.Generator
	logger *logrus.Logger
)

/*
Job is the struct representation of a build job.  Intended to be
created with NewJob, but exported so it can be used for tests.
*/
type Job struct {
	Account            string         `json:"account,omitempty"`
	Bobfile            string         `json:"bobfile,omitempty"`
	Completed          time.Time      `json:"completed,omitempty"`
	Created            time.Time      `json:"created"`
	Error              error          `json:"error,omitempty"`
	GitCloneDepth      string         `json:"clone_depth,omitempty"`
	GitHubAPIToken     string         `json:"-"`
	ID                 string         `json:"id,omitempty"`
	LogRoute           string         `json:"log_route,omitempty"`
	Logger             *logrus.Logger `json:"-"`
	Ref                string         `json:"ref,omitempty"`
	Repo               string         `json:"repo,omitempty"`
	Status             string         `json:"status"`
	Workdir            string         `json:"-"`
	infoRoute          string         `json:"-"`
	logDir             string         `json:"-"`
	logFile            *os.File       `json:"-"`
	clonedRepoLocation string         `json:"-"`
}

/*
Config contains global configuration options for jobs
*/
type Config struct {
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
func NewJob(cfg *Config, spec *Spec) *Job {
	gen = uuid.NewUUIDGenerator(!TestMode)
	id, err := gen.NextUUID()
	if err != nil {
		cfg.Logger.WithField("error", err).Error("error creating uuid")
	}

	bobfile := spec.Bobfile
	if bobfile == "" {
		bobfile = defaultBobfile
	}

	ret := &Job{
		Bobfile:        bobfile,
		ID:             id,
		Account:        spec.RepoOwner,
		GitHubAPIToken: spec.GitHubAPIToken,
		Ref:            spec.GitRef,
		Repo:           spec.RepoName,
		Workdir:        cfg.Workdir,
		infoRoute:      "/jobs/" + id,
		LogRoute:       "/jobs/" + id + "/tail?n=" + defaultTail,
		logDir:         cfg.Workdir + "/" + id,
		Status:         "created",
		Created:        time.Now(),
	}

	out := io.MultiWriter(os.Stdout)

	if err = fileutils.MkdirP(ret.logDir, 0755); err != nil {
		cfg.Logger.WithField("error", err).Error("error creating log dir")
		id = ""
	} else {
		file, err := os.Create(ret.logDir + "/log.log")
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

func (job *Job) build() error {

	job.Logger.Debug("attempting to create a builder")
	bobfile := filepath.Join(job.clonedRepoLocation, job.Bobfile)

	bob, err := builder.NewBuilder(job.Logger, true)
	if err != nil {
		job.Logger.WithField("error", err).Error("issue creating a builder")
		return err
	}

	job.Logger.WithField("file", bobfile).Info("building from file")

	err = bob.BuildFromFile(bobfile)

	return err
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

	defer func() {
		if job.logFile != nil {
			job.logFile.Close()
		}
	}()

	// step 1: clone
	job.Status = "cloning"
	path, err := job.clone()
	if err != nil {
		job.Status = "errored"
		job.Error = err
		return err
	}
	job.clonedRepoLocation = path

	// step 2: build
	job.Status = "building"
	if err = job.build(); err != nil {
		job.Status = "errored"
		job.Error = err
		return err
	}

	job.Status = "completed"
	job.Completed = time.Now()
	fileutils.RmRF(path)
	return nil
}

func (job *Job) processTestMode() error {
	// If this function is used correctly,
	// we should never see this warning message.
	job.Logger.Warn("processing job in test mode")

	// set status to validating in anticipation of performing validation step
	job.Status = "validating"

	// set clone path to fixtures dir
	job.clonedRepoLocation = specFixturesRepoDir

	// log something for test purposes
	levelBefore := job.Logger.Level
	job.Logger.Level = logrus.DebugLevel
	job.Logger.Debug("FOO")
	job.Logger.Level = levelBefore

	// mark job as completed
	job.Status = "completed"
	job.Completed = time.Now()

	return nil
}
