package kamino_test

import (
	"github.com/Sirupsen/logrus"
	"github.com/modcloth/kamino"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	kamino.Logger = logrus.New()
	kamino.Logger.Level = logrus.PanicLevel
	RunSpecs(t, "Kamino Spec")
}
