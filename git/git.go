package git

import (
	"os/exec"
	"strings"
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
