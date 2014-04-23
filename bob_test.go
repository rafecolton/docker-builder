package bob

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//. "github.com/rafecolton/bob"
	"testing"
)

import (
	//"github.com/rafecolton/bob/builderfile"
	"github.com/rafecolton/bob/parser"
)

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Builder Specs")
}

var _ = Describe("Setup", func() {
	var (
		branch          string
		rev             string
		short           string
		top             string
		expectedFiles   []string
		subject         *Builder
		baseSubSequence = &parser.SubSequence{
			Metadata: &parser.SubSequenceMetadata{
				Name:       "base",
				Dockerfile: "Dockerfile.base",
				Excluded:   []string{"spec", "tmp"},
				Included:   []string{"Gemfile", "Gemfile.lock"},
			},
			SubCommand: []exec.Cmd{
				*&exec.Cmd{
					Path: "docker",
					Args: []string{"docker", "build", "-t", "quay.io/modcloth/style-gallery:035c4ea0-d73b-5bde-7d6f-c806b04f2ec3", "--rm", "--no-cache", "."},
				},
				*&exec.Cmd{
					Path: "docker",
					Args: []string{"docker", "tag", "<IMG>", "quay.io/modcloth/style-gallery:base"},
				},
				*&exec.Cmd{
					Path: "docker",
					Args: []string{"docker", "push", "quay.io/modcloth/style-gallery"},
				},
			},
		}
		appSubSequence = &parser.SubSequence{
			Metadata: &parser.SubSequenceMetadata{
				Name:       "app",
				Dockerfile: "Dockerfile",
				Excluded:   []string{"spec", "tmp"},
				Included:   []string{},
			},
			SubCommand: []exec.Cmd{
				*&exec.Cmd{
					Path: "docker",
					Args: []string{"docker", "build", "-t", "quay.io/modcloth/style-gallery:035c4ea0-d73b-5bde-7d6f-c806b04f2ec3", "--rm", "--no-cache", "."},
				},
				*&exec.Cmd{
					Path: "docker",
					Args: []string{"docker", "tag", "<IMG>", fmt.Sprintf("quay.io/modcloth/style-gallery:%s", branch)},
				},
				*&exec.Cmd{
					Path: "docker",
					Args: []string{"docker", "tag", "<IMG>", fmt.Sprintf("quay.io/modcloth/style-gallery:%s", rev)},
				},
				*&exec.Cmd{
					Path: "docker",
					Args: []string{"docker", "tag", "<IMG>", fmt.Sprintf("quay.io/modcloth/style-gallery:%s", short)},
				},
				*&exec.Cmd{
					Path: "docker",
					Args: []string{"docker", "push", "quay.io/modcloth/style-gallery"},
				},
			},
		}
	)

	BeforeEach(func() {
		subject = NewBuilder(nil, false)
		top = os.ExpandEnv("${PWD}")
		git, _ := exec.LookPath("git")
		// branch
		branchCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "rev-parse", "-q", "--abbrev-ref", "HEAD"},
		}

		branchBytes, _ := branchCmd.Output()
		branch = string(branchBytes)[:len(branchBytes)-1]

		// rev
		revCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "rev-parse", "-q", "HEAD"},
		}
		revBytes, _ := revCmd.Output()
		rev = string(revBytes)[:len(revBytes)-1]

		// short
		shortCmd := &exec.Cmd{
			Path: git,
			Dir:  top,
			Args: []string{git, "describe", "--always"},
		}
		shortBytes, _ := shortCmd.Output()
		short = string(shortBytes)[:len(shortBytes)-1]
	})

	AfterEach(func() {
		subject.CleanWorkdir()
	})

	Context("with the base container sequence", func() {
		It("places the correct files in the workdir", func() {
			subject.SetNextSubSequence(baseSubSequence)
			subject.CleanWorkdir()
			subject.Setup()

			expectedFiles = []string{
				"Dockerfile.base",
				"Gemfile",
				"Gemfile.lock",
				"README.txt",
			}

			files, _ := ioutil.ReadDir(subject.Workdir())
			fileNames := []string{}
			for _, v := range files {
				fileNames = append(fileNames, v.Name())
			}

			sort.Strings(fileNames)
			sort.Strings(expectedFiles)

			Expect(fileNames).To(Equal(expectedFiles))
		})
	})

	Context("with the app container sequence", func() {
		It("places the correct files in the workdir", func() {
			subject.SetNextSubSequence(appSubSequence)
			subject.CleanWorkdir()
			subject.Setup()

			expectedFiles = []string{
				"Dockerfile",
				"Dockerfile.base",
				"Gemfile",
				"Gemfile.lock",
				"foo",
				"README.txt",
				"other_file.txt",
			}

			files, _ := ioutil.ReadDir(subject.Workdir())
			fileNames := []string{}
			for _, v := range files {
				fileNames = append(fileNames, v.Name())
			}

			sort.Strings(fileNames)
			sort.Strings(expectedFiles)

			Expect(fileNames).To(Equal(expectedFiles))
		})
	})
})
