package tag

import (
	"os"
	"os/exec"
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
			top: os.Getenv("PWD"),
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
	git, _ := exec.LookPath("git")

	branchCmd := &exec.Cmd{
		Path: git,
		Dir:  top,
		Args: []string{git, "rev-parse", "-q", "--abbrev-ref", "HEAD"},
	}
	branchBytes, _ := branchCmd.Output()
	revCmd := &exec.Cmd{
		Path: git,
		Dir:  top,
		Args: []string{git, "rev-parse", "-q", "HEAD"},
	}
	revBytes, _ := revCmd.Output()
	shortCmd := &exec.Cmd{
		Path: git,
		Dir:  top,
		Args: []string{git, "describe", "--always"},
	}
	shortBytes, _ := shortCmd.Output()

	// remove trailing newline
	branch := string(branchBytes)[:len(branchBytes)-1]
	rev := string(revBytes)[:len(revBytes)-1]
	short := string(shortBytes)[:len(shortBytes)-1]

	switch gt.tag {
	case "git:branch":
		return branch
	case "git:rev":
		return rev
	case "git:short":
		return short
	default:
		return ""
	}
}

/*
Tag returns the string in args["tag"], which is the string provided as-is in
the config file
*/
func (tag *stringTag) Tag() string {
	return tag.tag
}
