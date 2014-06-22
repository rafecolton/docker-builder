package job

import (
	"path/filepath"

	"github.com/modcloth/docker-builder/builder"
	"github.com/modcloth/go-fileutils"
	"github.com/modcloth/kamino"

	"github.com/Sirupsen/logrus"
)

type job struct {
	Workdir        string
	GitHubAPIToken string
	GitCloneDepth  string
	Account        string
	Ref            string
	Repo           string
	Logger         *logrus.Logger
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
	ret := &job{
		Account:        spec.RepoOwner,
		GitHubAPIToken: spec.GitHubAPIToken,
		Ref:            spec.GitRef,
		Repo:           spec.RepoName,
		Workdir:        cfg.Workdir,
		Logger:         cfg.Logger,
	}

	if ret.GitHubAPIToken == "" {
		ret.GitHubAPIToken = cfg.GitHubAPIToken
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
	// step 1: clone
	path, err := job.clone()
	if err != nil {
		return err
	}

	// step 2: build
	if err = job.build(filepath.Join(path, "Bobfile")); err != nil {
		return err
	}

	// step 3: cleanup
	return fileutils.RmRF(path)
}
