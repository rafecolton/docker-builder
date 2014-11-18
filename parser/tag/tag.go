package tag

import (
	"github.com/sylphon/build-runner/git"
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
Tag, for a special set of macros (currently `git:branch`, `git:sha`,
& `git:tag`) returns git information from the directory in which bob was run.
These macros are specified in args["tag"], and to work properly, args["top"]
must be supplied as well.  If any of the conditions are not met, Tag returns
"".
*/
func (gt *gitTag) Tag() string {
	var top = gt.top

	switch gt.tag {
	case "git:branch":
		return git.Branch(top)
	case "git:rev", "git:sha":
		return git.Sha(top)
	case "git:short", "git:tag":
		return git.Tag(top)
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
