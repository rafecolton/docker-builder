package tag

import (
	"github.com/libgit2/git2go"
	"os"
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
			top: os.ExpandEnv("${PWD}"),
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
	repo, err := git.OpenRepository(gt.top)
	if err != nil {
		return ""
	}

	ref, err := repo.LookupReference("HEAD")
	if err != nil {
		return ""
	}
	ref, err = ref.Resolve()
	if err != nil {
		return ""
	}

	switch gt.tag {
	case "git:branch":
		return ref.Shorthand()
	case "git:rev":
		return ref.Target().String()
	case "git:short":
		return ref.Target().String()[0:7]
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
