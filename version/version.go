package version

import (
	"os"
	"path"
)

import (
	. "github.com/rafecolton/builder/config"
	. "github.com/wsxiaoys/terminal/color"
)

var (
	BranchString  string
	VersionString string
	RevString     string
	branchFlag    = false
	versionFlag   = false
	revFlag       = false
)

type VersionTrick struct {
	branchString      string
	revString         string
	versionString     string
	programnameString string
	runtime           Runtime
	BranchString      string
	RevString         string
	ProgramnameString string
	VersionString     string
}

func Init(runtime *Runtime) *VersionTrick {
	return &VersionTrick{
		branchString:      BranchString,
		revString:         RevString,
		versionString:     VersionString,
		programnameString: path.Base(os.Args[0]),
		runtime:           *runtime,
	}
}

func (me *VersionTrick) VersionAndExit() {
	if me.versionString == "" {
		me.versionString = "<unknown>"
	}
	me.runtime.Println(Sprint("@{!w}" + me.versionString))
	os.Exit(0)
}

func (me *VersionTrick) VersionFullAndExit() {
	if me.versionString == "" {
		me.versionString = "<unknown>"
	}
	me.runtime.Println(Sprintf("@{!w}%s, %s", me.programnameString, me.versionString))
	os.Exit(0)
}

func (me *VersionTrick) RevAndExit() {
	if me.revString == "" {
		me.revString = "<unknown>"
	}
	me.runtime.Println(Sprint("@{!w}" + me.revString))
	os.Exit(0)
}

func (me *VersionTrick) BranchAndExit() {
	if me.branchString == "" {
		me.branchString = "<unknown>"
	}
	me.runtime.Println(Sprint("@{!w}" + me.branchString))
	os.Exit(0)
}
