package main_test

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	cmd := exec.Cmd{
		Path: os.Getenv("GOPATH") + "/bin/docker-builder",
		Args: []string{"docker-builder", "serve"},
	}
	if err := cmd.Start(); err != nil {
		t.Errorf("error start server: %s\n", err.Error())
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(cmd.Process.Pid, syscall.SIGUSR1)
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}()

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 166 {
					t.Errorf("server exited unexpectedly and/or did not start correctly")
				}
			} else {
				t.Errorf("unable to get process exit code")
			}
		} else {
			t.Errorf("error waiting on server: %s\n", err.Error())
		}
	}
}

func TestServerWithBasicAuth(t *testing.T) {
	cmd := exec.Cmd{
		Path: os.Getenv("GOPATH") + "/bin/docker-builder",
		Args: []string{"docker-builder", "serve", "--username", "foo", "--password", "bar"},
	}
	if err := cmd.Start(); err != nil {
		t.Errorf("error start server: %s\n", err.Error())
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(cmd.Process.Pid, syscall.SIGUSR1)
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}()

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 166 {
					t.Errorf("server exited unexpectedly and/or did not start correctly")
				}
			} else {
				t.Errorf("unable to get process exit code")
			}
		} else {
			t.Errorf("error waiting on server: %s\n", err.Error())
		}
	}
}
