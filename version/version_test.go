package version_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rafecolton/builder/version"
	color "github.com/wsxiaoys/terminal/color"
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

		subject = New()
	})

	It("prints the correct branch", func() {
		Expect(subject.Branch).To(Equal(color.Sprint("@{!w}bogus-branch")))
	})

	It("prints the correct rev", func() {
		Expect(subject.Rev).To(Equal(color.Sprint("@{!w}1234567890")))
	})

	It("prints the correct version", func() {
		Expect(subject.Version).To(Equal(color.Sprint("@{!w}12345-test")))
	})

	It("prints the correct full version", func() {
		Expect(subject.VersionFull).To(Equal(color.Sprintf("@{!w}%s %s", "version.test", subject.Version)))
	})
})
