package version

import (
	"os"
	"path"
)

var (
	// BranchString is set by -ldflags in the Makefile - it contains the branch
	// at the time the project was built.
	BranchString string
	// VersionString is set by -ldflags in the Makefile - it contains the
	// latest version (tag or short rev) at the time the project was built.
	VersionString string
	// RevString is set by -ldflags in the Makefile - it contains the full git
	// revision at the time the project was built.
	RevString string
)

/*
The Version struct contains data that is set when building via `-ldflags` in
the Makefile.  The struct member names indicate what data is included in this
list.
*/
type Version struct {
	Branch      string
	Rev         string
	Programname string
	Version     string
}

/*
NewVersion returns a Version instance with values set based on the `-ldflags`
and some sensible defaults.
*/
func NewVersion() *Version {
	ver := &Version{
		Programname: path.Base(os.Args[0]),
	}

	if BranchString == "" {
		ver.Branch = "<unknown>"
	} else {
		ver.Branch = BranchString
	}

	if RevString == "" {
		ver.Rev = "<unknown>"
	} else {
		ver.Rev = RevString
	}

	if VersionString == "" {
		ver.Version = ""
	} else {
		ver.Version = VersionString
	}

	return ver
}
