package analyzer

import (
	"github.com/rafecolton/go-gitutils"
	"github.com/winchman/builder-core/unit-config"

	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

/*
An Analysis offers functions that provide data about a given directory. This is
then used to populate an example Bobfile for `builder init .` commands.
*/
type Analysis interface {
	RemoteAccount() string
	DockerfilePresent() bool
	IsGitRepo() bool
	RepoBasename() string
}

/*
ParseAnalysisFromDir is a handy function that combines NewAnalysis with
ParseAnalysis to make things a little easier.
*/
func ParseAnalysisFromDir(dir string) (*unitconfig.UnitConfig, error) {
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

/*
RemoteAccount returns the output of `git remote -v` when run on the directory being
analyized.  If the remotes cannot be determined (i.e. if the directory is not a
git repo), an empty string is returned.
*/
func (ra *RepoAnalysis) RemoteAccount() string {
	return git.RemoteAccount(ra.repoDir)
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
	cmd := exec.Command("git", "rev-parse")
	cmd.Dir = ra.repoDir
	return cmd.Run() == nil
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
func ParseAnalysis(analysis Analysis) (*unitconfig.UnitConfig, error) {
	if !analysis.DockerfilePresent() {
		return nil, errors.New("uh-oh, can't initialize without a Dockerfile")
	}

	ret := &unitconfig.UnitConfig{
		Version: 1,
		Docker: *&unitconfig.Docker{
			TagOpts: []string{"--force"},
		},
		ContainerArr: []*unitconfig.ContainerSection{},
	}

	registry := "my-registry"
	tags := []string{"latest"}
	if analysis.IsGitRepo() {
		registry = analysis.RemoteAccount()
		tags = append(tags, []string{"{{ branch }}", "{{ sha }}", "{{ tag }}"}...)
	}

	appContainer := &unitconfig.ContainerSection{
		Name:       "app",
		Registry:   registry,
		Dockerfile: "Dockerfile",
		SkipPush:   false,
		Project:    analysis.RepoBasename(),
		Tags:       tags,
	}

	ret.ContainerArr = append(ret.ContainerArr, appContainer)

	return ret, nil
}
