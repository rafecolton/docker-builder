package job

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/modcloth/go-fileutils"
	"github.com/modcloth/kamino"
	gouuid "github.com/nu7hatch/gouuid"
	"github.com/winchman/builder-core"
	"github.com/winchman/builder-core/unit-config"

	"github.com/rafecolton/docker-builder/conf"
)

const (
	defaultTail         = "100"
	defaultBobfile      = "Bobfile"
	specFixturesRepoDir = "./_testing/fixtures/repodir"
)

var (
	// TestMode monkeys with certain things for tests so bad things don't happen
	TestMode bool

	// SkipPush indicates whether or not a global --skip-push directive has been given
	SkipPush bool

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
	InfoRoute          string         `json:"info_route,omitempty"`
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

func (job *Job) addHostToRoutes(req *http.Request) {
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

	if strings.HasPrefix(job.LogRoute, "/") {
		job.LogRoute = scheme + "://" + host + job.LogRoute
	}

	if strings.HasPrefix(job.InfoRoute, "/") {
		job.InfoRoute = scheme + "://" + host + job.InfoRoute
	}
}

/*
NewJob creates a new job from the config as well as a job spec.  After creating
the job, calling job.Process() will actually perform the work.
*/
func NewJob(cfg *Config, spec *Spec, req *http.Request) *Job {
	var idUUID *gouuid.UUID
	var err error
	idUUID, err = gouuid.NewV4()
	if TestMode {
		idUUID, err = gouuid.NewV5(gouuid.NamespaceURL, []byte("0"))
	}

	if err != nil {
		cfg.Logger.WithField("error", err).Error("error creating uuid")
	}
	id := idUUID.String()

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
		InfoRoute:      "/jobs/" + id,
		LogRoute:       "/jobs/" + id + "/tail?n=" + defaultTail,
		logDir:         cfg.Workdir + "/" + id,
		Status:         "created",
		Created:        time.Now(),
	}
	ret.addHostToRoutes(req)

	out, file, err := newMultiWriter(ret.logDir)
	if err != nil {
		cfg.Logger.WithField("error", err).Error("error creating log dir")
		id = ""
	}

	ret.logFile = file
	ret.Logger = &logrus.Logger{
		Formatter: cfg.Logger.Formatter,
		Level:     cfg.Logger.Level,
		Out:       out,
	}

	if ret.GitHubAPIToken == "" {
		ret.GitHubAPIToken = cfg.GitHubAPIToken
	}

	if id != "" {
		jobs[id] = ret
	}

	return ret
}

func newMultiWriter(logDir string) (io.Writer, *os.File, error) {
	var out io.Writer

	if err := fileutils.MkdirP(logDir, 0755); err != nil {
		return nil, nil, err
	}

	file, err := os.Create(logDir + "/log.log")
	if err != nil {
		return nil, nil, err
	}

	if TestMode {
		out = file
	} else {
		out = io.MultiWriter(os.Stdout, file)
	}

	return out, file, nil
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
	unitConfig, err := unitconfig.ReadFromFile(job.clonedRepoLocation + "/" + job.Bobfile)
	if err != nil {
		job.Logger.WithField("error", err).Error("issue parsing Bobfile")
		return err
	}

	globals := unitconfig.ConfigGlobals{
		SkipPush: SkipPush,
		CfgUn:    conf.Config.CfgUn,
		CfgPass:  conf.Config.CfgPass,
		CfgEmail: conf.Config.CfgEmail,
	}

	unitConfig.SetGlobals(globals)

	job.Logger.WithField("file", job.Bobfile).Info("building from file")

	return runner.RunBuildSynchronously(runner.Options{
		UnitConfig: unitConfig,
		ContextDir: job.clonedRepoLocation,
		LogLevel:   job.Logger.Level,
	})
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
		job.Logger.WithField("error", err).Error("unable to process job synchronously")
		job.Status = "errored"
		job.Error = err
		return err
	}
	job.clonedRepoLocation = path

	// step 2: build
	job.Status = "building"
	if err = job.build(); err != nil {
		job.Logger.WithField("error", err).Error("unable to process job synchronously")
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
