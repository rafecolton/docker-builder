package job_test

import (
	. "github.com/modcloth/docker-builder/job"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
)

var _ = Describe("NewJobSpec()", func() {
	var (
		args interface{}
	)

	Context("when required args are missing", func() {
		It("returns an error when account is not provided", func() {
			args = makeArg(`{
			  "repo": "kamino-test",
			  "ref": "master"
			}`)
			spec, err := NewJobSpec(args)

			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("returns an error when repo is not provided", func() {
			args = makeArg(`{
			  "account": "modcloth-labs",
			  "ref": "master"
			}`)
			spec, err := NewJobSpec(args)

			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("returns an error when args are empty", func() {
			emptyArgs := []interface{}{}
			spec, err := NewJobSpec(emptyArgs...)

			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("returns an error when no ref is provided", func() {
			args = makeArg(`{
			  "account": "modcloth-labs",
			  "repo": "kamino-test"
			}`)
			spec, err := NewJobSpec(args)

			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("returns an error when provided args are not valid json", func() {
			args = makeArg(`foo`)
			spec, err := NewJobSpec(args)

			Expect(spec).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Context("when required args are present", func() {
		It("returns a valid job spec", func() {
			args = makeArg(`{
			  "account": "modcloth-labs",
			  "repo": "kamino-test",
			  "ref": "master"
			}`)
			spec, err := NewJobSpec(args)

			Expect(spec).ToNot(BeNil())
			Expect(err).To(BeNil())
		})
	})
})

func makeArg(rawJSON string) interface{} {

	type jsonType map[string]interface{}

	ret := &jsonType{}

	_ = json.Unmarshal([]byte(rawJSON), ret)

	return interface{}(*ret)
}
