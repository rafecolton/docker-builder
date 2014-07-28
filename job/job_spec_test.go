package job_test

import (
	. "github.com/modcloth/docker-builder/job"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewSpec()", func() {
	var (
		args []byte
	)

	Context("when required args are missing", func() {
		It("returns an error when account is not provided", func() {
			args = []byte(`{
			  "repo": "kamino-test",
			  "ref": "master"
			}`)
			spec, _ := NewSpec(args)
			err := spec.Validate()

			Expect(err).ToNot(BeNil())
		})

		It("returns an error when repo is not provided", func() {
			args = []byte(`{
			  "account": "modcloth-labs",
			  "ref": "master"
			}`)
			spec, _ := NewSpec(args)
			err := spec.Validate()

			Expect(err).ToNot(BeNil())
		})

		It("returns an error when no ref is provided", func() {
			args = []byte(`{
			  "account": "modcloth-labs",
			  "repo": "kamino-test"
			}`)
			spec, _ := NewSpec(args)
			err := spec.Validate()

			Expect(err).ToNot(BeNil())
		})

		It("returns an error when provided args are not valid json", func() {
			args = []byte(`foo`)
			spec, err := NewSpec(args)

			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Context("when required args are present", func() {
		It("returns a valid job spec", func() {
			args = []byte(`{
			  "account": "modcloth-labs",
			  "repo": "kamino-test",
			  "ref": "master"
			}`)
			spec, err := NewSpec(args)

			Expect(spec).ToNot(BeNil())
			Expect(err).To(BeNil())
		})
	})
})
