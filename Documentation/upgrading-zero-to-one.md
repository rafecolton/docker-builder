# Upgrading Bobfile from Version 0 to Version 1

## Why

The previous Bobfile data format was affected by a [change in Go 1.3](http://golang.org/doc/go1.3#map).

This had two affects.  First, it caused one of the parser tests to fail
randomly.  Second, and more important, it caused the container sections
in the Bobfile to be processed out of order.

## How

Here is a full Bobfile example, marked up to indicate how to perform the
conversion.

```toml
# Bobfile version 0
[docker]
build_opts = [
  "--rm",
  "--no-cache"
]
tag_opts = ["--force"]

[containers]

[containers.global]
registry = "rafecolton"
project = "my-app"
excluded = ["spec"]
tags = [
  "git:branch",
  "git:rev",
  "git:short"
]

[containers.base]
Dockerfile = "Dockerfile.base"
included = [
  "Gemfile",
  "Gemfile.lock"
]
tags = ["base"]
skip_push = true

[containers.app]
Dockerfile = "Dockerfile.app"
```

```toml
# Bobfile version 1

version = 1 # <--- this line gets added

[docker]
build_opts = [
  "--rm",
  "--no-cache"
]
tag_opts = ["--force"]

# [containers] # <--- this line is no longer necessary

[container_globals] # <--- this line changed, used to be "[containers.global]"
registry = "rafecolton"
project = "my-app"
excluded = ["spec"]
tags = [
  "git:branch",
  "git:rev",
  "git:short"
]

[[container]] # <--- this line changed, used to be [containers.base]
name = "base" # <--- this line got added
Dockerfile = "Dockerfile.base"
included = [
  "Gemfile",
  "Gemfile.lock"
]
tags = ["base"]
skip_push = true

[[container]] # <--- this line changed, used to be [containers.app]
name = "app" # <--- this line got added
Dockerfile = "Dockerfile.app"

# vim:ft=toml
```
