package parser

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

import (
	"fmt"
	"github.com/rafecolton/bob/builderfile"
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
			commands: map[string]exec.Cmd{
				"baseBuild": *&exec.Cmd{
					Path: "docker",
					Args: []string{
						"docker",
						"build",
						"-t",
						"quay.io/modcloth/style-gallery:latest",
						"--rm",
						"--no-cache",
					},
					Stdout: nil,
					Stderr: nil,
				},
				"baseTag0": *&exec.Cmd{
					Path: "docker",
					Args: []string{
						"docker",
						"tag",
						"<IMAGE>",
						"quay.io/modcloth/style-gallery:base",
					},
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
	})

	Context("with a valid Builderfile", func() {

		It("is an openable file", func() {
			subject = NewParser(validFile, nil)
			Expect(subject.IsOpenable()).To(Equal(true))
		})

		It("returns a non empty string and a nil error as raw data", func() {
			subject = NewParser(validFile, nil)
			raw, err := subject.getRaw()
			Expect(len(raw)).ToNot(Equal(0))
			Expect(err).To(BeNil())
		})

		It("returns a fully parsed Builderfile", func() {
			subject = NewParser(validFile, nil)
			actual, err := subject.rawToStruct()
			Expect(expectedBuilderfile).To(Equal(actual))
			Expect(err).To(BeNil())
		})

		It("further processes the Builderfile into an InstructionSet", func() {
			subject = NewParser(validFile, nil)
			actual, err := subject.structToInstructionSet()
			Expect(expectedInstructionSet).To(Equal(actual))
			Expect(err).To(BeNil())
		})

		It("further processes the InstructionSet into an CommandSequence", func() {
			subject = NewParser(validFile, nil)
			actual, err := subject.instructionSetToCommandSequence()
			Expect(expectedCommandSequence).To(Equal(actual))
			Expect(err).To(BeNil())
		})
	})

	Context("with an invalid Builderfile", func() {
		It("returns an error", func() {
			subject = NewParser(invalidFile, nil)
			Expect(subject.IsOpenable()).To(Equal(false))
		})

		It("returns an empty string as raw data", func() {
			subject = NewParser(invalidFile, nil)
			raw, _ := subject.getRaw()
			Expect(raw).To(Equal(""))
		})

		It("returns a non-nil error", func() {
			subject = NewParser(invalidFile, nil)
			_, err := subject.getRaw()
			Expect(err).ToNot(BeNil())
		})
	})
})
