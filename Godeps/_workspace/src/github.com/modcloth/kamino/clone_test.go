package kamino_test

import (
	. "github.com/modcloth/kamino"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/modcloth/go-fileutils"
)

var (
	requestedSHA    = "df66a4216affe8fe29af354f78e9016781e7bb8e"
	nonRequestedSHA = "9830dc808697ba1c7db91df908cb99eb9b3062fe"
	preCloneGenome  *Genome
	cloneGenome     *Genome
	cacheDirSuffix  = "modcloth-labs/kamino-test"
	subject         *CloneFactory
	tmpdir          string
	err             error
	path            string
	cachePath       string
)

var _ = Describe("no cache", func() {
	BeforeEach(func() {
		genome := &Genome{
			Depth:    "50",
			Account:  "modcloth-labs",
			Repo:     "kamino-test",
			UseCache: "no",
			Ref:      requestedSHA,
		}

		tmpdir, _ = ioutil.TempDir("", "kamino-test")
		subject, _ = NewCloneFactory(tmpdir)
		cachePath = fmt.Sprintf("%s/%s", tmpdir, cacheDirSuffix)

		path, _ = subject.Clone(genome)
	})

	AfterEach(func() {
		fileutils.RmRF(tmpdir)
	})

	It("does not create a cache directory", func() {
		Expect(path).ToNot(Equal(cachePath))
	})

	It("clones to the correct ref", func() {
		ref, _ := GetRef(path)

		Expect(ref).To(Equal(requestedSHA))
	})
})

var _ = Describe("create cache", func() {
	BeforeEach(func() {
		cloneGenome = &Genome{
			Depth:    "50",
			Account:  "modcloth-labs",
			Repo:     "kamino-test",
			UseCache: "create",
			Ref:      requestedSHA,
		}

		tmpdir, _ = ioutil.TempDir("", "kamino-test")
		subject, _ = NewCloneFactory(tmpdir)
		cachePath = fmt.Sprintf("%s/%s", tmpdir, cacheDirSuffix)

		path, _ = subject.Clone(cloneGenome)
	})

	AfterEach(func() {
		fileutils.RmRF(tmpdir)
	})

	It("creates the cache directory", func() {
		Expect(path).To(Equal(cachePath))
	})

	It("correctly clones to the correct ref", func() {
		ref, _ := GetRef(path)

		Expect(ref).To(Equal(requestedSHA))
	})
})

var _ = Describe("force cache", func() {
	BeforeEach(func() {
		preCloneGenome = &Genome{
			Depth:    "50",
			Account:  "modcloth-labs",
			Repo:     "kamino-test",
			UseCache: "create",
			Ref:      nonRequestedSHA,
		}
		cloneGenome = &Genome{
			Depth:    "50",
			Account:  "modcloth-labs",
			Repo:     "kamino-test",
			UseCache: "force",
			Ref:      requestedSHA,
		}
		tmpdir, _ = ioutil.TempDir("", "kamino-test")
		subject, _ = NewCloneFactory(tmpdir)
	})

	AfterEach(func() {
		fileutils.RmRF(tmpdir)
	})

	Context("cache directory exists prior to cloning", func() {
		It("successfully clones the to the correct ref", func() {

			//////////////////////////////////////////
			// make sure the directory exists first //
			//////////////////////////////////////////
			subject.Clone(preCloneGenome)
			/////////////////////////////////////////

			cachePath = fmt.Sprintf("%s/%s", tmpdir, cacheDirSuffix)

			subject.Clone(cloneGenome)

			ref, _ := GetRef(cachePath)

			Expect(ref).To(Equal(requestedSHA))
		})
	})

	Context("cache directory does not exist prior to cloning", func() {
		It("returns an error", func() {
			path, err = subject.Clone(cloneGenome)

			Expect(err).ToNot(BeNil())
		})
	})
})

var _ = Describe("use cache if available", func() {
	BeforeEach(func() {
		preCloneGenome = &Genome{
			Depth:    "50",
			Account:  "modcloth-labs",
			Repo:     "kamino-test",
			UseCache: "create",
			Ref:      nonRequestedSHA,
		}
		cloneGenome = &Genome{
			Depth:    "50",
			Account:  "modcloth-labs",
			Repo:     "kamino-test",
			UseCache: "if_available",
			Ref:      requestedSHA,
		}

		tmpdir, _ = ioutil.TempDir("", "kamino-test")
		subject, _ = NewCloneFactory(tmpdir)
		cachePath = fmt.Sprintf("%s/%s", tmpdir, cacheDirSuffix)
	})

	AfterEach(func() {
		fileutils.RmRF(tmpdir)
	})

	Context("cache directory exists prior to cloning", func() {
		It("successfully clones the to the correct ref", func() {

			//////////////////////////////////////////
			// make sure the directory exists first //
			//////////////////////////////////////////
			subject.Clone(preCloneGenome)
			/////////////////////////////////////////

			path, _ := subject.Clone(cloneGenome)
			ref, _ := GetRef(cachePath)

			Expect(ref).To(Equal(requestedSHA))
			Expect(path).To(Equal(cachePath))
		})
	})

	Context("cache directory does not exist prior to cloning", func() {
		It("does not create a cache directory", func() {
			path, _ = subject.Clone(cloneGenome)
			Expect(path).ToNot(Equal(cachePath))
		})
	})
})

func GetRef(path string) (string, error) {
	git, _ := fileutils.Which("git")

	cmd := &exec.Cmd{
		Path: git,
		Dir:  path,
		Args: []string{"git", "rev-parse", "-q", "HEAD"},
	}

	refBytes, err := cmd.Output()

	if err != nil || len(refBytes) == 0 {
		return "", err
	}

	return string(refBytes)[:len(refBytes)-1], nil
}
