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
REPO_BRANCH := $(shell git rev-parse -q --abbrev-ref HEAD) # FIXME: will be "HEAD" if not on branch
GOBUILD_VERSION_ARGS := -ldflags "\
  -X $(REV_VAR) $(REPO_REV) \
  -X $(VERSION_VAR) $(REPO_VERSION) \
  -X $(BRANCH_VAR) $(REPO_BRANCH) \
  -w \
"

BATS_INSTALL_DIR ?= $(PWD)/Specs/bats

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
build: binclean get monkey-patch-drone
	CGO_ENABLED=0 go install -a $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(PACKAGES)

.PHONY: monkey-patch-drone
monkey-patch-drone:
	@if [[ "$(DRONE)" == "true" ]] && [[ "$(CI)" == "true" ]] ; then rm -f $(GOROOT)/src/pkg/os/error_posix.go ; fi

.PHONY: release
release: binclean gox-build
	open ./Release

.PHONY: gox-build
gox-build: get $(GOPATH)/bin/gox
	CGO_ENABLED=0 $(GOPATH)/bin/gox -output="Release/docker-builder-$(REPO_VERSION)-{{ .OS }}-{{ .Arch }}" -osarch="darwin/amd64 linux/amd64" $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) $(B)
	for file in $$(find ./Release -type f -name 'docker-builder-*') ; do openssl sha256 -out $$file-SHA256SUM $$file ; done

.PHONY: .test
.test: fmtpolice bats
	go test ./...

.PHONY: test
test:
	@GO_TAG_ARGS="-tags netgo -tags integration" $(MAKE) build
	@DOCKER_BUILDER_TEST_MODE=1 $(MAKE) .test

.PHONY: fmtpolice
fmtpolice: $(PWD)/Specs/bin/fmtpolice
	./Specs/bin/fmtpolice

$(PWD)/Specs/bin/fmtpolice:
	curl -sL https://raw.githubusercontent.com/rafecolton/fmtpolice/master/fmtpolice -o $@ && \
	  chmod +x $@

.PHONY: bats
bats: $(BATS_INSTALL_DIR)/bin/bats
	@echo "----------"
	$(BATS_INSTALL_DIR)/bin/bats $(BATS_OUT_FORMAT) $(shell find . -type f -name '*.bats')

$(BATS_INSTALL_DIR)/bin/bats:
	git clone https://github.com/sstephenson/bats.git && \
		(cd bats && ./install.sh $(BATS_INSTALL_DIR)) && \
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
	$(GOPATH)/bin/deppy restore

$(PWD)/Specs/bin/coverage:
	curl -sL https://raw.githubusercontent.com/rafecolton/fmtpolice/master/coverage -o $@ && \
	  chmod +x $@

.PHONY: coverage
coverage: $(PWD)/Specs/bin/coverage
	go get -u code.google.com/p/go.tools/cmd/cover || go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/axw/gocov/gocov
	./Specs/bin/coverage

.PHONY: goveralls
goveralls: coverage
	go get -u github.com/mattn/goveralls
	@echo "goveralls -coverprofile=gover.coverprofile -repotoken <redacted>"
	@goveralls -coverprofile=gover.coverprofile -repotoken $(GOVERALLS_REPO_TOKEN)
