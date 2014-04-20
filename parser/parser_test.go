package parser

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

import (
	"fmt"
	"github.com/rafecolton/bob/builderfile"
	"github.com/rafecolton/bob/log"
	"os"
	"os/exec"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Specs")
}

var _ = Describe("Parse", func() {

	var (
		subject                 *Parser
		validFile               string
		invalidFile             string
		top                     = os.ExpandEnv("${PWD}")
		expectedCommandSequence = &CommandSequence{
			commands: []exec.Cmd{
				*&exec.Cmd{
					Path:   "docker",
					Args:   []string{"docker", "build", "-t", "quay.io/modcloth/style-gallery:035c4ea0-d73b-5bde-7d6f-c806b04f2ec3", "--rm", "--no-cache"},
					Stdout: nil,
					Stderr: nil,
				},
				*&exec.Cmd{
					Path:   "docker",
					Args:   []string{"docker", "tag", "<IMG>", "quay.io/modcloth/style-gallery:base"},
					Stdout: nil,
					Stderr: nil,
				},
				*&exec.Cmd{
					Path:   "docker",
					Args:   []string{"docker", "push", "quay.io/modcloth/style-gallery"},
					Stdout: nil,
					Stderr: nil,
				},
				*&exec.Cmd{
					Path:   "docker",
					Args:   []string{"docker", "build", "-t", "quay.io/modcloth/style-gallery:latest", "--rm", "--no-cache"},
					Stdout: nil,
					Stderr: nil,
				},
				*&exec.Cmd{
					Path:   "docker",
					Args:   []string{"docker", "tag", "<IMG>", "quay.io/modcloth/style-gallery:035c4ea0-d73b-5bde-7d6f-c806b04f2ec3"},
					Stdout: nil,
					Stderr: nil,
				},
				*&exec.Cmd{
					Path:   "docker",
					Args:   []string{"docker", "tag", "<IMG>", "quay.io/modcloth/style-gallery<TAG>:"},
					Stdout: nil,
					Stderr: nil,
				},
				*&exec.Cmd{
					Path:   "docker",
					Args:   []string{"docker", "tag", "<IMG>", "quay.io/modcloth/style-gallery<TAG>:"},
					Stdout: nil,
					Stderr: nil,
				},
				*&exec.Cmd{
					Path:   "docker",
					Args:   []string{"docker", "push", "quay.io/modcloth/style-gallery"},
					Stdout: nil,
					Stderr: nil,
				},
			},
		}
		expectedInstructionSet = &InstructionSet{
			DockerBuildOpts: []string{"--rm", "--no-cache"},
			DockerTagOpts:   []string{},
			Containers: map[string]builderfile.ContainerSection{
				"base": *&builderfile.ContainerSection{
					Dockerfile: "Dockerfile.base",
					Included:   []string{},
					Excluded: []string{
						"spec",
						"tmp",
					},
					Registry: "quay.io/modcloth",
					Project:  "style-gallery",
					Tags: []string{
						"base",
					},
				},
				"app": *&builderfile.ContainerSection{
					Dockerfile: "Dockerfile",
					Included:   []string{},
					Excluded: []string{
						"spec",
						"tmp",
					},
					Registry: "quay.io/modcloth",
					Project:  "style-gallery",
					Tags: []string{
						"git describe --always",
						"git rev-parse -q --abbrev-ref HEAD",
						"git rev-parse -q HEAD",
					},
				},
			},
		}
		expectedBuilderfile = &builderfile.Builderfile{
			Docker: *&builderfile.Docker{
				BuildOpts: []string{"--rm", "--no-cache"},
				TagOpts:   []string{},
			},
			Containers: map[string]builderfile.ContainerSection{
				"global": *&builderfile.ContainerSection{
					Dockerfile: "",
					Included:   []string{},
					Excluded:   []string{"spec", "tmp"},
					Registry:   "quay.io/modcloth",
					Project:    "style-gallery",
					Tags: []string{
						"git describe --always",
						"git rev-parse -q --abbrev-ref HEAD",
						"git rev-parse -q HEAD",
					},
				},
				"base": *&builderfile.ContainerSection{
					Dockerfile: "Dockerfile.base",
					Included:   []string{},
					Excluded:   []string{},
					Registry:   "",
					Project:    "",
					Tags:       []string{"base"},
				},
				"app": *&builderfile.ContainerSection{
					Dockerfile: "Dockerfile",
					Included:   []string{},
					Excluded:   []string{},
					Registry:   "",
					Project:    "",
					Tags:       nil,
				},
			},
		}
	)

	BeforeEach(func() {
		validFile = fmt.Sprintf("%s/spec/fixtures/Builderfile", top)
		invalidFile = fmt.Sprintf("%s/specs/fixtures/foodoesnotexist", top)
		subject = nil
	})

	Context("with a valid Builderfile", func() {

		It("produces an openable file", func() {
			subject, _ := NewParser(validFile, &log.NullLogger{})
			Expect(subject.IsOpenable()).To(Equal(true))
		})

		It("returns a non empty string and a nil error as raw data", func() {
			subject, _ := NewParser(validFile, &log.NullLogger{})
			raw, err := subject.getRaw()
			Expect(len(raw)).ToNot(Equal(0))
			Expect(err).To(BeNil())
		})

		It("returns a fully parsed Builderfile", func() {
			subject, _ := NewParser(validFile, &log.NullLogger{})
			actual, err := subject.rawToStruct()
			Expect(expectedBuilderfile).To(Equal(actual))
			Expect(err).To(BeNil())
		})

		It("further processes the Builderfile into an InstructionSet", func() {
			subject, _ := NewParser(validFile, &log.NullLogger{})
			actual, err := subject.structToInstructionSet()
			Expect(expectedInstructionSet).To(Equal(actual))
			Expect(err).To(BeNil())
		})

		XIt("further processes the InstructionSet into an CommandSequence", func() {
			subject, _ := NewParser(validFile, &log.NullLogger{})
			actual, err := subject.instructionSetToCommandSequence()
			Expect(expectedCommandSequence).To(Equal(actual))
			Expect(err).To(BeNil())
		})
	})

	Context("with an invalid Builderfile", func() {
		It("returns an empty string and error as raw data", func() {
			subject, _ := NewParser(invalidFile, &log.NullLogger{})
			raw, err := subject.getRaw()
			Expect(raw).To(Equal(""))
			Expect(err).ToNot(BeNil())
		})
	})
})
