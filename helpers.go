package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

func exitErr(exitCode int, message string, args interface{}) {
	var fields logrus.Fields

	switch args.(type) {
	case error:
		fields = logrus.Fields{
			"error": args,
		}
	case map[string]interface{}:
		fields = args.(map[string]interface{})
	}

	Logger.WithFields(fields).Error(message)
	gocleanup.Exit(exitCode)
}

func setServerVars(c *cli.Context) {
	/// lowest priority

	// ENV
	un = config.Username
	pwd = config.Password
	apiToken = config.APIToken
	travisToken = config.TravisToken
	githubSecret = config.GitHubSecret

	// command line
	cliUn := c.String("username")
	cliPwd := c.String("password")
	cliAPIToken := c.String("api-token")
	cliTravisToken := c.String("travis-token")
	cliGitHubSecret := c.String("github-secret")

	if cliTravisToken != "" {
		travisToken = cliTravisToken
	}

	if cliGitHubSecret != "" {
		githubSecret = cliGitHubSecret
	}

	// if username passed on command line, use cl one instead
	if cliUn != "" {
		un = cliUn
	}

	// if password passed on command line, used cl one instead
	if cliPwd != "" {
		pwd = cliPwd
	}

	// get api token
	if cliAPIToken != "" {
		apiToken = cliAPIToken
	}

	// get port
	portString = fmt.Sprintf(":%d", c.Int("port"))

	// get skip-push
	skipPush = c.Bool("skip-push") || config.SkipPush

	// check if builds should be async or not
	syncDefault = c.Bool("sync-build") || config.SyncBuild

	// check if should travis
	shouldTravis = !c.Bool("no-travis") && !config.NoTravis

	// check if should github
	shouldGitHub = !c.Bool("no-github") && !config.NoGitHub

	shouldBasicAuth = (un != "" && pwd != "")
	shouldTravisAuth = (travisToken != "")
	shouldGitHubAuth = (githubSecret != "")

	/// highest priority
}
