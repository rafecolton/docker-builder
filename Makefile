SHELL := /bin/bash
SUDO ?= sudo
DOCKER ?= docker
B := github.com/rafecolton/builder
#TEST_LIBS := $(B)/spec
TARGETS := \
  $(B)/builder \
  $(B)/builderfile \
  $(B)/log \
  $(B)/parser \
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

GOPATH := $(PWD)/Godeps/_workspace
GOBIN := $(GOPATH)/bin
BATS_INSTALL_DIR := /usr/local

export GOPATH
export GOBIN
export BATS_INSTALL_DIR

help:
	@echo "Usage: make [target]"
	@echo
	@echo "Options:"
	@echo
	@echo "  help/default: display this message"
	@echo
	@echo "  all: clean build test"
	@echo
	@echo "  quick: build + invokes builder a couple times (good for debugging)"
	@echo
	@echo "  build: gvm linkthis plus installing libs plus installing deps"
	@echo
	@echo "  test: build fmtpolice and ginkgotests"
	@echo

all: clean build test

clean:
	go clean -x -i $(TARGETS)
	rm -rf $${GOPATH%%:*}/src/github.com/rafecolton/builder
	rm -f $${GOPATH%%:*}/bin/builder
	rm -rf Godeps/_workspace/*

quick: build
	@echo "----------"
	@builder --version
	@echo "----------"
	@builder --help
	@echo "----------"
	@builder
	@echo "----------"

build: linkthis deps
	rm -f $${GOPATH%%:*}/bin/builder
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(TARGETS)

linkthis:
	which gvm >/dev/null && (test -d $${GOPATH%%:*}/src/github.com/rafecolton/builder || gvm linkthis github.com/rafecolton/builder)

deps: godep
	$(GOBIN)/godep restore
	go get -x github.com/golang/lint/golint
	go get -x github.com/onsi/ginkgo/ginkgo
	go get -x github.com/onsi/gomega
	if ! which bats >/dev/null ; then git clone https://github.com/sstephenson/bats.git && (cd bats && $(SUDO) ./install.sh ${BATS_INSTALL_DIR}) && rm -rf bats ; fi

godep:
	go get -x github.com/tools/godep

test: build fmtpolice ginkgo bats

fmtpolice: deps
	@echo "----------"
	set -e ; for f in $(shell git ls-files '*.go'); do gofmt $$f | diff -u $$f - ; done
	@echo "----------"
	fail=0 ; for f in $(shell git ls-files '*.go'); do v="$$($(GOBIN)/golint $$f)" ; if [ ! -z "$$v" ] ; then echo "$$v" ; fail=1 ; fi ; done ; [ $$fail -eq 0 ]

ginkgo:
	@echo "----------"
	$(GOBIN)/ginkgo -nodes=10 -noisyPendings -r -race -v .

bats:
	@echo "----------"
	$(BATS_INSTALL_DIR)/bin/bats $(shell git ls-files '*.bats')

container:
	#TODO: docker build

.PHONY: godep test
default: help
