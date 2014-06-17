# Welcome to Bob!

[![Build Status](https://travis-ci.org/modcloth/docker-builder.svg?branch=master)](https://travis-ci.org/modcloth/docker-builder)
[![GoDoc](https://godoc.org/github.com/modcloth/docker-builder?status.png)](https://godoc.org/github.com/modcloth/docker-builder)

Bob, the `docker-builder`, builds, tags, and pushes multiple Docker images, all
from a friendly `Bobfile` config file.

## Quick Start

Steps to quick start:

```bash
# install docker-builder
go get github.com/modcloth/docker-builder

# create a Bobfile
docker-builder init .

# build the aforementioned Bobfile
docker-builder build Bobfile
```

## More Reasonably-Paced Start
0. Install `docker-builder`
0. Explore the [Writing a Bobfile](docs/writing-a-bobfile.md) doc
0. Consider running `docker-builder -h` and `docker-builder help build` for more
   options

### First, Install `docker-builder`

#### Easiest

```bash
go get github.com/modcloth/docker-builder
```

#### From pre-compiled binaries

```bash
# on Mac OS X
curl -sL https://github.com/modcloth/docker-builder/releases/download/v0.3.2/docker-builder-v0.3.2-darwin-amd64.tar.gz | \
  tar -xzf - -C /usr/local --strip-components=1

# on Linux, note: you may need sudo
curl -sL https://github.com/modcloth/docker-builder/releases/download/v0.3.2/docker-builder-v0.3.2-linux-amd64.tar.gz | \
  sudo tar -xzf - -C /usr/local --strip-components=1
```

These commands will place `docker-builder` at
`/usr/local/bin/docker-builder`, so to use `docker-builder`, make sure
and check that `/usr/local/bin` is in your `$PATH` or change the `-C`
option.

**NOTE:** You may see some junk output when running `tar -xzf`.  This
has something to do with the archives being build on Mac OSX.  The
output is harmless and safe to ignore.

#### From source

To build from source, run `make build`.  You may have to install some
things first, such as `go`.


## Why?

Bob was created out of the need to more easily build, tag, and push
layered docker images.  Beyond what a normal `docker build` would offer,
Bob offers the following:

0. **Build from multiple "Dockerfiles"**
  - In order for a docker build to have
    [context](http://docs.docker.io/reference/builder/), the
`Dockerfile` must be present in the code repo and must be named
"Dockerfile".  Bob makes this possible by performing your builds in a
temporary directory, so you can name your `Dockerfile` whatever you
want.

0. **Includes &amp; Excludes**
  - Sometimes, you want to tailor which of your application's files end
    up in your container, but writing an explicit `ADD` command for each
file and directory is very tedious.  Instead, by using Includes and
Excludes, your temporary build directory will have only exactly the
files you want.  That way, instead of adding each file individually, you
can simply `ADD . <dir>`

0. **Tagging macros**
  - More often than not, in addition to a static tag, it is desirable to
    tag a docker container dynamically with, for example, the git
revision of the associated code repo.  Bob makes this easy for you with
tagging macros.

0. **Seamless, reliable build, tag, &amp; push process**
  - A typical docker build workflow can be a bit tedious and nuanced.
    Bob aims to abstract all of this and make the process much simpler
  - simply write your `Dockerfile` and let Bob take care of the rest!
