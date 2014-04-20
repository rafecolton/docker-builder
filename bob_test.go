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
	var (
		subject *Builder
	)

	BeforeEach(func() {
		subject = nil
	})

	Context("when", func() {
		It("", func() {
		})
	})
})
