SHELL := /bin/bash
DOCKER ?= docker
B := github.com/rafecolton/builder
TEST_LIBS := $(B)/spec
TARGETS := \
  $(B)/builder \
  $(B)/linter \
  $(B)/log \
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

export GOPATH
export GOBIN

help:
	@echo "Usage: TODO"

all: clean build test

clean:
	go clean -x -i $(TARGETS)
	rm -rf $${GOPATH%%:*}/src/github.com/rafecolton/builder
	rm -rf $${GOPATH%%:*}/bin/builder

quick: build
	@echo "----------"
	@builder --version
	@echo "----------"
	@builder --help
	@echo "----------"

build: linkthis deps
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(TARGETS)

linkthis:
	which gvm >/dev/null && (test -d $${GOPATH%%:*}/src/github.com/rafecolton/builder || gvm linkthis github.com/rafecolton/builder)

deps: godep
	$(GOBIN)/godep restore
	go get -x github.com/onsi/ginkgo/ginkgo
	go get -x github.com/onsi/gomega

godep:
	go get -x github.com/tools/godep

test: build fmtpolice
	$(GOBIN)/ginkgo -nodes=10 -noisyPendings -r -race -v .

fmtpolice:
	set -e ; for f in $(shell git ls-files '*.go'); do gofmt $$f | diff -u $$f - ; done

container:
	#TODO: docker build

.PHONY: godep test
default: help
