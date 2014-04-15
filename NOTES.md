# THIS DOC IS A WIP!

## Goals / Notes

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

## TODO MVP:
* tag `0.1.0`
* https://github.com/sstephenson/bats for integration tests
* https://github.com/mitchellh/gox for building binaries
