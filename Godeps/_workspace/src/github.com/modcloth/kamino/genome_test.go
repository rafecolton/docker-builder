package kamino_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/modcloth/kamino"
)

var _ = Describe("NewGenome()", func() {

	var genomeSubject *Genome

	BeforeEach(func() {
		genomeSubject = &Genome{
			Depth:    "50",
			APIToken: "abc123",
			Account:  "modcloth-labs",
			Repo:     "kamino-test",
			Ref:      "123",
		}
	})

	Context("with a non-integer depth", func() {
		It("returns an error", func() {
			genomeSubject.Depth = "foo"
			err := ValidateGenome(genomeSubject)

			Expect(err).ToNot(BeNil())
		})
	})

	Context("with no account specified", func() {
		It("returns an error", func() {
			genomeSubject.Account = ""
			err := ValidateGenome(genomeSubject)

			Expect(err).ToNot(BeNil())
		})
	})

	Context("with no repo specified", func() {
		It("returns an error", func() {
			genomeSubject.Repo = ""
			err := ValidateGenome(genomeSubject)

			Expect(err).ToNot(BeNil())
		})
	})

	Context("with an invalid cache option", func() {
		It("returns an error", func() {
			genomeSubject.UseCache = "foo"
			err := ValidateGenome(genomeSubject)

			Expect(err).ToNot(BeNil())
		})
	})

	Context("with no ref specified", func() {
		It("returns an error", func() {
			genomeSubject.Ref = ""
			err := ValidateGenome(genomeSubject)

			Expect(err).ToNot(BeNil())
		})
	})
})
