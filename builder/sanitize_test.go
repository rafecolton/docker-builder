package builder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rafecolton/docker-builder/builder"

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

	Context("when the path is bogus", func() {
		It(`returns an error when the path contains ".."`, func() {
			config, _ := NewTrustedFilePath(dotDotPath, "..")
			_, err := SanitizeTrustedFilePath(config)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(DotDotSanitizeErrorMessage))
		})

		It("returns an error when the path contains symlinks", func() {
			config, _ := NewTrustedFilePath(symlinkPath, "..")
			_, err := SanitizeTrustedFilePath(config)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(SymlinkSanitizeErrorMessage))
		})

		It("returns an error when the path is invalid", func() {
			config, _ := NewTrustedFilePath(bogusPath, "..")
			_, err := SanitizeTrustedFilePath(config)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(InvalidPathSanitizeErrorMessage))
		})
	})

	Context("when the path is valid", func() {
		It("does not return an error", func() {
			config, _ := NewTrustedFilePath(validPath, "..")
			_, err := SanitizeTrustedFilePath(config)
			Expect(err).To(BeNil())
		})

		It("returns a cleaned version of the path", func() {
			config, _ := NewTrustedFilePath(validPath, "..")
			path, _ := SanitizeTrustedFilePath(config)
			Expect(path).To(Equal(cleanedValidPath))
		})
	})
})
