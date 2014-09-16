SHELL := /bin/bash
SUDO ?= sudo
DOCKER ?= docker
B := github.com/rafecolton/docker-builder
PACKAGES := ./...
REV_VAR := $(B)/version.RevString
VERSION_VAR := $(B)/version.VersionString
BRANCH_VAR := $(B)/version.BranchString
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
REPO_BRANCH := $(shell git rev-parse -q --abbrev-ref HEAD)
GOBUILD_VERSION_ARGS := -ldflags "\
  -X $(REV_VAR) $(REPO_REV) \
  -X $(VERSION_VAR) $(REPO_VERSION) \
  -X $(BRANCH_VAR) $(REPO_BRANCH) \
  -w \
"

BATS_INSTALL_DIR ?= /usr/local
GINKGO_PATH ?= "."

BATS_OUT_FORMAT=$(shell bash -c "echo $${CI+--tap}")
GOPATH := $(shell echo $${GOPATH%%:*})

# go build args
GO_TAG_ARGS := -tags netgo

export GOPATH

.PHONY: all
all: binclean clean build test

.PHONY: default help
default: help
help:
	@echo "Usage: make [target]"
	@echo
	@echo "Some Useful Options:"
	@echo
	@echo "  help: display this message"
	@echo
	@echo "  all: binclean clean build test"
	@echo
	@echo "  quick: build + invokes docker-builder a couple times (good for debugging)"
	@echo
	@echo "  build: installing libs plus installing deps"
	@echo
	@echo "  test: build fmtpolice, ginkgo tests, and bats tests"
	@echo
	@echo "  dev: set up the dev toolchain (deps + gox)"

.PHONY: clean
clean:
	go clean -i -r $(PACKAGES) || true
	rm -f $(GOPATH)/bin/docker-builder

.PHONY: quick
quick: build
	@echo "----------"
	@docker-builder --version
	@echo "----------"
	@docker-builder --help
	@echo "----------"
	@docker-builder
	@echo "----------"

.PHONY: binclean
binclean:
	rm -f $(GOPATH)/bin/docker-builder
	rm -rf ./Release/*
	touch ./Release/.gitkeep

.PHONY: build
build: binclean
	CGO_ENABLED=0 go install -a $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(PACKAGES)

.PHONY: release
release: binclean gox-linux gox-darwin
	#open ./Release

.PHONY: gox-linux
gox-linux: build dev
	CGO_ENABLED=0 gox -output="Release/docker-builder-$(REPO_VERSION)-linux-amd64" -osarch="linux/amd64" $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(B)
	# TODO: make checksum

.PHONY: gox-darwin
gox-darwin: build dev
	CGO_ENABLED=0 gox -output="Release/docker-builder-$(REPO_VERSION)-docker-amd64" -osarch="darwin/amd64" $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(B)
	# TODO: make checksum

.PHONY: deps
deps:
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
	  if [[ "$$($(GOPATH)/bin/golint $$file)" =~ ^[[:blank:]]*$$ ]] ; then \
	  echo yayyy >/dev/null ; \
	  else $(MAKE) lintv && exit 1 ; fi \
	  done

.PHONY: lintv
lintv:
	@echo "----------"
	@for file in $(shell git ls-files '*.go') ; do $(GOPATH)/bin/golint $$file ; done

.PHONY: ginkgo
ginkgo:
	@echo "----------"
	@if [[ "$(GINKGO_PATH)" == "." ]] ; then \
	  echo "$(GOPATH)/bin/ginkgo -nodes=10 -noisyPendings -race -r ." && \
	  $(GOPATH)/bin/ginkgo -nodes=10 -noisyPendings -race -r . ; \
	  else echo "$(GOPATH)/bin/ginkgo -nodes=10 -noisyPendings -race --v $(GINKGO_PATH)" && \
	  $(GOPATH)/bin/ginkgo -nodes=10 -noisyPendings -race --v $(GINKGO_PATH) ; \
	  fi

.PHONY: bats
bats:
	@echo "----------"
	$(BATS_INSTALL_DIR)/bin/bats $(BATS_OUT_FORMAT) $(shell find . -type f -name '*.bats')

.PHONY: gox
gox:
	@if which gox >/dev/null ; then \
	  echo "not installing gox, gox already installed." ; \
	  else \
	  go get github.com/mitchellh/gox ; \
	  gox -build-toolchain -osarch="linux/amd64 darwin/amd64" ; \
	  fi \

.PHONY: gopath
gopath:
	@echo  "\$$GOPATH = $(GOPATH)"

.PHONY: get
get:
	go get ./...

.PHONY: dev
dev: deps gox
