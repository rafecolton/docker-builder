package bob_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rafecolton/bob"
	"testing"
)

import (
//"github.com/rafecolton/bob/builderfile"
//"github.com/rafecolton/bob/parser"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Builder Specs")
}

var _ = Describe("Build", func() {
	var ()

	BeforeEach(func() {
	})

	Context("when", func() {
		XIt("", func() {
		})
	})
})

var _ = Describe("CommandSequence", func() {
	var (
		subject *Builder
		//invalidInstructions  *builderfile.Builderfile
		//expectedInstructions = &parser.InstructionSet{
		//DockerBuildOpts: []string{"--rm", "--no-cache"},
		//DockerTagOpts:   []string{},
		//Containers: map[string]builderfile.ContainerSection{
		//"base": *&builderfile.ContainerSection{
		//Dockerfile: "Dockerfile.base",
		//Included:   []string{},
		//Excluded: []string{
		//"spec",
		//"tmp",
		//},
		//Registry: "quay.io/modcloth",
		//Project:  "style-gallery",
		//Tags: []string{
		//"base",
		//},
		//},
		//"app": *&builderfile.ContainerSection{
		//Dockerfile: "Dockerfile",
		//Included:   []string{},
		//Excluded: []string{
		//"spec",
		//"tmp",
		//},
		//Registry: "quay.io/modcloth",
		//Project:  "style-gallery",
		//Tags: []string{
		//"git describe --always",
		//"git rev-parse -q --abbrev-ref HEAD",
		//"git rev-parse -q HEAD",
		//},
		//},
		//},
		//}
	)

	BeforeEach(func() {
		subject = NewBuilder()
	})

	Context("when determining commands from a valid instruction set", func() {
		XIt("produces the correct command sequence", func() {
			//sequence, _ := subject.CommandSequence(validInstructions)
			//Expect(sequence).To(Equal([]string{}))
		})

		XIt("does not produce an error", func() {
			//_, err := subject.CommandSequence(validInstructions)
			//Expect(err).To(BeNil())
		})
	})

	Context("when determining commands from an invalid instruction set", func() {
		XIt("produces an empty command sequence", func() {
			//sequence, _ := subject.CommandSequence(invalidInstructions)
			//Expect(sequence).To(Equal([]string{}))
		})

		XIt("produces an error", func() {
			//_, err := subject.CommandSequence(invalidInstructions)
			//Expect(err).ToNot(BeNil())
		})
	})
})
