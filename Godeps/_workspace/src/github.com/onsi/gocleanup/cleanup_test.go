package gocleanup_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gocleanup"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleanup", func() {
	Context("with no cleanup functions registered", func() {
		It("should do nothing when told to clean up", func() {
			Cleanup() //ha!
		})
	})

	Context("with cleanup functions registered", func() {
		var calls []string
		BeforeEach(func() {
			calls = []string{}
			Register(func() {
				calls = append(calls, "A")
			})
			Register(func() {
				calls = append(calls, "B")
			})
		})

		It("should call the functions when told to clean up", func() {
			Cleanup()
			Ω(calls).Should(Equal([]string{"A", "B"}))
		})

		It("should unregister the functions upon cleanup", func() {
			Cleanup()
			Ω(calls).Should(Equal([]string{"A", "B"}))

			calls = []string{}
			Cleanup()
			Ω(calls).Should(BeEmpty())
		})
	})
})
