package analyzer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rafecolton/docker-builder/analyzer"
	"testing"

	"github.com/winchman/builder-core/unit-config"
)

type SpecRepoAnalysis struct {
	remoteAccount     string
	dockerfilePresent bool
	isGitRepo         bool
	repoBasename      string
}

func (sra *SpecRepoAnalysis) RemoteAccount() string {
	return sra.remoteAccount
}

func (sra *SpecRepoAnalysis) DockerfilePresent() bool {
	return sra.dockerfilePresent
}

func (sra *SpecRepoAnalysis) IsGitRepo() bool {
	return sra.isGitRepo
}

func (sra *SpecRepoAnalysis) RepoBasename() string {
	return sra.repoBasename
}

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Analyzer Specs")
}

var _ = Describe("Analysis Parsing", func() {
	var (
		subject *SpecRepoAnalysis
		outfile *unitconfig.UnitConfig
	)

	BeforeEach(func() {
		subject = &SpecRepoAnalysis{
			remoteAccount:     "rafecolton",
			dockerfilePresent: true,
			isGitRepo:         true,
			repoBasename:      "fake-repo",
		}
		outfile = &unitconfig.UnitConfig{
			Version: 1,
			Docker: *&unitconfig.Docker{
				TagOpts: []string{"--force"},
			},
			ContainerArr: []*unitconfig.ContainerSection{
				&unitconfig.ContainerSection{
					Name:     "app",
					Registry: "rafecolton",
					Project:  "fake-repo",
					Tags: []string{
						"latest",
						"{{ branch }}",
						"{{ sha }}",
						"{{ tag }}",
					},
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
			outfile.ContainerArr = []*unitconfig.ContainerSection{
				&unitconfig.ContainerSection{
					Name:       "app",
					Registry:   "my-registry",
					Project:    "fake-repo",
					Tags:       []string{"latest"},
					Dockerfile: "Dockerfile",
					SkipPush:   false,
				},
			}
			out, err := ParseAnalysis(subject)

			Expect(out).To(Equal(outfile))
			Expect(err).To(BeNil())

		})
	})
})
