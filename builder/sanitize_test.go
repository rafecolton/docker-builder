package builder_test

import (
	. "github.com/modcloth/docker-builder/builder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path/filepath"
)

var _ = Describe("Sanitize Builderfile path", func() {

	var (
		dotDotPath       = "Specs/fixtures/../fixtures/repodir/foo/bar/Bobfile"
		symlinkPath      = "Specs/fixtures/repodir/foo/symlink/Bobfile"
		bogusPath        = "foobarbaz"
		validPath        = "Specs/fixtures/repodir/foo/bar/Bobfile"
		absValidPath, _  = filepath.Abs("../" + validPath)
		cleanedValidPath = filepath.Clean(absValidPath)
	)

	BeforeEach(func() {
		os.Chdir("..")
	})

	Context("when the path is bogus", func() {
		It(`returns an error when the path contains ".."`, func() {
			_, err := SanitizeBuilderfilePath(dotDotPath)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(DotDotSanitizeErrorMessage))
		})

		It("returns an error when the path contains symlinks", func() {
			_, err := SanitizeBuilderfilePath(symlinkPath)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(SymlinkSanitizeErrorMessage))
		})

		It("returns an error when the path is invalid", func() {
			_, err := SanitizeBuilderfilePath(bogusPath)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(InvalidPathSanitizeErrorMessage))
		})
	})

	Context("when the path is valid", func() {
		It("does not return an error", func() {
			_, err := SanitizeBuilderfilePath(validPath)
			Expect(err).To(BeNil())
		})

		It("returns a cleaned version of the path", func() {
			path, _ := SanitizeBuilderfilePath(validPath)
			Expect(path).To(Equal(cleanedValidPath))
		})
	})
})
