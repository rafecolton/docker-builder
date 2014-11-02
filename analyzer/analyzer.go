package analyzer

import (
	"github.com/modcloth/go-fileutils"
	"github.com/rafecolton/docker-builder/builderfile"

	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var gitRemoteRegex = regexp.MustCompile("^([^\t\n\f\r ]+)[\t\n\v\f\r ]+(git@github\\.com:|http[s]?:\\/\\/github\\.com\\/)([a-zA-Z0-9]{1}[a-zA-Z0-9-]*)\\/([a-zA-Z0-9_.-]+)\\.git.*$")

/*
An Analysis offers functions that provide data about a given directory. This is
then used to populate an example Bobfile for `builder init .` commands.
*/
type Analysis interface {
	GitRemotes() string
	DockerfilePresent() bool
	IsGitRepo() bool
	RepoBasename() string
}

/*
ParseAnalysisFromDir is a handy function that combines NewAnalysis with
ParseAnalysis to make things a little easier.
*/
func ParseAnalysisFromDir(dir string) (*builderfile.Builderfile, error) {
	if dir == "" {
		dir = "."
	}

	a, err := NewAnalysis(dir)
	if err != nil {
		return nil, err
	}

	b, err := ParseAnalysis(a)
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
NewAnalysis creates an Analysis of the provided directory.
*/
func NewAnalysis(dir string) (Analysis, error) {
	//get absolute path
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// make sure the dir exists
	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("provided dir (%q) does not exist", dir)
		}

		return nil, err
	}

	// make sure the dir is a directory
	if !info.IsDir() {
		return nil, errors.New("provided repo dir must be a directory")
	}

	return &RepoAnalysis{
		repoDir: abs,
	}, nil
}

/*
A RepoAnalysis implements Analysis and returns real data about a given
directory.  It is also the type that is returned by NewAnalysis()
*/
type RepoAnalysis struct {
	repoDir             string
	gitRemotes          string
	gitRemotesPopulated bool
}

func (ra *RepoAnalysis) populateGitRemotes() {
	git, err := fileutils.Which("git")
	if err != nil {
		ra.gitRemotes = ""
	}

	cmd := &exec.Cmd{
		Path: git,
		Args: []string{"git", "remote", "-v"},
		Dir:  ra.repoDir,
	}

	out, err := cmd.Output()
	if err != nil {
		ra.gitRemotes = ""
	} else {
		ra.gitRemotes = string(out)
	}

	ra.gitRemotesPopulated = true
}

/*
GitRemotes returns the output of `git remote -v` when run on the directory being
analyized.  If the remotes cannot be determined (i.e. if the directory is not a
git repo), an empty string is returned.
*/
func (ra *RepoAnalysis) GitRemotes() string {
	if !ra.gitRemotesPopulated {
		ra.populateGitRemotes()
	}

	return ra.gitRemotes
}

/*
DockerfilePresent whether or not a file named Dockerfile is present at the top
level of the directory being analyzed.
*/
func (ra *RepoAnalysis) DockerfilePresent() bool {
	dockerfilePath := filepath.Join(ra.repoDir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); err != nil {
		return false
	}

	return true
}

/*
IsGitRepo returns whether or not the directory being analyzed appears to be a
valid git repo. Validity is determined by whether or not `git remote -v`
returns a non-zero exit code.
*/
func (ra *RepoAnalysis) IsGitRepo() bool {
	if !ra.gitRemotesPopulated {
		ra.populateGitRemotes()
	}

	return ra.gitRemotes != ""
}

/*
RepoBasename returns the basename of the repo being analyzed.
*/
func (ra *RepoAnalysis) RepoBasename() string {
	return filepath.Base(ra.repoDir)
}

/*
ParseAnalysis takes the results of the analysis of a directory and produces a
Builderfile with some educated guesses.  This is later written to a file named
"Bobfile" upon running `builder init .`
*/
func ParseAnalysis(analysis Analysis) (*builderfile.Builderfile, error) {
	if !analysis.DockerfilePresent() {
		return nil, errors.New("uh-oh, can't initialize without a Dockerfile")
	}

	ret := &builderfile.Builderfile{
		Version: 1,
		Docker: *&builderfile.Docker{
			TagOpts: []string{"--force"},
		},
		ContainerArr: []*builderfile.ContainerSection{},
	}

	var appContainer *builderfile.ContainerSection

	if analysis.IsGitRepo() {
		// get registry
		appContainer = &builderfile.ContainerSection{
			Name:       "app",
			Registry:   registryFromRemotes(analysis.GitRemotes()),
			Dockerfile: "Dockerfile",
			SkipPush:   false,
			Project:    analysis.RepoBasename(),
			Tags: []string{
				"git:branch",
				"git:sha",
				"git:tag",
				"latest",
			},
		}
	} else {
		appContainer = &builderfile.ContainerSection{
			Name:       "app",
			Registry:   "my-registry",
			Dockerfile: "Dockerfile",
			SkipPush:   false,
			Project:    analysis.RepoBasename(),
			Tags:       []string{"latest"},
		}
	}

	ret.ContainerArr = append(ret.ContainerArr, appContainer)

	return ret, nil
}

func registryFromRemotes(remotes string) string {
	lines := strings.Split(remotes, "\n")

	var ret string

	for _, line := range lines {
		matches := gitRemoteRegex.FindStringSubmatch(line)
		if len(matches) == 5 {
			ret = matches[3]
			if matches[1] == "origin" {
				break
			}
		}
	}

	return ret
}
