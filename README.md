bob
===

[![Build Status](https://travis-ci.org/rafecolton/bob.svg?branch=master)](https://travis-ci.org/rafecolton/bob)
[![GoDoc](https://godoc.org/github.com/rafecolton/bob?status.png)](https://godoc.org/github.com/rafecolton/bob)

builds a docker image from an arbitrary file

## Why?

* base layers
* can only add in one "Dockerfile"
* can't easily exclude dirs / long file names (aufs limitation?)
* generally make it more efficient to build in smaller layers to make
  pushing and pulling faster, easier, more reliable, etc

## Installing

When the binaries are available, installing will look something like this:

```bash
# on Mac OS X
curl -sL https://github.com/rafecolton/bob/releases/download/0.0.1-alpha/darwin-amd64.tar.gz | \
  tar -xzf - -C /usr/local --strip-components=1
# on Linux, note: you may need sudo
curl -sL https://github.com/rafecolton/bob/releases/download/0.0.1-alpha/linux-amd64.tar.gz | \
  sudo tar -xzf - -C /usr/local --strip-components=1
```

## Building Containers

Example usage:

```bash
# verify your file bob.toml file is valid
builder --lint bob.toml

# build your containers
builder --build bob.toml
```

For other usages:

```bash
builder -h/--help

# Usage:
#   builder [OPTIONS]
# 
# Application Options:
#   -v             Print version and exit
#       --version  Print long version and exit
#       --branch   Print branch and exit
#       --rev      Print revision and exit
#   -q, --quiet    Produce no output, only exit codes (false)
#   -l, --lint=    Lint the provided file. Compatible with -q/--quiet
#   -b, --build=   The configuration file for Builder
# 
# Help Options:
#   -h, --help     Show this help message
```

## Hacking

```bash
make help

# Usage: make [target]
#
# Options:
#
#   help/default: display this message
#
#   all: clean build test
#
#   quick: build + invokes builder a couple times (good for debugging)
#
#   build: gvm linkthis plus installing libs plus installing deps
#
#   test: build fmtpolice and ginkgotests
#
#   dev: set up the dev tool chain
```
