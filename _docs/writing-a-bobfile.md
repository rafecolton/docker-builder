# Writing a Bobfile - Version 1

You may also want to look at the following:

* [Writing a Bobfile - Version 0 (deprecated)](writing-a-bobfile-version-zero.md)
* [Upgrading Bobfile from Version 0 to Version 1: How and Why](upgrading-zero-to-one.md)

The basic ingredients to a Bob-build are the `docker-builder` executable and a
`Bobfile` config file.  This page assumes that you have already
downloaded or built the `docker-builder` executable.  If you have not, please
check out the other docs referenced

## Writing a `Bobfile` Config

Here is an example `Bobfile` config file with each section explained.
Bob config files are written in [toml](https://github.com/mojombo/toml).
**NOTE:** The file does not have to be named `Bobfile`.  It can be named
whatever you want, and you can have as many bob config files as you
want.  The name `Bobfile` is just a convention.

```toml
# Bobfile

version = 1 # this is required, or your Bobfile may be parsed incorrectly

[docker]
build_opts = [
  "--rm",
  "--no-cache"
]
tag_opts = ["--force"]

[container_globals]
registry = "rafecolton"
dockercfg_un = "foo"
dockercfg_pass = "bar"
dockercfg_email = "baz@example.com"
project = "my-app"
tags = [
  "{{ branch }}",
  "{{ sha }}",
  "{{ tag }}"
]

[[container]]
name = "base"
Dockerfile = "Dockerfile.base"
tags = ["base"]
skip_push = true

[[container]]
name = "app"
Dockerfile = "Dockerfile.app"
```

### The `[docker]` Section

The `[docker]` section is used for declaring options that will be passed
to the `docker build` and `docker tag` commands.  The following stanzas
are available:

* `build_opts` - Array
* `tag_opts` - Array

### The `[container_globals]` Section

The `[container_globals]` section is a special section that will get
merged into each of the other container sections, with the values in the
individual container section taking precedence over the global section.
For example, the above `Bobfile` could be rewritten as follows:

```toml
# Bobfile
[docker]
build_opts = [
  "--rm",
  "--no-cache"
]
tag_opts = ["-f"]

[[container]]
name = "base"
Dockerfile = "Dockerfile.base"
registry = "rafecolton"
dockercfg_un = "foo"
dockercfg_pass = "bar"
dockercfg_email = "baz@example.com"
project = "my-app"
tags = ["base"]
skip_push = true

[[container]]
name = "app"
Dockerfile = "Dockerfile.app"
registry = "rafecolton"
dockercfg_un = "foo"
dockercfg_pass = "bar"
dockercfg_email = "baz@example.com"
project = "my-app"
tags = [
  "{{ branch }}",
  "{{ sha }}",
  "{{ tag }}"
]
```

### The `[[container]]` Sections

The following stanzas are available in a `[[container]]` section:

* `name` - String (required) - the name of the section
* `Dockerfile` - String (required) - the file to be used as the
  "Dockerfile" for the build
* `registry` - String
* `project` - String
* <del>`excluded` - Array</del> **deprecated**
* <del>`included` - Array</del> **deprecated**
* `tags` - Array
* `skip_push` - Bool - don't run `docker push...` after building this
  container
* Docker registry credentials (for authenticating to a registry when pushing images)
  - `dockercfg_un` - auth username
  - `dockercfg_pass` - auth password
  - `dockercfg_email` - auth email

#### The `tags` Stanza

All tags are evaluated using the golang template package.  There are a
few special functions that may be interpolated in the tag string:

- `{{ branch }}` - the git branch of the app repo (`git rev-parse -q --abbrev-ref HEAD`)
- `{{ sha }}` - the full git sha of the app repo (`git rev-parse -q HEAD`)
- `{{ tag }}` - the shortened version of the app repo rev (`git describe --always --dirty --tags`)

For example, the following tag `myapp-{{ tag }}` might produce the tag
`myapp-v0.1.0`

## Linting &amp; Building

Once you have written your `Bobfile` config file, linting and building
are both very simple.  First, place the `Bobfile` file at the top level
of your application repo. 

Then, to validate your config:

```bash
docker-builder lint <path>/<to>/Bobfile
```

and to build:

```bash
docker-builder build <path>/<to>/Bobfile.whatever

# or, if your config is just named "Bobfile", then from the repo top level...
docker-builder build
```
