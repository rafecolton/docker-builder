package git

import (
	"os/exec"
)

type commandRunner interface {
	BranchCommand(string) ([]byte, error)
	BranchCommand2(string) ([]byte, error)
	ShaCommand(string) ([]byte, error)
	TagCommand(string) ([]byte, error)
	CleanCommand(string) ([]byte, error)
	UpToDateLocal(string) ([]byte, error)
	UpToDateRemote(string) ([]byte, error)
	UpToDateBase(string) ([]byte, error)
	RemoteV(string) ([]byte, error)
}

type realRunner struct{}
type fakeRunner struct {
	branch         string
	branch2        string
	sha            string
	tag            string
	clean          string
	upToDateLocal  string
	upToDateRemote string
	upToDateBase   string
	remoteV        string
}

// branch 1

func (r *realRunner) BranchCommand(top string) ([]byte, error) {
	cmd := exec.Command("git", "rev-parse", "-q", "--abbrev-ref", "HEAD")
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) BranchCommand(top string) ([]byte, error) {
	return []byte(f.branch), nil
}

// branch 2

func (r *realRunner) BranchCommand2(top string) ([]byte, error) {
	cmd := exec.Command("git", "branch", "-ra", "--contains", Sha(top))
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) BranchCommand2(top string) ([]byte, error) {
	return []byte(f.branch2), nil
}

// sha

func (r *realRunner) ShaCommand(top string) ([]byte, error) {
	cmd := exec.Command("git", "rev-parse", "-q", "HEAD")
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) ShaCommand(top string) ([]byte, error) {
	return []byte(f.sha), nil
}

// tag

func (r *realRunner) TagCommand(top string) ([]byte, error) {
	cmd := exec.Command("git", "describe", "--always", "--dirty", "--tags")
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) TagCommand(top string) ([]byte, error) {
	return []byte(f.tag), nil
}

// clean

func (r *realRunner) CleanCommand(top string) ([]byte, error) {
	cmd := exec.Command("git", "diff", "--shortstat")
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) CleanCommand(top string) ([]byte, error) {
	return []byte(f.clean), nil
}

// up to date local

func (r *realRunner) UpToDateLocal(top string) ([]byte, error) {
	cmd := exec.Command("git", "rev-parse", "@")
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) UpToDateLocal(top string) ([]byte, error) {
	return []byte(f.upToDateLocal), nil
}

// up to date remote

func (r *realRunner) UpToDateRemote(top string) ([]byte, error) {
	cmd := exec.Command("git", "rev-parse", "@{u}")
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) UpToDateRemote(top string) ([]byte, error) {
	return []byte(f.upToDateRemote), nil
}

// up to date base

func (r *realRunner) UpToDateBase(top string) ([]byte, error) {
	cmd := exec.Command("git", "merge-base", "@", "@{u}")
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) UpToDateBase(top string) ([]byte, error) {
	return []byte(f.upToDateBase), nil
}

// remote

func (r *realRunner) RemoteV(top string) ([]byte, error) {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = top
	return cmd.Output()
}

func (f *fakeRunner) RemoteV(top string) ([]byte, error) {
	return []byte(f.remoteV), nil
}
