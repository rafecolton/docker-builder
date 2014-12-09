# Welcome to Bob!

[![Build Status](https://drone.io/github.com/rafecolton/docker-builder/status.png)](https://drone.io/github.com/rafecolton/docker-builder/latest)
[![Build Status](https://travis-ci.org/rafecolton/docker-builder.svg?branch=master)](https://travis-ci.org/rafecolton/docker-builder)
[![GoDoc](https://godoc.org/github.com/rafecolton/docker-builder?status.png)](https://godoc.org/github.com/rafecolton/docker-builder)
[![Coverage Status](https://img.shields.io/coveralls/rafecolton/docker-builder.svg)](https://coveralls.io/r/rafecolton/docker-builder?branch=master)

Bob, the `docker-builder`, builds, tags, and pushes multiple Docker images, all
from a friendly `Bobfile` config file.

## About

This repo contains documentation on the server and CLI features of
docker-builder.  For documentation on how to write a Bobfile, visit
[github.com/winchman/builder-core](https://github.com/winchman/builder-core)

Other useful docs:

* [Original Motivation for Docker Builder](_docs/why.md)
* [Docker Builder Server](_docs/advanced-usage.md)
* [Writing a Bobfile](_docs/writing-a-bobfile.md)
* [Using with TLS (or `boot2docker` on Mac OS X](_docs/using-with-tls.md)

## Getting Started

0. Install `docker-builder`
0. Explore the [Writing a Bobfile](_docs/writing-a-bobfile.md) doc
0. Run`docker-builder --help`
0. Run `docker-builder help build`

### Installing `docker-builder`

#### From source

To build from source, run `make build`.  You may have to install some
things first, such as `go`

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

## New Command: `docker-builder enqueue`

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
