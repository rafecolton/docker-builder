package version_test

import (
	. "github.com/rafecolton/docker-builder/version"
	"testing"
)

func TestBranch(t *testing.T) {
	BranchString = "bogus-branch"
	var subject = NewVersion()
	if subject.Branch != BranchString {
		t.Errorf("expected %q, got %q", BranchString, subject.Branch)
	}
}

func TestRev(t *testing.T) {
	RevString = "1234567890"
	var subject = NewVersion()
	if subject.Rev != RevString {
		t.Errorf("expected %q, got %q", RevString, subject.Rev)
	}
}

func TestVersion(t *testing.T) {
	VersionString = "12345-test"
	var subject = NewVersion()
	if subject.Version != VersionString {
		t.Errorf("expected %q, got %q", VersionString, subject.Version)
	}
}
