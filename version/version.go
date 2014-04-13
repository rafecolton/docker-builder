package version

import (
	"os"
	"path"
)

import (
	color "github.com/wsxiaoys/terminal/color"
)

var (
	BranchString  string
	VersionString string
	RevString     string
)

type Version struct {
	Branch      string
	Rev         string
	Programname string
	Version     string
	VersionFull string
}

func New() *Version {
	ver := &Version{
		Programname: path.Base(os.Args[0]),
	}

	if BranchString == "" {
		ver.Branch = "<unknown>"
	} else {
		ver.Branch = color.Sprintf("@{!w}%s", BranchString)
	}

	if RevString == "" {
		ver.Rev = "<unknown>"
	} else {
		ver.Rev = color.Sprintf("@{!w}%s", RevString)
	}

	if VersionString == "" {
		ver.Version = ""
		ver.VersionFull = ""
	} else {
		ver.Version = color.Sprintf("@{!w}%s", VersionString)
		ver.VersionFull = color.Sprintf("@{!w}%s %s", ver.Programname, ver.Version)
	}

	return ver
}
