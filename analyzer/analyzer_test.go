package analyzer_test

import (
	. "github.com/modcloth/docker-builder/analyzer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

import (
	"github.com/modcloth/docker-builder/builderfile"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Analyzer Specs")
}

var _ = Describe("Analysis Parsing", func() {
	var (
		subject *SpecRepoAnalysis
		outfile *builderfile.Builderfile
	)

	BeforeEach(func() {
		subject = &SpecRepoAnalysis{
			remotes: `origin	git@github.com:modcloth/bob.git (fetch)
					  origin	git@github.com:modcloth/bob.git (push)`,
			dockerfilePresent: true,
			isGitRepo:         true,
			repoBasename:      "fake-repo",
		}
		outfile = &builderfile.Builderfile{
			Docker: *&builderfile.Docker{
				BuildOpts: []string{"--rm", "--no-cache"},
				TagOpts:   []string{"--force"},
			},
			Containers: map[string]builderfile.ContainerSection{
				"global": *&builderfile.ContainerSection{
					Registry: "modcloth",
					Project:  "fake-repo",
					Tags: []string{
						"git:branch",
						"git:rev",
						"git:short",
						"latest",
					},
				},
				"app": *&builderfile.ContainerSection{
					Dockerfile: "Dockerfile",
					SkipPush:   false,
				},
			},
		}
	})

	Context("when given valid data", func() {
		It("correctly parses the repo analysis results", func() {
			out, err := ParseAnalysis(subject)

			Expect(out).To(Equal(outfile))
			Expect(err).To(BeNil())
		})
	})

	Context("when no Dockerfile is present", func() {
		It("produces an error", func() {
			subject.dockerfilePresent = false
			out, err := ParseAnalysis(subject)

			Expect(out).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Context("when the given directory is not a git repo", func() {
		It("only has `latest` tag and default registry", func() {
			subject.isGitRepo = false
			subject.remotes = ""
			outfile.Containers["global"] = *&builderfile.ContainerSection{
				Registry: "my-registry",
				Project:  "fake-repo",
				Tags:     []string{"latest"},
			}
			out, err := ParseAnalysis(subject)

			Expect(out).To(Equal(outfile))
			Expect(err).To(BeNil())

		})
	})
})
