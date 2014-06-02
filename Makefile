SHELL := /bin/bash
SUDO ?= sudo
DOCKER ?= docker
B := github.com/modcloth/docker-builder
TARGETS := \
  $(B)/builder \
  $(B)/builderfile \
  $(B)/dclient \
  $(B)/log \
  $(B)/parser \
  $(B)/parser/tag \
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

BATS_OUT_FORMAT=$(shell bash -c "echo $${CI+--tap}")
GOPATH := $(shell echo $${GOPATH%%:*})

export BATS_INSTALL_DIR
export GINKGO_PATH
export GOPATH

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
	@echo "  quick: build + invokes builder a couple times (good for debugging)"
	@echo
	@echo "  build: installing libs plus installing deps"
	@echo
	@echo "  test: build fmtpolice, ginkgo tests, and bats tests"
	@echo
	@echo "  dev: set up the dev toolchain (deps + gox)"

.PHONY: all
all: binclean clean build test

.PHONY: clean
clean:
	go clean -i -r $(TARGETS) || true
	rm -f $(GOPATH)/bin/builder

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
	rm -f $(GOPATH)/bin/builder
	rm -rf ./releases/*
	touch ./releases/.gitkeep

.PHONY: build
build: binclean deps
	go build -o $(GOPATH)/bin/builder $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(B)
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(TARGETS)

.PHONY: release
release: binclean gox-linux gox-darwin

.PHONY: gox-linux
gox-linux: build dev
	mkdir -p ./releases/linux/bin
	gox -output="releases/linux/bin/builder" -arch="amd64" -os="linux" $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(B)
	pushd releases >/dev/null && \
	  tar -czf docker-builder-$(REPO_VERSION)-linux-amd64.tar.gz linux/ && \
	  popd >/dev/null

.PHONY: gox-darwin
gox-darwin: build dev
	mkdir -p ./releases/darwin/bin
	gox -output="releases/darwin/bin/builder" -arch="amd64" -os="darwin" $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(B)
	pushd releases >/dev/null && \
	  tar -czf docker-builder-$(REPO_VERSION)-darwin-amd64.tar.gz darwin/ && \
	  popd >/dev/null

.PHONY: godep
godep:
	go get github.com/tools/godep

.PHONY: deps
deps: godep
	@echo "godep restoring..."
	$(GOPATH)/bin/godep restore
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

.PHONY: save
save:
	godep save -copy=false ./...

.PHONY: dev
dev: deps gox
