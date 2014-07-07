## Motivation for Docker Builder

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
