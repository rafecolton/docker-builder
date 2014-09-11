package fileutils

import (
	"os"
	"os/exec"
	"path/filepath"
)

// ChmodR is like `chmod -R`
func ChmodR(name string, mode os.FileMode) error {
	return filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chmod(path, mode)
		}
		return err
	})
}

// ChownR is like `chown -R`
func ChownR(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}

// MkdirP is `mkdir -p` / os.MkdirAll
func MkdirP(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Mv is `mv` / os.Rename
func Mv(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// Rm is `rm` / os.Remove
func Rm(name string) error {
	return os.Remove(name)
}

// RmRF is `rm -rf` / os.RemoveAll
func RmRF(path string) error {
	return os.RemoveAll(path)
}

// Which is `which` / exec.LookPath
func Which(file string) (string, error) {
	return exec.LookPath(file)
}
