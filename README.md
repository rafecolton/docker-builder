# Welcome to Bob!

[![Build Status](https://travis-ci.org/rafecolton/docker-builder.svg?branch=master)](https://travis-ci.org/rafecolton/docker-builder)
[![GoDoc](https://godoc.org/github.com/rafecolton/docker-builder?status.png)](https://godoc.org/github.com/rafecolton/docker-builder)
[![Coverage Status](https://img.shields.io/coveralls/rafecolton/docker-builder.svg)](https://coveralls.io/r/rafecolton/docker-builder?branch=master)

Bob, the `docker-builder`, builds, tags, and pushes multiple Docker images, all
from a friendly `Bobfile` config file.

## Other Topics

* [Original Motivation for Docker Builder](Documentation/why.md)
* [Docker Builder Server](Documentation/advanced-usage.md)
* [Writing a Bobfile](Documentation/writing-a-bobfile.md)

## Quick Start

Steps to quick start:

```bash
# install docker-builder
go get github.com/rafecolton/docker-builder

$GOPATH/bin/docker-builder --help
```

## More Reasonably-Paced Start
0. Install `docker-builder`
0. Explore the [Writing a Bobfile](Documentation/writing-a-bobfile.md) doc
0. Run`docker-builder --help`
0. Run `docker-builder help build`

### Installing `docker-builder`

#### Easiest

```bash
go get github.com/rafecolton/docker-builder
```

#### From pre-compiled binaries

```bash
# on Mac OS X
curl -sL https://github.com/rafecolton/docker-builder/releases/download/v0.10.0/docker-builder-v0.10.0-darwin-amd64 \
  -o /usr/local/bin/docker-builder && chmod +x /usr/local/bin/docker-builder

# on Linux, note: you may need sudo
curl -sL https://github.com/rafecolton/docker-builder/releases/download/v0.10.0/docker-builder-v0.10.0-linux-amd64 \
  -o /usr/local/bin/docker-builder && chmod +x /usr/local/bin/docker-builder
```

**NOTE:** Checksums available on the [release page](https://github.com/rafecolton/docker-builder/releases)

#### From source

To build from source, run `make build`.  You may have to install some
things first, such as `go`

### Using with TLS

If you are using a version of `docker` with TLS enabled (supported in
`docker` `v1.3.0` and up, enabled by default with `boot2docker`), you
will need to use `docker-builder` `v0.9.2` or greater.

Additionally, you must set the following environment variables:

```bash
# all values are the boot2docker defaults
export DOCKER_CERT_PATH="$HOME/.boot2docker/certs/boot2docker-vm"
export DOCKER_TLS_VERIFY=1
export DOCKER_HOST="tcp://127.0.0.1:2376"
```

**NOTE:** `docker-builder` will automatically set the correct url scheme
for TLS if you are using port 2376.  If you are using another port and
wish to enable TLS, you must set the following additional environment
variable:

```bash
export DOCKER_HOST_SCHEME="https"
```

## Experimental Feature: Enqueue

Docker-builder has an new, experimental command-line feature: `enqueue`

Use `docker-builder enqueue` to push a build for your *current working
directory* to your Docker build server.  To use locally, run the server
in one tab and `docker-builder enqueue` in another.  For example:

```bash
# terminal window 1
docker-builder serve

# terminal window 2
docker-builder enqueue
```

Or, you may push directly to your build server by setting the
docker-build-server host:

```bash
# via tne environment
export DOCKER_BUILDER_HOST="http://username:password@build-server-host.example.com:5000" && \
  docker-builder enqueue

# via the command line
docker-builder enqueue --host "http://username:password@build-server-host.example.com:5000"
```

## Contributing

**Pull requests welcome!**
