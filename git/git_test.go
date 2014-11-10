package git_test

import (
	"github.com/rafecolton/docker-builder/git"
	"testing"
)

func TestRemoteAccountHTTP(t *testing.T) {
	var remoteV = `origin  http://github.com/rafecolton/docker-builder.git (fetch)`
	parseRemoteAccount(t, remoteV)
}

func TestRemoteAccountHTTPS(t *testing.T) {
	var remoteV = `origin  https://github.com/rafecolton/docker-builder.git (fetch)`
	parseRemoteAccount(t, remoteV)
}

func TestRemoteAccountSSH(t *testing.T) {
	var remoteV = `origin  git@github.com:rafecolton/docker-builder.git (fetch)`
	parseRemoteAccount(t, remoteV)
}

func TestRemoteAccountGit(t *testing.T) {
	var remoteV = `origin  git://github.com/rafecolton/docker-builder.git (fetch)`
	parseRemoteAccount(t, remoteV)
}

func parseRemoteAccount(t *testing.T, remoteV string) {
	var expected = "rafecolton"
	var actual = git.AccountFromRemotes(remoteV)
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
