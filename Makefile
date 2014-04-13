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

#SRC_DIR := src/github.com/rafecolton

#BUILD_DIR := $(SRC_DIR)/builder

#GO_TAG_ARGS ?= -tags full

GOPATH := "${PWD}/$(godep path)"

help:
	@echo "Usage: TODO"

all: clean build test

clean:
	go clean -x -i $(TARGETS)
	rm -rf $${GOPATH%%:*}/src/github.com/rafecolton/builder
	rm -rf $${GOPATH%%:*}/bin/builder

linkthis:
	which gvm >/dev/null && (test -d $${GOPATH%%:*}/src/github.com/rafecolton/builder || gvm linkthis github.com/rafecolton/builder)
	#if [ -n ${CI} ] ; then gvm linkthis github.com/rafecolton/builder ; fi

quick: build
	@echo "----------"
	@builder --version
	@echo "----------"
	@builder --help
	@echo "----------"

build: linkthis deps
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(TARGETS)

godep:
	go get -x github.com/tools/godep

deps: godep
	godep restore
	go get -x github.com/onsi/ginkgo/ginkgo
	go get -x github.com/onsi/gomega

savedeps:
	godep save -copy=false $(TEST_LIBRARIES) $(TARGETS)

test: deps fmtpolice
	$${GOPATH%%:*}/bin/ginkgo -nodes=10 -noisyPendings -r -race .

fmtpolice:
	set -e ; for f in $(shell git ls-files '*.go'); do gofmt $$f | diff -u $$f - ; done
	#fail=0 ; for f in $(shell git ls-files '*.go'); do v="$$(golint $$f)" ; if [ ! -z "$$v" ] ; then echo "$$v" ; fail=1 ; fi ; done ; [ $$fail = 0 ]

container:
	#TODO: docker build

.PHONY: godep test
default: help
