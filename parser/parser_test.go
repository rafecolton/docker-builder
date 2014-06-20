package parser

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/modcloth/docker-builder/builderfile"
	"github.com/modcloth/go-fileutils"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Specs")
}

var nullLogger = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Level:     logrus.Panic,
}

var _ = Describe("Parse", func() {

	var (
		subject                 *Parser
		validFile               string
		invalidFile             string
		branch                  string
		rev                     string
		short                   string
		top                     string
		expectedCommandSequence *CommandSequence
		expectedInstructionSet  = &InstructionSet{
			DockerBuildOpts: []string{"--rm", "--no-cache"},
			DockerTagOpts:   []string{"--force"},
			Containers: map[string]builderfile.ContainerSection{
				"base": *&builderfile.ContainerSection{
					Name:       "base",
					Dockerfile: "Dockerfile.base",
					Included:   []string{"Gemfile", "Gemfile.lock"},
					Excluded:   []string{"spec", "tmp"},
					Registry:   "quay.io/modcloth",
					Project:    "style-gallery",
					Tags:       []string{"base"},
					SkipPush:   true,
				},
				"app": *&builderfile.ContainerSection{
					Name:       "app",
					Dockerfile: "Dockerfile",
					Included:   []string{},
					Excluded:   []string{"spec", "tmp"},
					Registry:   "quay.io/modcloth",
					Project:    "style-gallery",
					Tags:       []string{"git:branch", "git:rev", "git:short"},
					SkipPush:   false,
				},
			},
		}
		expectedBuilderfile = &builderfile.Builderfile{
			Docker: *&builderfile.Docker{
				BuildOpts: []string{"--rm", "--no-cache"},
				TagOpts:   []string{"--force"},
			},
			Containers: map[string]builderfile.ContainerSection{
				"global": *&builderfile.ContainerSection{
					Dockerfile: "",
					Included:   nil,
					Excluded:   []string{"spec", "tmp"},
					Registry:   "quay.io/modcloth",
					Project:    "style-gallery",
					Tags:       []string{"git:branch", "git:rev", "git:short"},
					SkipPush:   false,
				},
				"base": *&builderfile.ContainerSection{
					Dockerfile: "Dockerfile.base",
					Included:   []string{"Gemfile", "Gemfile.lock"},
					Excluded:   nil,
					Registry:   "",
					Project:    "",
					Tags:       []string{"base"},
					SkipPush:   true,
				},
				"app": *&builderfile.ContainerSection{
					Dockerfile: "Dockerfile",
					Included:   nil,
					Excluded:   nil,
					Registry:   "",
					Project:    "",
					Tags:       nil,
					SkipPush:   false,
				},
			},
		}
	)

	BeforeEach(func() {
		top = os.Getenv("PWD")
		git, _ := fileutils.Which("git")
		validFile = fmt.Sprintf("%s/spec/fixtures/bob.toml", top)
		invalidFile = fmt.Sprintf("%s/specs/fixtures/foodoesnotexist", top)
		subject = nil
		// branch
		branchCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "rev-parse", "-q", "--abbrev-ref", "HEAD"},
		}

		branchBytes, _ := branchCmd.Output()
		branch = string(branchBytes)[:len(branchBytes)-1]

		// rev
		revCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "rev-parse", "-q", "HEAD"},
		}
		revBytes, _ := revCmd.Output()
		rev = string(revBytes)[:len(revBytes)-1]

		// short
		shortCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "describe", "--always"},
		}
		shortBytes, _ := shortCmd.Output()
		short = string(shortBytes)[:len(shortBytes)-1]
		expectedCommandSequence = &CommandSequence{
			Commands: []*SubSequence{
				&SubSequence{
					Metadata: &SubSequenceMetadata{
						Name:       "base",
						Dockerfile: "Dockerfile.base",
						Excluded:   []string{"spec", "tmp"},
						Included:   []string{"Gemfile", "Gemfile.lock"},
						UUID:       "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
					},
					SubCommand: []exec.Cmd{
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"build",
								"-t",
								"quay.io/modcloth/style-gallery:035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
								"--rm",
								"--no-cache",
								".",
							},
						},
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"tag",
								"--force",
								"<IMG>",
								"quay.io/modcloth/style-gallery:base",
							},
						},
					},
				},
				&SubSequence{
					Metadata: &SubSequenceMetadata{
						Name:       "app",
						Dockerfile: "Dockerfile",
						Excluded:   []string{"spec", "tmp"},
						Included:   []string{},
						UUID:       "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
					},
					SubCommand: []exec.Cmd{
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"build",
								"-t",
								"quay.io/modcloth/style-gallery:035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
								"--rm",
								"--no-cache",
								".",
							},
						},
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"tag",
								"--force",
								"<IMG>",
								fmt.Sprintf("quay.io/modcloth/style-gallery:%s", branch),
							},
						},
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"tag",
								"--force",
								"<IMG>",
								fmt.Sprintf("quay.io/modcloth/style-gallery:%s", rev),
							},
						},
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"tag",
								"--force",
								"<IMG>",
								fmt.Sprintf("quay.io/modcloth/style-gallery:%s", short),
							},
						},
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"push",
								fmt.Sprintf("quay.io/modcloth/style-gallery:%s", branch),
							},
						},
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"push",
								fmt.Sprintf("quay.io/modcloth/style-gallery:%s", rev),
							},
						},
						*&exec.Cmd{
							Path: "docker",
							Args: []string{
								"docker",
								"push",
								fmt.Sprintf("quay.io/modcloth/style-gallery:%s", short),
							},
						},
					},
				},
			},
		}
	})

	Context("with a valid Builderfile", func() {

		It("produces an openable file", func() {
			subject, _ := NewParser(validFile, nullLogger)
			Expect(subject.IsOpenable()).To(Equal(true))
		})

		It("returns a non empty string and a nil error as raw data", func() {
			subject, _ := NewParser(validFile, nullLogger)
			raw, err := subject.getRaw()
			Expect(len(raw)).ToNot(Equal(0))
			Expect(err).To(BeNil())
		})

		It("returns a fully parsed Builderfile", func() {
			subject, _ := NewParser(validFile, nullLogger)
			actual, err := subject.rawToStruct()
			Expect(expectedBuilderfile).To(Equal(actual))
			Expect(err).To(BeNil())
		})

		It("further processes the Builderfile into an InstructionSet", func() {
			subject, _ := NewParser(validFile, nullLogger)
			actual, err := subject.structToInstructionSet()
			Expect(expectedInstructionSet).To(Equal(actual))
			Expect(err).To(BeNil())
		})

		It("further processes the InstructionSet into an CommandSequence", func() {
			subject, _ := NewParser(validFile, nullLogger)
			subject.SeedUUIDGenerator()
			actual, err := subject.instructionSetToCommandSequence()
			Expect(expectedCommandSequence).To(Equal(actual))
			Expect(err).To(BeNil())
		})
	})

	Context("with an invalid Builderfile", func() {
		It("returns an empty string and error as raw data", func() {
			subject, _ := NewParser(invalidFile, nullLogger)
			raw, err := subject.getRaw()
			Expect(raw).To(Equal(""))
			Expect(err).ToNot(BeNil())
		})
	})
})
