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
#   -l, --lint=
#   -f, --builderfile=
#
# Help Options:
#   -h, --help         Show this help message
```

-----

# The rest of this doc is a WIP

## Goals

Ideas:

* like `docker build` but on steroids
* reliably and repeatly build docker containers
* use the docker api / daemon
* smart about detecting if a container is already built or not
* config via env or config file
* add some extra options like excluding directories
* can have a generic build file
* http api
* it uses its own api - can run as a daemon or as a single instance
i.e. its own client calls

~~~~~~

runs docker inside a container


## Technical Requirements

Requirements:

* written in Go
* works on linux and mac
* TDD using Ginkgo
* manage packages (Godep?)
* easy to build
* prebuilt binaries


Builderfile = TOML config
* global options
* per project options
* per project docker file filename if desired
  - or just single script to be run (must be executable)
* list of ports
* app name
* directory in the container and purpose (we have a standard place to
  put them)


Accepts requests via:
HTTP (MVP)
git hook
