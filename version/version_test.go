package version_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rafecolton/docker-builder/version"
	"testing"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Version Specs")
}

var _ = Describe("Version", func() {
	var (
		subject *Version
	)

	BeforeEach(func() {
		BranchString = "bogus-branch"
		RevString = "1234567890"
		VersionString = "12345-test"

		subject = NewVersion()
	})

	It("prints the correct branch", func() {
		Expect(subject.Branch).To(Equal("bogus-branch"))
	})

	It("prints the correct rev", func() {
		Expect(subject.Rev).To(Equal("1234567890"))
	})

	It("prints the correct version", func() {
		Expect(subject.Version).To(Equal("12345-test"))
	})
})
