package builder

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Version Specs")
}

var _ = Describe("Version", func() {
	var (
	//version = "0.1.0"
	)
})
