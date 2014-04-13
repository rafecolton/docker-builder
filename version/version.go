package version

import (
	"fmt"
	"os"
	"path"
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
}

func Init() *VersionTrick {
	return &VersionTrick{
		branchString:      BranchString,
		revString:         RevString,
		versionString:     VersionString,
		programnameString: path.Base(os.Args[0]),
	}
}

func (me *VersionTrick) VersionAndExit() {
	if me.versionString == "" {
		me.versionString = "<unknown>"
	}
	fmt.Println(me.versionString)
	os.Exit(0)
}

func (me *VersionTrick) VersionFullAndExit() {
	if me.versionString == "" {
		me.versionString = "<unknown>"
	}
	fmt.Printf("%s, %s\n", me.programnameString, me.versionString)
	os.Exit(0)
}

func (me *VersionTrick) RevAndExit() {
	if me.revString == "" {
		me.revString = "<unknown>"
	}
	fmt.Println(me.revString)
	os.Exit(0)
}

func (me *VersionTrick) BranchAndExit() {
	if me.branchString == "" {
		me.branchString = "<unknown>"
	}
	fmt.Println(me.branchString)
	os.Exit(0)
}
