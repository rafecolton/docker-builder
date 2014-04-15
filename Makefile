SHELL := /bin/bash
SUDO ?= sudo
DOCKER ?= docker
B := github.com/rafecolton/bob
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

BATS_INSTALL_DIR := /usr/local

GOPATH := $(PWD)/Godeps/_workspace
GOBIN := $(GOPATH)/bin
PATH := $(GOPATH):$(PATH)

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
	rm -rf $${GOPATH%%:*}/src/github.com/rafecolton/bob
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

binclean:
	rm -f $${GOPATH%%:*}/bin/builder

build: linkthis deps binclean
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(TARGETS)

linkthis:
	if which gvm >/dev/null && \
	  [[ ! -d $${GOPATH%%:*}/src/github.com/rafecolton/bob ]] ; then \
	  gvm linkthis github.com/rafecolton/bob ; \
	  fi

godep:
	go get -x github.com/tools/godep

deps: godep
	$(GOBIN)/godep restore
	go get -x github.com/golang/lint/golint
	go get -x github.com/onsi/ginkgo/ginkgo
	go get -x github.com/onsi/gomega
	if ! which bats >/dev/null ; then \
	  git clone https://github.com/sstephenson/bats.git && \
	  (cd bats && $(SUDO) ./install.sh ${BATS_INSTALL_DIR}) && \
	  rm -rf bats ; \
	  fi

test: build fmtpolice ginkgo bats

fmtpolice: deps fmt lint

fmt: line
	@echo "checking fmt"
	@set -e ; \
	  for f in $(shell git ls-files '*.go'); do \
	  gofmt $$f | diff -u $$f - ; \
	  done

lint: line
	@echo "checking lint"
	@for file in $(shell git ls-files '*.go') ; do \
	  if [[ "$$($(GOBIN)/golint $$file)" =~ ^[[:blank:]]*$$ ]] ; then \
	  echo yayyy >/dev/null ; \
	  else exit 1 ; fi \
	  done

ginkgo: line
	$(GOBIN)/ginkgo -nodes=10 -noisyPendings -r -race -v .

bats: line
	$(BATS_INSTALL_DIR)/bin/bats $(shell git ls-files '*.bats')

line:
	@echo "----------"

container:
	#TODO: docker build

.PHONY: container line bats ginkgo
.PHONY:	lint fmt fmtpolice test deps
.PHONY:	linkthis build quick clean all help
default: help
