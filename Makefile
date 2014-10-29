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
GO_TAG_ARGS ?= -tags netgo

export GOPATH

.PHONY: all
all: binclean clean build test

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
build: binclean get
	CGO_ENABLED=0 go install -a $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(PACKAGES)

.PHONY: release
release: binclean gox-build
	open ./Release

.PHONY: gox-build
gox-build: get $(GOPATH)/bin/gox
	CGO_ENABLED=0 $(GOPATH)/bin/gox -output="Release/docker-builder-$(REPO_VERSION)-{{ .OS }}-{{ .Arch }}" -osarch="darwin/amd64 linux/amd64" $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(B)
	for file in $$(find ./Release -type f -name 'docker-builder-*') ; do openssl sha256 -out $$file-SHA256SUM $$file ; done

.PHONY: install-ginkgo
install-ginkgo:
	go get -u github.com/onsi/ginkgo/ginkgo

.PHONY: .test
.test: fmtpolice ginkgo bats

.PHONY: test
test:
	@GO_TAG_ARGS="-tags netgo -tags integration" $(MAKE) build
	@$(MAKE) .test

.PHONY: fmtpolice
fmtpolice: fmt lint

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
	go get -u github.com/golang/lint/golint

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
ginkgo: install-ginkgo
	@echo "----------"
	@if [[ "$(GINKGO_PATH)" == "." ]] ; then \
	  echo "$(GOPATH)/bin/ginkgo -nodes=10 -noisyPendings -race -r ." && \
	  $(GOPATH)/bin/ginkgo -nodes=10 -noisyPendings -race -r . ; \
	  else echo "$(GOPATH)/bin/ginkgo -nodes=10 -noisyPendings -race --v $(GINKGO_PATH)" && \
	  $(GOPATH)/bin/ginkgo -nodes=10 -noisyPendings -race --v $(GINKGO_PATH) ; \
	  fi

.PHONY: bats
bats: $(BATS_INSTALL_DIR)/bin/bats
	@echo "----------"
	$(BATS_INSTALL_DIR)/bin/bats $(BATS_OUT_FORMAT) $(shell find . -type f -name '*.bats')

$(BATS_INSTALL_DIR)/bin/bats:
	git clone https://github.com/sstephenson/bats.git && \
		(cd bats && $(SUDO) ./install.sh $(BATS_INSTALL_DIR)) && \
		rm -rf bats

$(GOPATH)/bin/gox:
	go get github.com/mitchellh/gox
	$(GOPATH)/bin/gox -build-toolchain -osarch="linux/amd64 darwin/amd64"

.PHONY: gopath
gopath:
	@echo  "\$$GOPATH = $(GOPATH)"

$(GOPATH)/bin/deppy:
	go get github.com/hamfist/deppy

.PHONY: get
get: $(GOPATH)/bin/deppy
	go get -t ./...
	deppy restore
