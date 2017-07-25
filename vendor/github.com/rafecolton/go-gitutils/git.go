package git

import (
	"bytes"
	"regexp"
	"strings"
)

// Status is a type for reporting the status of a git repo
type Status uint8

const (
	// StatusUpToDate means the local repo matches origin
	StatusUpToDate Status = iota

	// StatusNeedToPull means the local repo needs to pull from the remote
	StatusNeedToPull

	// StatusNeedToPush means the local repo needs to be pushed to the remote
	StatusNeedToPush

	// StatusDiverged means the local repo has diverged from the remote
	StatusDiverged
)

func (s Status) String() string {
	switch s {
	case StatusUpToDate:
		return "StatusUpToDate"
	case StatusNeedToPull:
		return "StatusNeedToPull"
	case StatusNeedToPush:
		return "StatusNeedToPush"
	case StatusDiverged:
		return "StatusDiverged"
	default:
		panic("invalid status")
	}
}

// GitRemoteRegex is a regex for pulling account owner from the output of `git remote -v`
var GitRemoteRegex = regexp.MustCompile(`^([^\t\n\f\r ]+)[\t\n\v\f\r ]+(git@github\.com:|(http[s]?|git):\/\/github\.com\/)([a-zA-Z0-9]{1}[a-zA-Z0-9-]*)\/([a-zA-Z0-9_.-]+)(\.git|[^\t\n\f\r ])+.*$`)

var runner commandRunner

/*
Branch determines the git branch in the repo located at `top`.  Two attempts
are made to determine branch. First:

  git rev-parse -q --abbrev-ref HEAD

If the current working directory is not on a branch, the result will return
"HEAD". In this case, branch will be estimated by parsing the output of the
following:

  git branch -ar --contains $(git rev-parse -q HEAD)
*/
func Branch(top string) string {
	initializeRunner()
	branchBytes, err := runner.BranchCommand(top)
	if err != nil {
		return ""
	}
	branch := strings.TrimRight(string(branchBytes), "\n")
	if branch == "HEAD" {
		branchBytes, err := runner.BranchCommand2(top)
		if err != nil {
			return branch
		}
		branches := strings.Split(strings.TrimRight(string(branchBytes), "\n"), "\n")
		for _, branchStr := range branches {
			if len(branchStr) < 1 || string(branchStr[0]) == "*" {
				continue
			}
			sections := strings.Split(strings.Trim(branchStr, " \n"), "/")
			return sections[len(sections)-1]
		}
	}
	return branch
}

// Sha produces the git branch at `top` as determined by `git rev-parse -q HEAD`
func Sha(top string) string {
	initializeRunner()
	shaBytes, err := runner.ShaCommand(top)
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(shaBytes), "\n")
}

// Tag produces the git tag at `top` as determined by `git describe --always --dirty --tags`
func Tag(top string) string {
	initializeRunner()
	shortBytes, err := runner.TagCommand(top)
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(shortBytes), "\n")
}

// IsClean returns true `git diff --shortstat` produces no output
func IsClean(top string) bool {
	initializeRunner()
	outBytes, err := runner.CleanCommand(top)
	if err != nil || len(outBytes) > 0 {
		return false
	}
	return true
}

//UpToDate returns the status of the repo as determined by the above constants
func UpToDate(top string) Status {
	initializeRunner()
	local, err := runner.UpToDateLocal(top)
	if err != nil {
		return StatusDiverged
	}

	remote, err := runner.UpToDateRemote(top)
	if err != nil {
		return StatusDiverged
	}

	if bytes.Compare(local, remote) == 0 {
		return StatusUpToDate
	}

	base, err := runner.UpToDateBase(top)
	if err != nil {
		return StatusDiverged
	}

	if bytes.Compare(local, base) == 0 {
		return StatusNeedToPull
	} else if bytes.Compare(remote, base) == 0 {
		return StatusNeedToPush
	}
	return StatusDiverged
}

/*
RemoteAccount returns the github account as determined by the output of `git
remote -v`
*/
func RemoteAccount(top string) string {
	initializeRunner()
	remotes, err := runner.RemoteV(top)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(remotes), "\n")

	for _, line := range lines {
		matches := GitRemoteRegex.FindStringSubmatch(line)
		if len(matches) == 7 && matches[1] == "origin" {
			return matches[4]
		}
	}

	return ""
}

func initializeRunner() {
	if runner == nil {
		runner = &realRunner{}
	}
}
