# Welcome to Bob!

[![Build Status](https://travis-ci.org/rafecolton/docker-builder.svg?branch=master)](https://travis-ci.org/rafecolton/docker-builder)
[![GoDoc](https://godoc.org/github.com/rafecolton/docker-builder?status.png)](https://godoc.org/github.com/rafecolton/docker-builder)
[![Go Report Card](https://goreportcard.com/badge/rafecolton/docker-builder)](https://goreportcard.com/report/rafecolton/docker-builder)
[![Coverage Status](https://coveralls.io/repos/github/rafecolton/docker-builder/badge.svg?branch=master)](https://coveralls.io/github/rafecolton/docker-builder?branch=master)

Bob, the `docker-builder`, builds, tags, and pushes multiple Docker images, all
from a friendly `Bobfile` config file.

## About

This repo contains documentation on the server and CLI features of
docker-builder.  For documentation on how to write a Bobfile, visit
[github.com/winchman/builder-core](https://github.com/winchman/builder-core)

Other useful docs:

* [Original Motivation for Docker Builder](_docs/why.md)
* [Using with TLS (and `boot2docker` on Mac OS X)](_docs/using-with-tls.md)

## Getting Started

0. Install `docker-builder`
0. Run`docker-builder --help`

### Installing `docker-builder`

#### From source

```bash
git clone https://github.com/rafecolton/docker-builder
cd docker-builder
make # clean build test
```

install dependencies as needed

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

## Subcommands

* [`docker-builder enqueue`](_docs/subcommands/enqueue.md) - enqueue a
  build with your cwd
* [`docker-builder serve`](_docs/subcommands/serve.md) - run
  docker-builder as an http server 
* `docker-builder -h/--help/help` - view all subcommands and flags

----

[CONTRIBUTING](CONTRIBUTING.md)
