package builder

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path"
	"testing"
)

import . "github.com/rafecolton/builder/version"

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Version Specs")
}

var _ = Describe("Version", func() {
	var (
		subject *VersionTrick
	)

	BeforeEach(func() {
		subject = &VersionTrick{
			BranchString:      "bogus-branch",
			RevString:         "",
			ProgramnameString: path.Base(os.Args[0]),
			VersionString:     "builder 12345",
		}
	})

	It("succeeds", func() {
		Expect("foo").To(Equal("foo"))
	})

	XIt("fails", func() {
		Expect("foo").To(Equal("bar"))
	})
})
