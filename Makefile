SHELL := /bin/bash
SUDO ?= sudo
DOCKER ?= docker
B := github.com/rafecolton/bob
TARGETS := \
  $(B)/builder \
  $(B)/builderfile \
  $(B)/dclient \
  $(B)/log \
  $(B)/parser \
  $(B)/parser/uuid \
  $(B)/version
REV_VAR := $(B)/version.RevString
VERSION_VAR := $(B)/version.VersionString
BRANCH_VAR := $(B)/version.BranchString
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
REPO_BRANCH := $(shell git rev-parse -q --abbrev-ref HEAD)
GOBUILD_VERSION_ARGS := -ldflags "\
  -X $(REV_VAR) $(REPO_REV) \
  -X $(VERSION_VAR) $(REPO_VERSION) \
  -X $(BRANCH_VAR) $(REPO_BRANCH)"

BATS_INSTALL_DIR ?= /usr/local
GINKGO_PATH ?= "."
GOPATH := $(PWD)/Godeps/_workspace
GOBIN := $(GOPATH)/bin
PATH := $(GOPATH):$(PATH)

export BATS_INSTALL_DIR
export GINKGO_PATH
export GOPATH
export GOBIN
export PATH

.PHONY: default help
default: help
help:
	@echo "Usage: make [target]"
	@echo
	@echo "Some Useful Options:"
	@echo
	@echo "  help: display this message"
	@echo
	@echo "  all: clean build test"
	@echo
	@echo "  quick: build + invokes builder a couple times (good for debugging)"
	@echo
	@echo "  build: gvm linkthis plus installing libs plus installing deps"
	@echo
	@echo "  test: build fmtpolice, ginkgo tests, and bats tests"
	@echo
	@echo "  dev: set up the dev toolchain (deps + gox)"

.PHONY: all
all: clean build test

.PHONY: clean
clean:
	go clean -i -r $(TARGETS) || true
	rm -rf $${GOPATH%%:*}/src/github.com/rafecolton/bob
	rm -f $${GOPATH%%:*}/bin/builder
	rm -rf Godeps/_workspace/*

.PHONY: quick
quick: build
	@echo "----------"
	@builder --version
	@echo "----------"
	@builder --help
	@echo "----------"
	@builder
	@echo "----------"

.PHONY: binclean
binclean:
	rm -f $${GOPATH%%:*}/bin/builder
	rm -f ./builds/builder-dev
	rm -f ./builds/darwin_amd64
	rm -f ./builds/linux_amd64

.PHONY: build
build: linkthis deps binclean
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(TARGETS)


.PHONY: gox-build
gox-build: linkthis deps binclean
	gox -osarch="darwin/amd64" -output "builds/builder-dev" $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(TARGETS)
	gox -output="builds/{{.OS}}_{{.Arch}}" -arch="amd64" -os="darwin linux" $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(TARGETS)


.PHONY: linkthis
linkthis:
	@echo "gvm linkthis'ing this..."
	@if which gvm >/dev/null && \
	  [[ ! -d $${GOPATH%%:*}/src/github.com/rafecolton/bob ]] ; then \
	  gvm linkthis github.com/rafecolton/bob ; \
	  fi

.PHONY: godep
godep:
	go get github.com/tools/godep


.PHONY: deps
deps: godep
	@echo "godep restoring..."
	$(GOBIN)/godep restore
	go get github.com/golang/lint/golint
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	@echo "installing bats..."
	@if ! which bats >/dev/null ; then \
	  git clone https://github.com/sstephenson/bats.git && \
	  (cd bats && $(SUDO) ./install.sh ${BATS_INSTALL_DIR}) && \
	  rm -rf bats ; \
	  fi


.PHONY: test
test: build fmtpolice ginkgo bats


.PHONY: fmtpolice
fmtpolice: deps fmt lint


.PHONY: fmt
fmt:
	@echo "----------"
	@echo "checking fmt"
	@set -e ; \
	  for f in $(shell git ls-files '*.go'); do \
	  gofmt $$f | diff -u $$f - ; \
	  done


.PHONY: linter
linter:
	go get github.com/golang/lint/golint


.PHONY: lint
lint: linter
	@echo "----------"
	@echo "checking lint"
	@for file in $(shell git ls-files '*.go') ; do \
	  if [[ "$$($(GOBIN)/golint $$file)" =~ ^[[:blank:]]*$$ ]] ; then \
	  echo yayyy >/dev/null ; \
	  else $(MAKE) lintv && exit 1 ; fi \
	  done


.PHONY: lintv
lintv:
	@echo "----------"
	@for file in $(shell git ls-files '*.go') ; do $(GOBIN)/golint $$file ; done


.PHONY: ginkgo
ginkgo:
	@echo "----------"
	@if [[ "$(GINKGO_PATH)" == "." ]] ; then \
	  echo "$(GOBIN)/ginkgo -nodes=10 -noisyPendings -race -r ." && \
	  $(GOBIN)/ginkgo -nodes=10 -noisyPendings -race -r . ; \
	  else echo "$(GOBIN)/ginkgo -nodes=10 -noisyPendings -race --v $(GINKGO_PATH)" && \
	  $(GOBIN)/ginkgo -nodes=10 -noisyPendings -race --v $(GINKGO_PATH) ; \
	  fi


.PHONY: bats
bats:
	@echo "----------"
	$(BATS_INSTALL_DIR)/bin/bats --pretty $(shell git ls-files '*.bats')


.PHONY: gox
gox:
	@if which gox ; then \
	  echo "not installing gox, gox already installed." ; \
	  else \
	  go get github.com/mitchellh/gox ; \
	  gox -build-toolchain ; \
	  fi \


.PHONY: dev
dev: deps gox


.PHONY: container
container:
	#TODO: docker build
