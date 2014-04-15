package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rafecolton/bob/parser"
	"testing"
)

import (
	"fmt"
	"github.com/rafecolton/bob/builderfile"
	"os"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Specs")
}

var _ = Describe("Parse", func() {

	var (
		subject     *Parser
		validFile   string
		invalidFile string
		top         = os.ExpandEnv("${PWD}")
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

		It("returns a non empty string as raw data", func() {
			subject = NewParser(validFile, nil)
			raw, _ := subject.ParseRaw()
			Expect(len(raw)).ToNot(Equal(0))
		})

		It("returns a nil error", func() {
			subject = NewParser(validFile, nil)
			_, err := subject.ParseRaw()
			Expect(err).To(BeNil())
		})

		It("returns a fully parsed Builderfile", func() {
			subject = NewParser(validFile, nil)
			actual, _ := subject.Parse(false)

			expected := &builderfile.Builderfile{
				Docker: *&builderfile.Docker{
					BuildOpts: "--rm --no-cache",
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

			Expect(expected).To(Equal(actual))
		})
	})

	Context("with an invalid Builderfile", func() {
		It("returns an error", func() {
			subject = NewParser(invalidFile, nil)
			Expect(subject.IsOpenable()).To(Equal(false))
		})

		It("returns an empty string as raw data", func() {
			subject = NewParser(invalidFile, nil)
			raw, _ := subject.ParseRaw()
			Expect(raw).To(Equal(""))
		})

		It("returns a non-nil error", func() {
			subject = NewParser(invalidFile, nil)
			_, err := subject.ParseRaw()
			Expect(err).ToNot(BeNil())
		})
	})
})
