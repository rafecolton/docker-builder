package tag

import (
	"os/exec"
	"strings"
)

/*
Tag is the interface for specifying tags for container builds.
*/
type Tag interface {
	Tag() string
}

/*
NewTag returns a Tag instance.  See function implementation for details on what
args to pass.
*/
func NewTag(version string, args map[string]string) Tag {
	switch version {
	case "null":
		return &nullTag{}
	case "git":
		return &gitTag{
			tag: args["tag"],
			top: args["top"],
		}
	default:
		return &stringTag{
			tag: args["tag"],
		}
	}
}

// used for empty tags for testing
type nullTag struct {
}

// used for git-based tags
type gitTag struct {
	tag string
	top string
}

// used for "as-is" tags
type stringTag struct {
	tag string
}

/*
Tag returns the fixed string "<TAG>" for a nullTag.
*/
func (tag *nullTag) Tag() string {
	return "<TAG>"
}

/*
Tag, for a special set of macros (currently `git:branch`, `git:rev`,
& `git:short`) returns git information from the directory in which bob was run.
These macros are specified in args["tag"], and to work properly, args["top"]
must be supplied as well.  If any of the conditions are not met, Tag returns
"".
*/
func (gt *gitTag) Tag() string {

	top := gt.top
	var gitexe = &gitexe{top: top}

	switch gt.tag {
	case "git:branch":
		return gitexe.branch()
	case "git:rev", "git:sha":
		return gitexe.sha()
	case "git:short", "git:tag":
		return gitexe.tag()
	default:
		return ""
	}
}

type gitexe struct {
	exe string
	top string
}

func (g gitexe) branch() string {
	var branchCmd = exec.Command("git", "rev-parse", "-q", "--abbrev-ref", "HEAD")
	branchCmd.Dir = g.top
	branchBytes, _ := branchCmd.Output()
	branch := string(branchBytes)[:len(branchBytes)-1]
	if branch == "HEAD" {
		branchCmd = exec.Command("git", "branch", "--contains", g.sha())
		branchCmd.Dir = g.top
		branchBytes, _ := branchCmd.Output()
		branches := strings.Split(string(branchBytes)[:len(branchBytes)-1], "\n")
		for _, branch := range branches {
			if string(branch[0]) == "*" {
				continue
			}
			return strings.Replace(branch, " ", "", -1)
		}
		return branch

	}
	return branch
}

func (g gitexe) sha() string {
	var revCmd = exec.Command("git", "rev-parse", "-q", "HEAD")
	revCmd.Dir = g.top
	revBytes, _ := revCmd.Output()
	rev := string(revBytes)[:len(revBytes)-1]
	return rev
}

func (g gitexe) tag() string {
	var shortCmd = exec.Command("git", "describe", "--always", "--dirty", "--tags")
	shortCmd.Dir = g.top
	shortBytes, _ := shortCmd.Output()
	short := string(shortBytes)[:len(shortBytes)-1]
	return short
}

/*
Tag returns the string in args["tag"], which is the string provided as-is in
the config file
*/
func (tag *stringTag) Tag() string {
	return tag.tag
}
