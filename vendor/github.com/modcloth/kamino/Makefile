SHELL := /bin/bash
SUDO ?= sudo
DOCKER ?= docker
K := github.com/modcloth/kamino
TARGETS := $(K)

GINKGO_PATH ?= "."
GOPATH := $(shell echo $${GOPATH%%:*})
GOBIN := $(GOPATH)/bin
PATH := $(GOBIN):$(PATH)

export GINKGO_PATH
export GOPATH
export GOBIN
export PATH

default: test

.PHONY: all
all: clean build test

.PHONY: clean
clean:
	go clean -i -r $(TARGETS) || true

.PHONY: build
build: deps
	go install $(TARGETS)

.PHONY: test
test: build fmtpolice ginkgo

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
	  echo "$(GOBIN)/ginkgo -nodes=15 -noisyPendings -race -r ." && \
	  $(GOBIN)/ginkgo -nodes=15 -noisyPendings -race -r . ; \
	  else echo "$(GOBIN)/ginkgo -nodes=15 -noisyPendings -race --v $(GINKGO_PATH)" && \
	  $(GOBIN)/ginkgo -nodes=15 -noisyPendings -race --v $(GINKGO_PATH) ; \
	  fi

.PHONY: save
save:
	godep save -copy=false ./...
