package tag

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

import (
	"os"
	"os/exec"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tag Specs")
}

var _ = Describe("String Tag", func() {
	var (
		subject Tag
		args    map[string]string
	)

	BeforeEach(func() {
		args = map[string]string{
			"tag": "foo",
		}
		subject = NewTag("default", args)
	})

	Context("when providing a string tag", func() {
		It("prints out the string provided", func() {
			Expect(subject.Tag()).To(Equal("foo"))
		})
	})
})

var _ = Describe("Git Tag", func() {
	var (
		subject Tag
		branch  string
		rev     string
		short   string
		top     string
	)

	BeforeEach(func() {
		top = os.ExpandEnv("${PWD}")
		git, _ := exec.LookPath("git")
		subject = NewTag("git", map[string]string{
			"top": top,
			"tag": "foo",
		})

		// branch
		branchCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "rev-parse", "-q", "--abbrev-ref", "HEAD"},
		}

		branchBytes, _ := branchCmd.Output()
		branch = string(branchBytes)

		// rev
		revCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "rev-parse", "-q", "HEAD"},
		}
		revBytes, _ := revCmd.Output()
		rev = string(revBytes)

		// short
		shortCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "describe", "--always"},
		}
		shortBytes, _ := shortCmd.Output()
		short = string(shortBytes)
	})

	Context("parsing git macros", func() {
		It("translates `git:branch` correctly", func() {
			subject = NewTag("git", map[string]string{
				"tag": "git:branch",
				"top": top,
			})
			Expect(subject.Tag()).To(Equal(branch))
		})

		It("translates `git:rev` correctly", func() {
			subject = NewTag("git", map[string]string{
				"tag": "git:rev",
				"top": top,
			})
			Expect(subject.Tag()).To(Equal(rev))
		})

		It("translates `git:short` correctly", func() {
			subject = NewTag("git", map[string]string{
				"tag": "git:short",
				"top": top,
			})
			Expect(subject.Tag()).To(Equal(short))
		})
	})
})

var _ = Describe("Null Tag", func() {
	var (
		subject Tag
	)

	BeforeEach(func() {
		subject = NewTag("null", nil)
	})

	It("prints a pre-determined fixed tag", func() {
		Expect(subject.Tag()).To(Equal("<TAG>"))
	})
})
