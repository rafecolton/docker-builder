SHELL := /bin/bash
SUDO ?= sudo
DOCKER ?= docker
B := github.com/artifactory/build-runner
PACKAGES := ./...

GOPATH := $(shell echo $${GOPATH%%:*})

# go build args
GO_TAG_ARGS ?= -tags netgo

export GOPATH

.PHONY: build
build: get monkey-patch-drone
	go install -a $(PACKAGES)

.PHONY: monkey-patch-drone
monkey-patch-drone:
	@if [[ "$(DRONE)" == "true" ]] && [[ "$(CI)" == "true" ]] ; then rm -f $(GOROOT)/src/pkg/os/error_posix.go ; fi

.PHONY: test
test: fmtpolice
	go test ./...

.PHONY: fmtpolice
fmtpolice: $(PWD)/fmtpolice
	bash fmtpolice

$(PWD)/fmtpolice:
	curl -sL https://raw.githubusercontent.com/rafecolton/fmtpolice/master/fmtpolice -o $@

$(GOPATH)/bin/deppy:
	go get github.com/hamfist/deppy

.PHONY: get
get: $(GOPATH)/bin/deppy
	go get -t ./...
	$(GOPATH)/bin/deppy restore

#$(PWD)/Specs/bin/coverage:
	#curl -sL https://raw.githubusercontent.com/rafecolton/fmtpolice/master/coverage -o $@ && \
	  #chmod +x $@

#.PHONY: coverage
#coverage: $(PWD)/Specs/bin/coverage
	#go get -u code.google.com/p/go.tools/cmd/cover || go get -u golang.org/x/tools/cmd/cover
	#go get -u github.com/axw/gocov/gocov
	#./Specs/bin/coverage

#.PHONY: goveralls
#goveralls: coverage
	#go get -u github.com/mattn/goveralls
	#@echo "goveralls -coverprofile=gover.coverprofile -repotoken <redacted>"
	#@goveralls -coverprofile=gover.coverprofile -repotoken $(GOVERALLS_REPO_TOKEN)
