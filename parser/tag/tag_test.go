package tag_test

import (
	. "github.com/rafecolton/docker-builder/parser/tag"

	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/modcloth/go-fileutils"
)

var args = map[string]string{"tag": "foo"}

func Test_StringTagPrintsString(t *testing.T) {
	var subject = NewTag("default", args)
	var expected = "foo"
	var actual = subject.Tag()
	if actual != expected {
		t.Error("expected " + expected + ", got " + actual)
	}
}

var top = os.Getenv("PWD")
var git, _ = fileutils.Which("git")

func getBranch() string {
	var git, err = fileutils.Which("git")
	if err != nil {
		fmt.Println(err)
	}
	branchCmd := &exec.Cmd{
		Path: git,
		Dir:  top,
		Args: []string{git, "rev-parse", "-q", "--abbrev-ref", "HEAD"},
	}
	branchBytes, err := branchCmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	branch := string(branchBytes)[:len(branchBytes)-1]
	if branch == "HEAD" {
		branchCmd = exec.Command("git", "branch", "--contains", getSha())
		branchCmd.Dir = top
		branchBytes, _ := branchCmd.Output()
		branches := strings.Split(string(branchBytes), "\n")
	Loop:
		for _, branchStr := range branches {
			if len(branchStr) > 0 && string(branchStr[0]) != "*" {
				branch = strings.Trim(branchStr, " ")
				break Loop
			}

		}
	}
	return branch
}

func getSha() (sha string) {
	shaCmd := &exec.Cmd{
		Path: git,
		Dir:  top,
		Args: []string{git, "rev-parse", "-q", "HEAD"},
	}
	shaBytes, _ := shaCmd.Output()
	sha = string(shaBytes)[:len(shaBytes)-1]
	return
}

func getTag() (tag string) {
	tagCmd := &exec.Cmd{
		Path: git,
		Dir:  top,
		Args: []string{git, "describe", "--always", "--dirty", "--tags"},
	}
	tagBytes, _ := tagCmd.Output()
	tag = string(tagBytes)[:len(tagBytes)-1]
	return
}

func Test_GitTagBranch(t *testing.T) {
	var subject = NewTag("git", map[string]string{
		"tag": "git:branch",
		"top": top,
	})
	var actual = subject.Tag()
	var expected = getBranch()
	if actual != expected {
		t.Error("expected " + expected + ", got " + actual)
	}
}

func Test_GitTagSha(t *testing.T) {
	var subject = NewTag("git", map[string]string{
		"tag": "git:sha",
		"top": top,
	})
	var actual = subject.Tag()
	var expected = getSha()
	if actual != expected {
		t.Error("expected " + expected + ", got " + actual)
	}
}

func Test_GitTagTag(t *testing.T) {
	var subject = NewTag("git", map[string]string{
		"tag": "git:tag",
		"top": top,
	})
	var actual = subject.Tag()
	var expected = getTag()
	if actual != expected {
		t.Error("expected " + expected + ", got " + actual)
	}
}

func Test_GitTagNull(t *testing.T) {
	var subject = NewTag("null", nil)
	var actual = subject.Tag()
	var expected = "<TAG>"
	if actual != expected {
		t.Error("expected " + expected + ", got " + actual)
	}
}
