package parser

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
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
			subject = New()
			subject.Builderfile = validFile
			Expect(subject.IsOpenable()).To(Equal(true))
		})

		It("returns a non empty string as raw data", func() {
			subject = New()
			subject.Builderfile = validFile
			raw, _ := subject.ParseRaw()
			Expect(len(raw)).ToNot(Equal(0))
		})

		It("returns a nil error", func() {
			subject = New()
			subject.Builderfile = validFile
			_, err := subject.ParseRaw()
			Expect(err).To(BeNil())
		})
	})

	Context("with an invalid Builderfile", func() {
		It("returns an error", func() {
			subject = New()
			subject.Builderfile = invalidFile
			Expect(subject.IsOpenable()).To(Equal(false))
		})

		It("returns an empty string as raw data", func() {
			subject = New()
			subject.Builderfile = invalidFile
			raw, _ := subject.ParseRaw()
			Expect(raw).To(Equal(""))
		})

		It("returns a non-nil error", func() {
			subject = New()
			subject.Builderfile = invalidFile
			_, err := subject.ParseRaw()
			Expect(err).ToNot(BeNil())
		})
	})
})
