SHELL := /bin/bash
GOPATH := $(shell echo $${GOPATH%%:*})

export GOPATH

.PHONY: test
test: .fmtpolice
	go test ./...
	@find . -type f -name '*.test' -exec rm {} \;

.PHONY: integration
integration:
	INTEGRATION=1 $(MAKE) test

fmtpolice:
	curl -sL https://raw.githubusercontent.com/rafecolton/fmtpolice/master/fmtpolice -o $@

.PHONY:
.fmtpolice: fmtpolice
	bash fmtpolice

.PHONY: get
get:
	go get -d -t ./...

.PHONY: coverage
coverage: $(PWD)/coverage
	go get -u code.google.com/p/go.tools/cmd/cover || go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/axw/gocov/gocov
	bash coverage

$(PWD)/coverage:
	curl -sL https://raw.githubusercontent.com/rafecolton/fmtpolice/master/coverage -o $@

.PHONY: goveralls
goveralls: coverage
	go get -u github.com/mattn/goveralls
	@echo "goveralls -coverprofile=gover.coverprofile -repotoken <redacted>"
	@goveralls -coverprofile=gover.coverprofile -repotoken $(GOVERALLS_TOKEN)
