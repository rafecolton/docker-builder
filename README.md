# Welcome to Bob!

[![Build Status](https://travis-ci.org/rafecolton/docker-builder.svg?branch=master)](https://travis-ci.org/rafecolton/docker-builder)
[![GoDoc](https://godoc.org/github.com/rafecolton/docker-builder?status.png)](https://godoc.org/github.com/rafecolton/docker-builder)

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

### First, Install `docker-builder`

#### Easiest

```bash
go get github.com/rafecolton/docker-builder
```

#### From pre-compiled binaries

```bash
# on Mac OS X
curl -sL https://github.com/rafecolton/docker-builder/releases/download/v0.9.0/docker-builder-v0.9.0-darwin-amd64 \
  -o /usr/local/bin/docker-builder && chmod +x /usr/local/bin/docker-builder

# on Linux, note: you may need sudo
curl -sL https://github.com/rafecolton/docker-builder/releases/download/v0.9.0/docker-builder-v0.9.0-linux-amd64 \
  -o /usr/local/bin/docker-builder && chmod +x /usr/local/bin/docker-builder
```

**NOTE:** Checksums available on the [release page](https://github.com/rafecolton/docker-builder/releases)

#### From source

To build from source, run `make build`.  You may have to install some
things first, such as `go`

## Contributing

**Pull requests welcome!**
