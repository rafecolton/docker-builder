package uuid_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rafecolton/docker-builder/parser/uuid"
	"testing"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UUID Specs")
}

var _ = Describe("GenerateUUID", func() {
	var (
		seeded Generator
		random Generator
	)

	Context("with a random uuid generator", func() {
		It("generates a different uuid every time", func() {
			random = NewUUIDGenerator(true)
			alpha, _ := random.NextUUID()
			beta, _ := random.NextUUID()
			Expect(alpha).ToNot(Equal(beta))
		})
	})

	Context("with a seeded uuid generator", func() {
		It("generates the same uuid every time", func() {
			seeded = NewUUIDGenerator(false)
			alpha, _ := seeded.NextUUID()
			beta, _ := seeded.NextUUID()
			Expect(alpha).To(Equal(beta))
		})
	})
})
