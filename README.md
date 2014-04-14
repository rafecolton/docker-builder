builder
=======

[![Build Status](https://travis-ci.org/rafecolton/builder.svg?branch=master)](https://travis-ci.org/rafecolton/builder)

builds a docker image from an arbitrary file

## Hacking

```bash
> make help
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
#
```

## Building

```bash
> builder -h/--help
# Usage:
#   builder [OPTIONS]
#
# Application Options:
#   -v                 Print version and exit
#       --version      Print long version and exit
#       --branch       Print branch and exit
#       --rev          Print revision and exit
#   -q, --quiet        Produce no output, only exit codes (false)
#   -l, --lint=        Lint the provided file. Compatible with -q/--quiet
#   -f, --builderfile= The configuration file for Builder
#
# Help Options:
#   -h, --help         Show this help message
```
