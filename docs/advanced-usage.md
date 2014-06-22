## Advanced Usage

### Running the Server

```bash
> docker-builder help serve

# NAME:
#    serve - serve <options> - start a small HTTP web server for receiving build requests
#
# USAGE:
#    command serve [command options] [arguments...]
#
# DESCRIPTION:
#    Start a small HTTP web server for receiving build requests.
#
# Configure through the environment:
#
# DOCKER_BUILDER_LOGLEVEL     =>     --log-level (global)
# DOCKER_BUILDER_LOGFORMAT    =>     --log-format (global)
# DOCKER_BUILDER_PORT         =>     --port
# DOCKER_BUILDER_APITOKEN     =>     --api-token
# DOCKER_BUILDER_SKIPPUSH     =>     --skip-push
#
#
# OPTIONS:
#    --port, -P '5000'  port on which to serve
#    --api-token, -T  GitHub API token
#    --skip-push    override Bobfile behavior and do not push any images (useful for testing)
```

### Example Request

```bash
#!/bin/bash

curl -XPOST -H 'Content-Type: application/json' 'http://localhost:5000/docker-build' -d '
{
  "account": "my-account",
  "repo": "my-repo",
  "ref": "master"
}
'
```

### Request Fields

Required Fields:

* `account / type: string` - the GitHub account for the repo being cloned
* `repo / type: string` - the name of the repo
* `ref / type: string` - the ref (can be any valid/unambiguous ref - a branch, tag, sha, etc)

Other Fields: 

* `api_token / type: string` - the GitHub api token (not required for public repos)
* `depth / type: string (must be int > 0)` - clone depth (default: no `--depth` argument passed to `git clone`)
