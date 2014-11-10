package parser

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/modcloth/go-fileutils"
	"github.com/rafecolton/docker-builder/builderfile"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Specs")
}

var nullLogger = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Level:     logrus.PanicLevel,
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
			Containers: []builderfile.ContainerSection{
				*&builderfile.ContainerSection{
					Name:       "base",
					Dockerfile: "Dockerfile.base",
					Registry:   "quay.io/modcloth",
					Project:    "style-gallery",
					Tags:       []string{"base"},
					SkipPush:   true,
					CfgUn:      "foo",
					CfgPass:    "bar",
					CfgEmail:   "baz",
				},
				*&builderfile.ContainerSection{
					Name:       "app",
					Dockerfile: "Dockerfile",
					Registry:   "quay.io/modcloth",
					Project:    "style-gallery",
					Tags:       []string{"git:branch", "git:rev", "git:short"},
					SkipPush:   false,
					CfgUn:      "foo",
					CfgPass:    "bar",
					CfgEmail:   "baz",
				},
			},
		}
		expectedBuilderfile = &builderfile.Builderfile{
			Version: 1,
			Docker: *&builderfile.Docker{
				BuildOpts: []string{"--rm", "--no-cache"},
				TagOpts:   []string{"--force"},
			},
			ContainerGlobals: &builderfile.ContainerSection{
				Registry: "quay.io/modcloth",
				Project:  "style-gallery",
				Tags:     []string{"git:branch", "git:rev", "git:short"},
				CfgUn:    "foo",
				CfgPass:  "bar",
				CfgEmail: "baz",
			},
			ContainerArr: []*builderfile.ContainerSection{
				&builderfile.ContainerSection{
					Name:       "base",
					Dockerfile: "Dockerfile.base",
					Registry:   "",
					Project:    "",
					Tags:       []string{"base"},
					SkipPush:   true,
				},
				&builderfile.ContainerSection{
					Name:       "app",
					Dockerfile: "Dockerfile",
					Registry:   "",
					Project:    "",
					Tags:       nil,
					SkipPush:   false,
				},
			},
		}
	)

	BeforeEach(func() {
		top = os.Getenv("PWD") + "/.."
		git, _ := fileutils.Which("git")
		validFile = top + "/Specs/fixtures/bob.toml"
		invalidFile = top + "/Specs/fixtures/foodoesnotexist"
		subject = nil

		// rev
		revCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "rev-parse", "-q", "HEAD"},
		}
		revBytes, _ := revCmd.Output()
		rev = string(revBytes)[:len(revBytes)-1]

		// branch
		branchCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "rev-parse", "-q", "--abbrev-ref", "HEAD"},
		}

		branchBytes, _ := branchCmd.Output()
		branch = string(branchBytes)[:len(branchBytes)-1]
		if branch == "HEAD" {
			branchCmd = exec.Command("git", "branch", "--contains", rev)
			branchCmd.Dir = top
			branchBytes, err := branchCmd.Output()
			if err == nil {
				branches := strings.Split(string(branchBytes), "\n")
				for _, branchStr := range branches {
					if len(branchStr) < 1 || string(branchStr[0]) == "*" {
						continue
					}
					branch = strings.Trim(branchStr, " ")
					break
				}
			}
		}

		// short
		shortCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "describe", "--always", "--dirty", "--tags"},
		}
		shortBytes, _ := shortCmd.Output()
		short = string(shortBytes)[:len(shortBytes)-1]
		expectedCommandSequence = &CommandSequence{
			Commands: []*SubSequence{
				&SubSequence{
					Metadata: &SubSequenceMetadata{
						Name:       "base",
						Dockerfile: "Dockerfile.base",
						UUID:       "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
					},
					SubCommand: []DockerCmd{
						&BuildCmd{
							buildOpts: docker.BuildImageOptions{
								Name:           "quay.io/modcloth/style-gallery:035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
								NoCache:        true,
								RmTmpContainer: true,
								Auth: docker.AuthConfiguration{
									Username: "foo",
									Password: "bar",
									Email:    "baz",
								},
								AuthConfigs: docker.AuthConfigurations{
									Configs: map[string]docker.AuthConfiguration{
										"quay.io/modcloth": docker.AuthConfiguration{
											Username:      "foo",
											Password:      "bar",
											Email:         "baz",
											ServerAddress: "quay.io/modcloth",
										},
									},
								},
								ContextDir: ".",
							},
							origBuildOpts: []string{"--rm", "--no-cache"},
						},
						&TagCmd{Repo: "quay.io/modcloth/style-gallery", Tag: "base", Force: true},
					},
				},
				&SubSequence{
					Metadata: &SubSequenceMetadata{
						Name:       "app",
						Dockerfile: "Dockerfile",
						UUID:       "035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
					},
					SubCommand: []DockerCmd{
						&BuildCmd{
							buildOpts: docker.BuildImageOptions{
								Name:           "quay.io/modcloth/style-gallery:035c4ea0-d73b-5bde-7d6f-c806b04f2ec3",
								NoCache:        true,
								RmTmpContainer: true,
								Auth: docker.AuthConfiguration{
									Username: "foo",
									Password: "bar",
									Email:    "baz",
								},
								AuthConfigs: docker.AuthConfigurations{
									Configs: map[string]docker.AuthConfiguration{
										"quay.io/modcloth": docker.AuthConfiguration{
											Username:      "foo",
											Password:      "bar",
											Email:         "baz",
											ServerAddress: "quay.io/modcloth",
										},
									},
								},
								ContextDir: ".",
							},
							origBuildOpts: []string{"--rm", "--no-cache"},
						},
						&TagCmd{Repo: "quay.io/modcloth/style-gallery", Tag: branch, Force: true},
						&TagCmd{Repo: "quay.io/modcloth/style-gallery", Tag: rev, Force: true},
						&TagCmd{Repo: "quay.io/modcloth/style-gallery", Tag: short, Force: true},
						&PushCmd{
							Image:     "quay.io/modcloth/style-gallery",
							Tag:       branch,
							Registry:  "quay.io/modcloth",
							AuthUn:    "foo",
							AuthPwd:   "bar",
							AuthEmail: "baz",
						},
						&PushCmd{
							Image:     "quay.io/modcloth/style-gallery",
							Tag:       rev,
							Registry:  "quay.io/modcloth",
							AuthUn:    "foo",
							AuthPwd:   "bar",
							AuthEmail: "baz",
						},
						&PushCmd{
							Image:     "quay.io/modcloth/style-gallery",
							Tag:       short,
							Registry:  "quay.io/modcloth",
							AuthUn:    "foo",
							AuthPwd:   "bar",
							AuthEmail: "baz",
						},
					},
				},
			},
		}
	})

	Context("with a valid Builderfile", func() {
		It("returns a non empty string and a nil error as raw data", func() {
			subject := NewParser(validFile, nullLogger)
			raw, err := subject.getRaw()
			Expect(len(raw)).ToNot(Equal(0))
			Expect(err).To(BeNil())
		})

		It("returns a fully parsed Builderfile", func() {
			subject := NewParser(validFile, nullLogger)
			actual, err := subject.rawToStruct()
			Expect(expectedBuilderfile).To(Equal(actual))
			Expect(err).To(BeNil())
		})

		It("further processes the Builderfile into an InstructionSet", func() {
			subject := NewParser(validFile, nullLogger)
			actual, err := subject.structToInstructionSet()
			Expect(expectedInstructionSet).To(Equal(actual))
			Expect(err).To(BeNil())
		})

		It("further processes the InstructionSet into an CommandSequence", func() {
			subject := NewParser(validFile, nullLogger)
			subject.SeedUUIDGenerator()
			actual, err := subject.instructionSetToCommandSequence()
			Expect(expectedCommandSequence).To(Equal(actual))
			Expect(err).To(BeNil())
		})
	})

	Context("with an invalid Builderfile", func() {
		It("returns an empty string and error as raw data", func() {
			subject := NewParser(invalidFile, nullLogger)
			raw, err := subject.getRaw()
			Expect(raw).To(Equal(""))
			Expect(err).ToNot(BeNil())
		})
	})
})
