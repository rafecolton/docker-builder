package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rafecolton/docker-builder/analyzer"
)

func Branch(top string) string {
	var branchCmd = exec.Command("git", "rev-parse", "-q", "--abbrev-ref", "HEAD")
	branchCmd.Dir = top
	branchBytes, err := branchCmd.Output()
	if err != nil {
		return ""
	}
	branch := string(branchBytes)[:len(branchBytes)-1]
	if branch == "HEAD" {
		branchCmd = exec.Command("git", "branch", "--contains", Sha(top))
		branchCmd.Dir = top
		branchBytes, err := branchCmd.Output()
		if err != nil {
			return branch
		}
		branches := strings.Split(string(branchBytes)[:len(branchBytes)-1], "\n")
		for _, branchStr := range branches {
			if len(branchStr) > 1 || string(branchStr[0]) == "*" {
				continue
			}
			return strings.Trim(branchStr, " ")
		}
	}
	return branch
}

func Sha(top string) string {
	var revCmd = exec.Command("git", "rev-parse", "-q", "HEAD")
	revCmd.Dir = top
	revBytes, err := revCmd.Output()
	if err != nil {
		return ""
	}
	rev := string(revBytes)[:len(revBytes)-1]
	return rev
}

func Tag(top string) string {

	var shortCmd = exec.Command("git", "describe", "--always", "--dirty", "--tags")
	shortCmd.Dir = top
	shortBytes, err := shortCmd.Output()
	if err != nil {
		return ""
	}
	short := string(shortBytes)[:len(shortBytes)-1]
	return short
}

func IsClean(top string) bool {
	var cmd = exec.Command("git", "diff", "--shortstat")
	cmd.Dir = top
	outBytes, err := cmd.Output()
	if err != nil {
		return false
	}
	if len(outBytes) > 0 {
		return false
	}
	return true
}

/*
0 - up to date
1 - need to pull
2 - need to push
3 - diverged or error status (considered not up to date)
*/
func UpToDate(top string) int {
	var cmdLocal = exec.Command("git", "rev-parse", "@")
	cmdLocal.Dir = top
	local, err := runCmd(cmdLocal)
	if err != nil {
		return 3
	}

	var cmdRemote = exec.Command("git", "rev-parse", "@{u}")
	cmdRemote.Dir = top
	remote, err := runCmd(cmdRemote)
	if err != nil {
		return 3
	}

	if local == remote {
		return 0
	}

	var cmdBase = exec.Command("git", "merge-base", "@", "@{u}")
	cmdBase.Dir = top
	base, err := runCmd(cmdBase)
	if err != nil {
		return 3
	}
	if local == base {
		return 1
	} else if remote == base {
		return 2
	}
	return 3
}

func RemoteAccount(top string) string {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = top
	outBytes, err := cmd.Output()
	if err != nil {
		return ""
	}
	remotes := string(outBytes)
	lines := strings.Split(remotes, "\n")

	var ret string

	for _, line := range lines {
		matches := analyzer.GitRemoteRegex.FindStringSubmatch(line)
		if len(matches) == 5 {
			ret = matches[3]
			if matches[1] == "origin" {
				break
			}
		}
	}

	return ret
}

func Repo(top string) string {
	return filepath.Base(os.Getenv("PWD"))
}

///////// HELPERS

func runCmd(cmd *exec.Cmd) (string, error) {
	outBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(outBytes)[:len(outBytes)-1], nil
}
