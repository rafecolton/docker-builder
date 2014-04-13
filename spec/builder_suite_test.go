package builder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//"github.com/rafecolton/builder"
	"testing"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Builder Suite")
}
